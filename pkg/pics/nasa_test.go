package pics

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGetJobsCount(t *testing.T) {
	type test struct {
		start  string
		end    string
		layout string
		want   int
	}
	tests := []test{
		{start: "2010-01-01", end: "2010-01-10", layout: "2006-01-02", want: 10},
		{start: "2010-01-01", end: "2010-01-01", layout: "2006-01-02", want: 1},
	}
	n := NewNASAFetcher(nil, nil)
	for _, tc := range tests {
		start, _ := time.Parse("2006-01-02", tc.start)
		end, _ := time.Parse("2006-01-02", tc.end)
		d, err := n.getJobs(start, end)
		if err != nil {
			t.Error(err.Error())
		}
		if len(d) != tc.want {
			t.Errorf("expected to get %d days instead got %d", tc.want, len(d))
		}
	}

}

func TestGetJobsError(t *testing.T) {
	start := "2010-01-03"
	end := "2010-01-01"
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)
	n := NewNASAFetcher(nil, nil)
	d, err := n.getJobs(startDate, endDate)
	if err == nil {
		t.Errorf("Expected error to not be nil, because start date %s cannot be after end date %s", start, end)
	}

	if d != nil {
		t.Errorf("Expected days to be nil as there are not valid days withing provided range of dates")
	}
}
func TestGetDaysCount(t *testing.T) {
	type test struct {
		start  string
		end    string
		layout string
		want   int
	}
	tests := []test{
		{start: "2010-01-01", end: "2010-01-10", layout: "2006-01-02", want: 10},
		{start: "2010-01-01", end: "2010-01-01", layout: "2006-01-02", want: 1},
	}
	for _, tc := range tests {
		start, _ := time.Parse("2006-01-02", tc.start)
		end, _ := time.Parse("2006-01-02", tc.end)
		d, err := getDays(start, end)
		if err != nil {
			t.Error(err.Error())
		}
		if len(d) != tc.want {
			t.Errorf("expected to get %d days instead got %d", tc.want, len(d))
		}
	}

}

func TestGetDaysError(t *testing.T) {
	start := "2010-01-03"
	end := "2010-01-01"
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)
	d, err := getDays(startDate, endDate)
	if err == nil {
		t.Errorf("Expected error to not be nil, because start date %s cannot be after end date %s", start, end)
	}

	if d != nil {
		t.Errorf("Expected days to be nil as there are not valid days withing provided range of dates")
	}
}

func TestBuildUrl(t *testing.T) {
	f := NewNASAFetcher(nil, nil)
	start := "2010-01-01"
	end := "2010-01-03"
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)

	url := f.buildUrl(startDate, endDate, Filter{
		key:   "copyright",
		value: "Foo Bar-Baz",
	})
	expected := fmt.Sprintf("%s?api_key=%s&copyright=Foo Bar-Baz&date=%s", f.api, f.apiKey, start)
	if url != expected {
		t.Errorf("Expected %s got %s", expected, url)
	}
}

func TestGetImages(t *testing.T) {
	c := NewMockClient(3, time.Second)
	f := NewNASAFetcher(nil, c)
	start := "2010-01-01"
	end := "2010-01-03"
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)
	jobs, _ := f.getJobs(startDate, endDate)
	imgs, err := f.GetImages(context.Background(), startDate, endDate)
	if err != nil {
		t.Error(err.Error())
	}

	if imgs == nil {
		t.Error("with valid range provided, exected result to not be nil")
	}

	if imgs.Urls == nil {
		t.Error("with valid range provided, exected urls to not be nil")
	}

	if len(imgs.Urls) != len(jobs) {
		t.Errorf("Expected imgs count to be %d instead got %d", len(jobs), len(imgs.Urls))
	}
}

func TestGetImagesInvalidRange(t *testing.T) {
	c := NewMockClient(3, time.Second)
	f := NewNASAFetcher(nil, c)
	start := "2010-01-04"
	end := "2010-01-01"
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)

	imgs, err := f.GetImages(context.Background(), startDate, endDate)
	if err == nil {
		t.Error("Expected error due to invalid range")
	}
	if imgs != nil {
		t.Error("Expected images to be nil, as there is nothing to be returned due to invalid range error")
	}
}
