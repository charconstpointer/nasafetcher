package pics

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type client interface {
	Get(url string) ([]byte, error)
}

type NASAClient struct {
	tokens  chan struct{}
	timeout time.Duration
}

func NewNASAClient(maxConc int, timeout time.Duration) *NASAClient {
	c := NASAClient{
		tokens:  make(chan struct{}, maxConc),
		timeout: timeout,
	}
	for i := 0; i < maxConc; i++ {
		c.tokens <- struct{}{}
	}

	return &c
}

func (c *NASAClient) Get(url string) ([]byte, error) {
	select {
	case <-c.tokens:
		res, err := http.Get(url)

		if res.StatusCode == http.StatusTooManyRequests {

			return nil, &TooManyRequests{}
		}
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(res.Body)
	case <-time.After(c.timeout):
		return nil, errors.New("concurrency limit reached, all go routines are busy, please retry later")
	}
}
