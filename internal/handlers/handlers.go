package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/mraramalho/bloGo/internal/config"
	"github.com/mraramalho/bloGo/internal/render"
)

// Post representa um artigo do blog
type Post struct {
	Title   string
	Date    string
	Content template.HTML
}

type PostCard struct {
	Title   string
	Excerpt string
	Slug    string
}

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	m.App.Session.Put(r.Context(), "remore_ip", r.RemoteAddr)
	render.RenderTemplate(w, "index", map[string]string{"Title": "Home"})
}

func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	render.RenderTemplate(w, "services", map[string]string{"Title": "Serviços"})
}

func (m *Repository) ContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render.RenderTemplate(w, "contact", map[string]string{"Title": "Contato"})
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")
		// Lógica para enviar o email
		err := sendEmail(name, email, message)
		if err != nil {
			// Handle the error by showing an error message to the user
			render.RenderTemplate(w, "contact", map[string]interface{}{
				"Title":   "Contato",
				"Error":   "Não foi possível enviar o email. Por favor, tente novamente mais tarde.",
				"Name":    name,
				"Email":   email,
				"Message": message,
			})
			return
		}

		// Redirect on success
		http.Redirect(w, r, "/contact?success=true", http.StatusSeeOther)

	}
}

func (m *Repository) BlogHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obter os posts do blog
	files, err := filepath.Glob("posts/*.md")
	if err != nil {
		http.Error(w, "Erro ao ler diretório de posts", http.StatusInternalServerError)
		return
	}
	var posts []PostCard
	for _, file := range files {
		slug := strings.TrimPrefix(strings.TrimSuffix(file, ".md"), "posts\\")
		posts = append(posts, PostCard{
			Title:   strings.ToTitle(strings.ReplaceAll(slug, "-", " ")),
			Excerpt: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor.",
			Slug:    "posts/" + slug,
		})
	}

	render.RenderTemplate(w, "blog", map[string]any{
		"Title": "Blog",
		"Posts": posts,
	})
}

func (m *Repository) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "services", map[string]string{"Title": "Serviços"})
}

func (m *Repository) PostHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/posts/")
	filePath := "posts/" + slug + ".md"

	content, err := ConvertMarkdownToHTML(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	stat := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	cTime := time.Unix(0, stat.CreationTime.Nanoseconds())
	creationTime := cTime.Format("02 de janeiro de 2006")

	post := Post{
		Title:   strings.ToTitle(strings.ReplaceAll(slug, "-", " ")), // Título baseado no slug
		Date:    creationTime,
		Content: template.HTML(content), // Inserindo HTML seguro
	}

	render.RenderTemplate(w, "post", post)
}
