package pics

import (
	"encoding/json"
	"net/http"
	"time"
)

type Config struct {
	layout string
	conc   int
}

type FetchServer struct {
	fetcher Fetcher
	config  Config
	tokens  chan struct{}
}

func NewFetchServer() *FetchServer {
	s := FetchServer{
		config: Config{
			layout: "2006-01-02",
			conc:   5,
		},
		fetcher: NewNASAFetcher(5),
	}
	return &s
}
func (s *FetchServer) getPictures() {

}

func (s *FetchServer) handleGetPictures(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	start := r.URL.Query().Get("start_time")
	end := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(s.config.layout, start)
	endTime, err := time.Parse(s.config.layout, end)

	img, err := s.fetcher.GetImages(startTime, endTime, s.config.conc)

	b, err := json.Marshal(img)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)
}

func (s *FetchServer) Listen() error {
	http.HandleFunc("/pictures", s.handleGetPictures)
	return http.ListenAndServe(":8080", nil)
}
