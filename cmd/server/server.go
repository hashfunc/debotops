package main

import (
	"log"

	"github.com/hashfunc/debotops/server"
)

func main() {
	apiserver, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(
		apiserver.Start(),
	)
}
