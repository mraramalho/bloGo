package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mraramalho/bloGo/internal/config"
	"github.com/mraramalho/bloGo/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(noSurf)
	router.Use(sessionLoad)

	router.Get("/", handlers.Repo.HomeHandler)
	router.Get("/about", handlers.Repo.AboutHandler)
	router.Get("/contact", handlers.Repo.ContactHandler)
	router.Post("/contact", handlers.Repo.ContactHandler)
	router.Get("/blog", handlers.Repo.BlogHandler)
	router.Get("/posts/{slug}", handlers.Repo.PostHandler)
	router.Get("/services", handlers.Repo.ServiceHandler)

	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	return router

}
