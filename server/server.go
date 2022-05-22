package main

import (
	"log"
	"net/http"
)

func greetingHandler(writer http.ResponseWriter, _ *http.Request) {
	message := "Greeting from API server"
	if _, err := writer.Write([]byte(message)); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", greetingHandler)

	log.Fatal(
		http.ListenAndServe(":8386", nil),
	)
}
