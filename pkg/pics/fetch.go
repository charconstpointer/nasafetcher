package pics

import (
	"context"
	"time"
)

type Fetcher interface {
	GetImages(ctx context.Context, start time.Time, end time.Time, filters ...Filter) (*FetchResult, error)
}
type Filter struct {
	key   string
	value string
}
type FetchResult struct {
	Urls []string
}
