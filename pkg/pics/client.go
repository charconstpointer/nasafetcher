package pics

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

//client represents http client
type client interface {
	Get(ctx context.Context, url string) ([]byte, error)
}

//NASAClient is a implementation of client with additional support for concurrency control using semaphore like channel solution

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

//We probably should depend on http.Client instead of using a default one
func (c *NASAClient) Get(ctx context.Context, url string) ([]byte, error) {
	select {
	case <-c.tokens:
		defer c.returnToken()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		//We probably should depend on http.Client in *NASAClient instead of using a default one
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
