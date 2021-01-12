package pics

import "time"

type Fetcher interface {
	GetImages(start time.Time, end time.Time, concLimit int) (*FetchResult, error)
}

type FetchResult struct {
	Urls []string
}
