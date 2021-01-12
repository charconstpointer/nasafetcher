package pics

import (
	"errors"
	"time"
)

type MockClient struct {
	tokens  chan struct{}
	timeout time.Duration
}

func NewMockClient(maxConc int, timeout time.Duration) *MockClient {
	c := MockClient{
		tokens:  make(chan struct{}, maxConc),
		timeout: timeout,
	}
	for i := 0; i < maxConc; i++ {
		c.tokens <- struct{}{}
	}

	return &c
}

func (c *MockClient) Get(url string) ([]byte, error) {
	select {
	case <-c.tokens:
		res := []byte("{\"copyright\":\"Eric Coles\",\"date\":\"2020-01-07\",\"explanation\":\"Rippling dust and gas lanes give the Flaming Star Nebula its name.  The orange and purple colors of the nebula are present in different regions and are created by different processes.  The bright star AE Aurigae, visible toward the image left, is so hot it is blue, emitting light so energetic it knocks electrons away from surrounding gas.  When a proton recaptures an electron, red light is frequently emitted (depicted here in orange). The purple region's color is a mix of this red light and blue light emitted by AE Aurigae but reflected to us by surrounding dust. The two regions are referred to as emission nebula and reflection nebula, respectively.  Pictured here in the Hubble color palette, the Flaming Star Nebula, officially known as IC 405, lies about 1500 light years distant, spans about 5 light years, and is visible with a small telescope toward the constellation of the Charioteer (Auriga).\",\"hdurl\":\"https://apod.nasa.gov/apod/image/2001/IC405hp_ColesHelm_3447.jpg\",\"media_type\":\"image\",\"service_version\":\"v1\",\"title\":\"IC 405: The Flaming Star Nebula\",\"url\":\"https://apod.nasa.gov/apod/image/2001/IC405hp_ColesHelm_960.jpg\"}")
		time.Sleep(50 * time.Microsecond)
		return res, nil
	case <-time.After(c.timeout):
		return nil, errors.New("concurrency limit reached, all go routines are busy, please retry later")
	}
}
