package pics

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type server struct {
	timeLayout string
	port       int
	mux        *http.ServeMux
	fetcher    Fetcher
	tokens     chan struct{}
	logger     Logger
}

func NewServer(config *Config, fetcher Fetcher) *server {
	s := server{
		mux:        http.NewServeMux(),
		fetcher:    fetcher,
		logger:     config.Logger,
		port:       config.Port,
		timeLayout: config.Layout,
	}
	s.routes()
	return &s
}

func (s *server) routes() {
	s.mux.HandleFunc("/pictures", s.withLogging(s.handleGetPictures()))
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) handleGetPictures() http.HandlerFunc {
	type success struct {
		Urls []string `json:"urls"`
	}

	type failure struct {
		Error string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if r.Method != http.MethodGet {
			return
		}
		start := r.URL.Query().Get("start_time")
		end := r.URL.Query().Get("end_time")

		startTime, err := time.Parse(s.timeLayout, start)
		endTime, err := time.Parse(s.timeLayout, end)

		if err != nil {
			res := failure{
				Error: "cannot parse provide dates, please make sure they're in correct format",
			}
			s.respond(w, res, http.StatusBadRequest)
			return
		}

		img, err := s.fetcher.GetImages(ctx, startTime, endTime)

		if err != nil {
			statusCode := http.StatusBadRequest
			switch err.(type) {
			case *TooManyRequests:
				statusCode = http.StatusTooManyRequests
			default:
				statusCode = http.StatusBadRequest
			}
			res := failure{
				Error: err.Error(),
			}
			s.respond(w, res, statusCode)
			return
		}
		if img == nil {
			res := failure{
				Error: "could not find any images for provided range of dates",
			}
			s.respond(w, res, http.StatusNotFound)
		}
		res := success{
			Urls: img.Urls,
		}
		s.respond(w, res, http.StatusOK)
	}
}

func (s *server) respond(w http.ResponseWriter, data interface{}, code int) {
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(b)
	if err != nil {
		s.logger.Info(err.Error())
	}
}

func (s *server) withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("handling request on %s", r.URL.Path)
		h(w, r)
	}
}
