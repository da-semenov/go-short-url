package app

import (
	"log"
	"net/http"
)

func RunApp() {
	repo := NewStorage()
	service := NewService(repo)
	h := EncodeURLHandler(service)
	http.HandleFunc("/", h.Handler)
	log.Println("starting server on 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
