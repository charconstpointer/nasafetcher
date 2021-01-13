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
	cfg := Config{
		Conc:   10,
		Logger: NewPicsLogger(),
	}
	n := NewNASAFetcher(&cfg, nil)
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
	cfg := Config{
		Conc:   10,
		Logger: NewPicsLogger(),
	}
	n := NewNASAFetcher(&cfg, nil)
	d, err := n.getJobs(startDate, endDate)
	if err == nil {
		t.Errorf("Expected error to not be nil, because start date %s cannot be after end date %s", start, end)
	}

	if d != nil {
		t.Errorf("Expected days to be nil as there are not valid days withing provided range of dates")
	}
}

func TestBuildUrl(t *testing.T) {
	cfg := Config{
		Conc:   10,
		Logger: NewPicsLogger(),
	}
	n := NewNASAFetcher(&cfg, nil)
	start := "2010-01-01"
	end := "2010-01-03"
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)

	url := n.buildUrl(startDate, endDate, Filter{
		key:   "copyright",
		value: "Foo Bar-Baz",
	})
	expected := fmt.Sprintf("%s?api_key=%s&copyright=Foo Bar-Baz&date=%s", n.api, n.apiKey, start)
	if url != expected {
		t.Errorf("Expected %s got %s", expected, url)
	}
}

func TestGetImages(t *testing.T) {
	c := NewMockClient(3, time.Second, time.Millisecond*10)
	cfg := Config{
		Conc:   10,
		Logger: NewPicsLogger(),
	}
	n := NewNASAFetcher(&cfg, c)
	start := "2010-01-01"
	end := "2010-01-03"
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)
	jobs, _ := n.getJobs(startDate, endDate)
	imgs, err := n.GetImages(context.Background(), startDate, endDate)
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
	c := NewMockClient(3, time.Second, time.Millisecond*10)
	cfg := Config{
		Conc:   10,
		Logger: NewPicsLogger(),
	}
	n := NewNASAFetcher(&cfg, c)
	start := "2010-01-04"
	end := "2010-01-01"
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)

	imgs, err := n.GetImages(context.Background(), startDate, endDate)
	if err == nil {
		t.Error("Expected error due to invalid range")
	}
	if imgs != nil {
		t.Error("Expected images to be nil, as there is nothing to be returned due to invalid range error")
	}
}

func TestGetImagesTimeout(t *testing.T) {
	c := NewMockClient(3, time.Second, time.Millisecond*999)
	cfg := Config{
		Conc:   10,
		Logger: NewPicsLogger(),
	}
	n := NewNASAFetcher(&cfg, c)
	start := "2010-01-04"
	end := "2010-01-10"
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)

	imgs, err := n.GetImages(context.Background(), startDate, endDate)

	if err == nil {
		t.Error("expected error to not be nil when concurrency limit is reached")
	}
	if imgs != nil {
		t.Error("Expected images to be nil, as there is nothing to be returned due concurrency limit reached")
	}
}
