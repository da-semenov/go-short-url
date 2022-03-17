package app

import (
	"fmt"
	conf "github.com/da-semenov/go-short-url/internal/app/config"
	"github.com/da-semenov/go-short-url/internal/app/handlers"
	midlwr "github.com/da-semenov/go-short-url/internal/app/middleware"
	"github.com/da-semenov/go-short-url/internal/app/server"
	"github.com/da-semenov/go-short-url/internal/app/storage"
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

	repo, err := storage.NewStorage(config.FileStorage)
	if err != nil {
		fmt.Println("can't init repository", err)
		return
	}
	service := server.NewService(repo, config.BaseURL)
	cs, _ := server.NewCryptoService()
	userService := server.NewUserService(repo)
	h := handlers.EncodeURLHandler(service)
	uh := handlers.NewUserHandler(userService, cs)
	router := chi.NewRouter()
	router.Use(middleware.CleanPath)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(midlwr.GzipHandle)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.GetMethodHandler)
		r.Get("/api/user/urls", uh.GetUserURLsHandler)
		r.Post("/api/shorten", h.PostShortenHandler)
		r.Post("/", h.PostMethodHandler)
		r.Put("/", h.DefaultHandler)
		r.Patch("/", h.DefaultHandler)
		r.Delete("/", h.DefaultHandler)
	})

	log.Println("starting server on 8080...")
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}
