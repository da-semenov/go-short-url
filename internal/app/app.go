package app

import (
	"context"
	"fmt"
	conf "github.com/da-semenov/go-short-url/internal/app/config"
	"github.com/da-semenov/go-short-url/internal/app/handlers"
	midlwr "github.com/da-semenov/go-short-url/internal/app/middleware"
	serv "github.com/da-semenov/go-short-url/internal/app/server"
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
		err = storage.ClearDatabase(context.Background(), postgresHandler)
		if err != nil {
			fmt.Println("can't clear database structure", err)
			return
		}
	}
	err = storage.InitDatabase(context.Background(), postgresHandler)
	if err != nil {
		fmt.Println("can't init database structure", err)
		return
	}

	postgresRepository, err := storage.NewPostgresRepository(postgresHandler)
	if err != nil {
		fmt.Println("can't init postgres repository", err)
		return
	}

	deleteRepository, err := storage.NewDeleteRepository(postgresHandler)
	if err != nil {
		fmt.Println("can't init delete repository", err)
		return
	}

	cryptoService, err := serv.NewCryptoService()
	if err != nil {
		fmt.Println("error in crypto-service", err)
		return
	}
	userService := serv.NewUserService(postgresRepository, fileRepository, config.BaseURL)
	deleteService := serv.NewDeleteService(deleteRepository)
	uh := handlers.NewUserHandler(userService, cryptoService, deleteService)
	router := chi.NewRouter()
	router.Use(middleware.CleanPath)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(midlwr.GzipHandle)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", uh.GetMethodHandler)
		r.Get("/api/user/urls", uh.GetUserURLsHandler)
		r.Get("/ping", uh.PingHandler)
		r.Post("/api/shorten", uh.PostShortenHandler)
		r.Post("/api/shorten/batch", uh.PostShortenBatchHandler)
		r.Delete("/api/user/urls", uh.AsyncDeleteHandler)
		r.Post("/", uh.PostMethodHandler)
		r.Put("/", uh.DefaultHandler)
		r.Patch("/", uh.DefaultHandler)
		r.Delete("/", uh.DefaultHandler)
	})

	log.Println("starting server on 8080...")
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}
