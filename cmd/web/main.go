package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mraramalho/bloGo/internal/config"
	"github.com/mraramalho/bloGo/internal/handlers"
)

var app = config.NewApp()

func main() {

	app.Session.Lifetime = 24 * time.Hour
	app.Session.Cookie.Persist = true
	app.Session.Cookie.SameSite = http.SameSiteLaxMode
	app.Session.Cookie.Secure = app.InProduction
	repo := handlers.NewRepo(app)
	handlers.NewHandlers(repo)
	srv := &http.Server{
		Addr:    app.Port,
		Handler: routes(app),
	}

	log.Printf("Servidor rodando em http://localhost%s\n", app.Port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
