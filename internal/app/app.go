package app

import (
	"fmt"
	conf "github.com/da-semenov/go-short-url/internal/app/config"
	midlwr "github.com/da-semenov/go-short-url/internal/app/custommiddleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func RunApp() {
	config := conf.NewConfig()
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := NewStorage(config.FileStorage)
	if err != nil {
		fmt.Println("can't init repository", err)
		return
	}
	service := NewService(repo, config.BaseURL)
	h := EncodeURLHandler(service)
	router := chi.NewRouter()
	router.Use(middleware.CleanPath)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(midlwr.Compress)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.GetMethodHandler)
		r.Post("/api/shorten", h.PostShortenHandler)
		r.Post("/", h.PostMethodHandler)
		r.Put("/", h.DefaultHandler)
		r.Patch("/", h.DefaultHandler)
		r.Delete("/", h.DefaultHandler)
	})

	log.Println("starting server on 8080...")
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}
