package pics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

//NASAImage is a single image response from api.nasa.gov
type NASAImage struct {
	Copyright       string `json:"copyright"`
	Date            string `json:"date"`
	Explanation     string `json:"explanation"`
	Hdurl           string `json:"hdurl"`
	Media_type      string `json:"media_type"`
	Service_version string `json:"service_version"`
	Title           string `json:"title"`
	Url             string `json:"url"`
}

type NASAFetcher struct {
	tokens chan struct{}
	apiKey string
	api    string
	c      client
	logger Logger
}

type TooManyRequests struct {
	Path string
}

func (e *TooManyRequests) Error() string {
	return fmt.Sprintf("too many requests")
}

func NewNASAFetcher(config *Config, nasaClient client) *NASAFetcher {
	var c client
	if nasaClient == nil {
		c = NewNASAClient(config)
	} else {
		c = nasaClient
	}
	if config == nil {
		config = newDefaultConfig()
	}
	fetcher := NASAFetcher{
		apiKey: config.APIKey,
		api:    "https://api.nasa.gov/planetary/apod",
		c:      c,
		logger: config.Logger,
	}

	return &fetcher
}

//if you want additional filtering, like mentioned copyright Filter is the way to do it
func (n *NASAFetcher) buildUrl(start time.Time, end time.Time, filters ...Filter) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s?api_key=%s", n.api, n.apiKey))

	for _, filter := range filters {
		//knowing that api_key is always present we can simply apend our filters with '&'%s=%s
		sb.WriteString(fmt.Sprintf("&%s=%s", filter.key, filter.value))
	}
	sb.WriteString(fmt.Sprintf("&date=%s", start.Format("2006-01-02")))
	return sb.String()
}

func (n *NASAFetcher) getJobs(start time.Time, end time.Time, filters ...Filter) ([]string, error) {
	if start.After(time.Now()) {
		return nil, errors.New("star time cannot be in the past")
	}
	if start.After(end) {
		return nil, errors.New("start time must be before end time")
	}

	jobs := make([]string, 0)
	for start.Before(end.Add(time.Hour * 24)) {
		if start.After(time.Now()) {
			n.logger.Info("trimming dates from the future")
			return jobs, nil
		}
		url := n.buildUrl(start, end, filters...)
		jobs = append(jobs, url)
		start = start.Add(time.Hour * 24)
	}
	return jobs, nil
}

func (n *NASAFetcher) execJobs(ctx context.Context, jobs []string) ([]*NASAImage, error) {
	images := make([]*NASAImage, 0)
	g, _ := errgroup.WithContext(context.Background())
	//could replace it with a result channel and collect it after execution
	mutex := sync.Mutex{}

	for _, job := range jobs {
		j := job
		g.Go(func() error {
			b, err := n.c.Get(ctx, j)
			if err != nil {
				n.logger.Info(err.Error())
				return err
			}
			var img NASAImage
			err = json.Unmarshal(b, &img)

			if err != nil {
				return err
			}

			mutex.Lock()
			images = append(images, &img)
			mutex.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		n.logger.Info(err.Error())
		return nil, err
	}

	return images, nil
}

func (n *NASAFetcher) GetImages(ctx context.Context, start time.Time, end time.Time, filters ...Filter) (*FetchResult, error) {
	//jobs are atomic api calls, given 3 day date range, each job represents single day
	jobs, err := n.getJobs(start, end, filters...)
	if err != nil {
		return nil, err
	}

	imgs, err := n.execJobs(ctx, jobs)
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0)
	for _, i := range imgs {
		urls = append(urls, i.Url)
	}
	return &FetchResult{Urls: urls}, nil
}
