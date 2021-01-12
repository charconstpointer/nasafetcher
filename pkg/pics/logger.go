package pics

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
}

type picsLogger struct {
}

func NewPicsLogger() *picsLogger {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	return &picsLogger{}
}
func (l *picsLogger) Info(args ...interface{}) {
	log.Info(args)
}

func (l *picsLogger) Infof(format string, args ...interface{}) {
	log.Infof(format, args)
}
