package main

import (
	"log"

	"github.com/hashfunc/debotops/server"
)

func main() {
	apiserver := server.NewServer()
	log.Fatal(apiserver.Start())
}
