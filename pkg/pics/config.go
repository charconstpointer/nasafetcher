package pics

import "time"

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
