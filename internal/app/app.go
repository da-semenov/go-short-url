package app

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
}

func RunApp() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ServerAddress %s\nBaseURL %s\n", cfg.ServerAddress, cfg.BaseURL)

	repo := NewStorage()
	service := NewService(repo)
	h := EncodeURLHandler(service)
	router := chi.NewRouter()
	router.Use(middleware.CleanPath)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.GetMethodHandler)
		r.Post("/api/shorten", h.PostShortenHandler)
		r.Post("/", h.PostMethodHandler)
		r.Put("/", h.DefaultHandler)
		r.Patch("/", h.DefaultHandler)
		r.Delete("/", h.DefaultHandler)
	})

	log.Println("starting server on 8080...")
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}
