package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func RunApp() {
	repo := NewStorage()
	service := NewService(repo)
	h := EncodeURLHandler(service)
	router := chi.NewRouter()
	router.Use(middleware.CleanPath)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.GetMethodHandler)
		r.Post("/", h.PostMethodHandler)
		r.Put("/", h.DefaultHandler)
		r.Patch("/", h.DefaultHandler)
		r.Delete("/", h.DefaultHandler)
	})

	log.Println("starting server on 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
