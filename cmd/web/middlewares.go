package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// sessionLoad loads and saves the session on every request
func sessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

// noSurf adds a CSRF cookie in the response
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.ExemptGlob("/webhook")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}
