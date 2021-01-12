package pics

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
}

func NewNASAFetcher(concLimit int) *NASAFetcher {
	fetcher := NASAFetcher{
		tokens: make(chan struct{}, 5),
	}
	for i := 0; i < 5; i++ {
		fetcher.tokens <- struct{}{}
	}

	return &fetcher
}

func getDays(start time.Time, end time.Time) ([]time.Time, error) {
	if start.After(end) {
		return nil, errors.New("start time must be before end time")
	}
	days := make([]time.Time, 0)
	for start.Before(end) {
		days = append(days, start)
		start = start.Add(time.Hour * 24)
	}
	return days, nil
}

func getJobs(start time.Time, end time.Time) ([]string, error) {
	if start.After(end) {
		return nil, errors.New("start time must be before end time")
	}
	jobs := make([]string, 0)
	for start.Before(end) {
		jobs = append(jobs, "https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY&date="+start.Format("2006-01-02"))
		start = start.Add(time.Hour * 24)
	}
	return jobs, nil
}

func (n *NASAFetcher) getImages(jobs []string) ([]*NASAImage, error) {
	log.Println(jobs)
	queue := make(chan string, len(jobs))
	for _, job := range jobs {
		queue <- job
	}
	images := make([]*NASAImage, 0)
	for len(images) < len(jobs) {
		select {
		case _ = <-n.tokens:
			job := <-queue
			res, err := http.Get(job)
			if err != nil {
				log.Println(err.Error())
			}
			var img NASAImage
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(b, &img)
			if err != nil {
				return nil, err
			}
			images = append(images, &img)
			n.tokens <- struct{}{}
		}
	}

	return images, nil
}

func (n *NASAFetcher) GetImages(start time.Time, end time.Time, concLimit int) (*FetchResult, error) {
	jobs, err := getJobs(start, end)
	if err != nil {
		return nil, err
	}
	imgs, err := n.getImages(jobs)
	urls := make([]string, 0)
	for _, i := range imgs {
		urls = append(urls, i.Url)
	}
	return &FetchResult{Urls: urls}, nil
}
