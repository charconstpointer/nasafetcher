package main

import (
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/charconstpointer/TWljaGFsIEdvZ29BcHBzIE5BU0E/pkg/pics"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	s := pics.NewFetchServer()
	s.Listen()

	<-sigs
}
