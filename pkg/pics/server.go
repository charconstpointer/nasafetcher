package pics

import (
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
	Layout string
	Conc   int
	Port   int
	Logger Logger
}

func NewServer(config Config) *server {
	s := server{
		fetcher:    NewNASAFetcher(nil),
		logger:     config.Logger,
		port:       config.Port,
		timeLayout: config.Layout,
	}
	return &s
}

func (s *server) handleGetPictures(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	start := r.URL.Query().Get("start_time")
	end := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(s.timeLayout, start)
	endTime, err := time.Parse(s.timeLayout, end)

	img, err := s.fetcher.GetImages(startTime, endTime)

	b, err := json.Marshal(img)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
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
