package render

import (
	"html/template"
	"log/slog"
	"net/http"
)

const (
	templateDir = "templates/"
	templateExt = ".page.html"
)

// renderTemplate renderiza os templates HTML
func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	tmplPath := templateDir + tmpl + templateExt
	t, err := template.ParseFiles("templates/base.html", tmplPath)
	if err != nil {
		http.Error(w, "Erro loading template", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		slog.Error("Error rendering template", "error", err)
		http.Error(w, "Erro rendering template", http.StatusInternalServerError)
		return
	}
}
