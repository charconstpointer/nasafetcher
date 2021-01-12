package pics

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type client interface {
	Get(ctx context.Context, url string) ([]byte, error)
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
func (c *NASAClient) returnToken() {
	c.tokens <- struct{}{}
}
func (c *NASAClient) Get(ctx context.Context, url string) ([]byte, error) {
	select {
	case <-c.tokens:
		defer c.returnToken()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		res, err := http.DefaultClient.Do(req)
		if res != nil && res.StatusCode == http.StatusTooManyRequests {
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
