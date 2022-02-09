package app

import (
	"log"
	"net/http"
)

func RunApp() {
	h := ShortURLHandler()
	http.HandleFunc("/", h.Handler)
	log.Println("Server started...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
