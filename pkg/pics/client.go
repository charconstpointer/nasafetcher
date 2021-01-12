package pics

import (
	"context"
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

func NewNASAClient(config *Config) *NASAClient {
	c := NASAClient{
		tokens:  make(chan struct{}, config.Conc),
		timeout: config.Timeout,
	}
	for i := 0; i < config.Conc; i++ {
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
	}
}
