package config

import (
	"html/template"
	"log"
	"os"

	"github.com/alexedwards/scs/v2"
)

const (
	version      = "1.0.0"
	cssVersion   = "1"
	port         = ":8888"
	inProduction = false
	useCache     = true
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	Port          string
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Session       *scs.SessionManager
	InProduction  bool
	version       string
	cssVersion    string
}

func NewApp() *AppConfig {
	return &AppConfig{
		UseCache:      useCache,
		TemplateCache: make(map[string]*template.Template),
		InfoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:      log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Session:       scs.New(),
		InProduction:  inProduction,
		version:       version,
		cssVersion:    cssVersion,
		Port:          port,
	}
}
