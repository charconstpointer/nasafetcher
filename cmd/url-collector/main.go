package main

import (
	"flag"
	"fmt"
	"net/http"

	"gitlab.com/charconstpointer/TWljaGFsIEdvZ29BcHBzIE5BU0E/pkg/pics"
)

var (
	port   = pics.GetEnvInt("PORT", 8080)
	conc   = pics.GetEnvInt("CONCURRENT_REQUESTS", 5)
	key    = pics.GetEnvString("API_KEY", "DEMO_KEY")
	layout = flag.String("layout", "2006-01-02", "date layout for time parsing")
)

//This project uses only go stdlib as i didn't really feel the need for external libs considering the size of this project
func main() {
	flag.Parse()

	cfg := pics.Config{
		Layout: *layout,
		Conc:   conc,
		Port:   port,
		APIKey: key,
		Logger: pics.NewPicsLogger(),
	}

	client := pics.NewNASAClient(&cfg)
	fetcher := pics.NewNASAFetcher(&cfg, client)
	server := pics.NewServer(&cfg, fetcher)

	addr := fmt.Sprintf(":%d", cfg.Port)

	if err := http.ListenAndServe(addr, server); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}

}
