package pics

import (
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
	n := NewNASAFetcher(nil)
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
	n := NewNASAFetcher(nil)
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
	f := NewNASAFetcher(nil)
	start := "2010-01-03"
	end := "2010-01-01"
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

}
