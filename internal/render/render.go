package render

import (
	"html/template"
	"net/http"
)

// renderTemplate renderiza os templates HTML
func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {

	tmplPath := "templates/" + tmpl + ".page.html"
	t, err := template.ParseFiles("templates/base.html", tmplPath)
	if err != nil {
		http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
