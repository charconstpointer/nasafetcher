package main

import (
	"flag"

	"gitlab.com/charconstpointer/TWljaGFsIEdvZ29BcHBzIE5BU0E/pkg/pics"
)

var (
	port   = flag.Int("port", 8080, "http port")
	conc   = flag.Int("conc", 5, "max concurrect go routines")
	layout = flag.String("layout", "2006-01-02", "date layout for time parsing")
)

//This project uses only go stdlib as i didn't really feel the need for external libs considering the size of this project
func main() {
	flag.Parse()

	cfg := pics.Config{
		Layout: *layout,
		Conc:   *conc,
		Port:   *port,
		Logger: pics.NewPicsLogger(),
	}

	s := pics.NewServer(cfg)
	s.Listen()
}
