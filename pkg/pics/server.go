package pics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type server struct {
	fetcher    Fetcher
	tokens     chan struct{}
	timeLayout string
	port       int
	logger     Logger
}

type Config struct {
	Layout  string
	Conc    int
	Port    int
	APIKey  string
	Logger  Logger
	Timeout time.Duration
}

func newDefaultConfig() *Config {
	return &Config{
		APIKey:  "DEMO_KEY",
		Logger:  NewPicsLogger(),
		Conc:    5,
		Port:    8080,
		Timeout: time.Second,
	}
}

func NewServer(config *Config, fetcher Fetcher) *server {
	f := fetcher
	s := server{
		fetcher:    f,
		logger:     config.Logger,
		port:       config.Port,
		timeLayout: config.Layout,
	}
	return &s
}

func (s *server) handleGetPictures(w http.ResponseWriter, r *http.Request) {
	type success struct {
		Urls []string `json:"urls"`
	}

	type failure struct {
		Error string `json:"error"`
	}
	if r.Method != http.MethodGet {
		return
	}
	start := r.URL.Query().Get("start_time")
	end := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(s.timeLayout, start)
	endTime, err := time.Parse(s.timeLayout, end)

	img, err := s.fetcher.GetImages(context.Background(), startTime, endTime)

	if err != nil {
		var res failure
		switch err.(type) {
		case *TooManyRequests:
			w.WriteHeader(http.StatusTooManyRequests)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
		res = failure{
			Error: err.Error(),
		}
		b, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
		return
	}
	res := success{
		Urls: img.Urls,
	}
	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (s *server) Listen() error {
	http.HandleFunc("/pictures", s.withLogging(s.handleGetPictures))
	s.logger.Infof("Server started listening on port %d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *server) withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("handling request on %s", r.URL.Path)
		h(w, r)
	}
}
