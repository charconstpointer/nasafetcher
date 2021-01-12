package pics

import (
	"errors"
	"io/ioutil"
	"log"
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
		log.Println("getting", url)
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(res.Body)
	case <-time.After(c.timeout):
		return nil, errors.New("concurrency limit reached, all go routines are busy, please retry later")
	}
}
