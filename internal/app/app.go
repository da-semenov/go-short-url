package app

import (
	"context"
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
	fileRepository, err := storage.NewFileStorage(config.FileStorage)
	if err != nil {
		fmt.Println("can't init file repository", err)
		return
	}
	postgresHandler, err := storage.NewPostgresHandler(context.Background(), config.DatabaseDSN)
	if err != nil {
		fmt.Println("can't init postgres handler", err)
		return
	}
	if config.ReInit {
		err = storage.ClearDatabase(postgresHandler)
		if err != nil {
			fmt.Println("can't clear database", err)
			return
		}
	}
	err = storage.InitDatabase(postgresHandler)
	if err != nil {
		fmt.Println("can't init database structure", err)
		return
	}
	postgresRepository, err := storage.NewPostgresRepository(postgresHandler)
	if err != nil {
		fmt.Println("can't init postgres repo", err)
		return
	}
	service := server.NewService(fileRepository, config.BaseURL)
	cs, _ := server.NewCryptoService()
	userService := server.NewUserService(postgresRepository)
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
		r.Get("/ping", uh.PingHandler)
		r.Post("/api/shorten", h.PostShortenHandler)
		r.Post("/", h.PostMethodHandler)
		r.Put("/", h.DefaultHandler)
		r.Patch("/", h.DefaultHandler)
		r.Delete("/", h.DefaultHandler)
	})

	log.Println("starting server on 8080...")
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}
