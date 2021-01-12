package pics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
}

func NewNASAFetcher(nasaClient client) *NASAFetcher {
	var c client
	if nasaClient == nil {
		c = NewNASAClient(5, time.Second)
	} else {
		c = nasaClient
	}

	fetcher := NASAFetcher{
		apiKey: "DEMO_KEY",
		api:    "https://api.nasa.gov/planetary/apod",
		c:      c,
	}

	return &fetcher
}

func getDays(start time.Time, end time.Time) ([]time.Time, error) {
	if start.After(end) {
		return nil, errors.New("start time must be before end time")
	}
	days := make([]time.Time, 0)
	for start.Before(end.Add(time.Hour * 24)) {
		days = append(days, start)
		start = start.Add(time.Hour * 24)
	}
	return days, nil
}
func (n *NASAFetcher) buildUrl(start time.Time, end time.Time, filters ...Filter) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s?api_key=%s", n.api, n.apiKey))

	for _, filter := range filters {
		sb.WriteString(fmt.Sprintf("&%s=%s", filter.key, filter.value))
	}
	sb.WriteString(fmt.Sprintf("&date=%s", start.Format("2006-01-02")))
	return sb.String()
}

func (n *NASAFetcher) getJobs(start time.Time, end time.Time, filters ...Filter) ([]string, error) {
	if start.After(end) {
		return nil, errors.New("start time must be before end time")
	}
	jobs := make([]string, 0)
	for start.Before(end.Add(time.Hour * 24)) {
		url := n.buildUrl(start, end, filters...)
		jobs = append(jobs, url)
		start = start.Add(time.Hour * 24)
	}
	return jobs, nil
}

func (n *NASAFetcher) getImages(jobs []string, ctx context.Context) ([]*NASAImage, error) {
	images := make([]*NASAImage, 0)

	g, err := errgroup.WithContext(ctx)

	if err != nil {
		// return nil, err
	}
	mutex := sync.Mutex{}
	for _, job := range jobs {

		g.Go(func() error {
			b, err := n.c.Get(job)
			if err != nil {
				return err
			}
			var img NASAImage
			err = json.Unmarshal(b, &img)
			log.Println(string(b))
			if err != nil {
				return err
			}
			mutex.Lock()
			images = append(images, &img)
			mutex.Unlock()
			return nil
		})
	}
	g.Wait()

	return images, nil
}

func (n *NASAFetcher) GetImages(start time.Time, end time.Time, filters ...Filter) (*FetchResult, error) {
	jobs, err := n.getJobs(start, end, filters...)
	if err != nil {
		return nil, err
	}
	imgs, err := n.getImages(jobs, context.Background())
	urls := make([]string, 0)
	for _, i := range imgs {
		urls = append(urls, i.Url)
	}
	return &FetchResult{Urls: urls}, nil
}
