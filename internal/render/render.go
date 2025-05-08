package render

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
)

const (
	templateExt = ".page.html"
)

var basePath = filepath.ToSlash("./templates")

// renderTemplate renderiza os templates HTML
func RenderTemplate(w http.ResponseWriter, tmpl string, data any) error {
	tmplPath := filepath.Join(basePath, tmpl+templateExt)
	baseLayout := filepath.Join(basePath, "base.html")

	t, err := template.ParseFiles(baseLayout, tmplPath)
	if err != nil {
		return fmt.Errorf("error loading template: %s", err)
	}

	if err := t.Execute(w, data); err != nil {
		slog.Error("error rendering template", "error", err)
		return fmt.Errorf("error rendering template: %s", err)
	}

	return nil
}
