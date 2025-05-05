package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/mraramalho/bloGo/internal/config"
	"github.com/mraramalho/bloGo/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(noSurf)
	router.Use(sessionLoad)

	router.Get("/", handlers.Repo.HomeHandler)
	router.Get("/about", handlers.Repo.AboutHandler)
	router.Get("/contact", handlers.Repo.ContactHandler)
	router.Post("/contact", handlers.Repo.ContactHandler)
	router.Get("/blog", handlers.Repo.BlogHandler)
	router.Get("/posts/{slug}", handlers.Repo.PostHandler)
	// router.Get("/services", handlers.Repo.ServiceHandler)

	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	return router

}
