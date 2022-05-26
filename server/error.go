package server

import (
	"log"
	"net/http"
)

func Error400(writer http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
	}
	http.Error(writer, "400 Bad Request", http.StatusBadRequest)
}

func Error405(write http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
	}
	http.Error(write, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func Error500(writer http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
	}
	http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
}
