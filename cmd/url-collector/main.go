package main

import (
	"gitlab.com/charconstpointer/TWljaGFsIEdvZ29BcHBzIE5BU0E/pkg/pics"
)

func main() {
	cfg := pics.Config{
		Layout: "2006-01-02",
		Conc:   5,
		Port:   8080,
		Logger: pics.NewPicsLogger(),
	}

	s := pics.NewServer(cfg)
	s.Listen()
}
