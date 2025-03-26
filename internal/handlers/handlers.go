package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mraramalho/bloGo/internal/config"
	"github.com/mraramalho/bloGo/internal/render"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v2"
)

// Post representa um artigo do blog
type Post struct {
	Title   string
	Date    string
	Content template.HTML
}

// PostCard representa um card de um artigo do blog
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
	// files, err := filepath.Glob("posts/*.md")
	files, err := filepath.Glob("posts/*.yaml")
	if err != nil {
		http.Error(w, "Erro ao ler diretório de posts", http.StatusInternalServerError)
		return
	}

	var postCards []PostCard
	for _, file := range files {
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, "Erro ao ler diretório yaml com dados dos posts", http.StatusInternalServerError)
			return
		}

		// Crie uma nova instância para cada arquivo
		yamlPostData := &config.YAMLPostData{}
		err = yaml.Unmarshal(yamlFile, yamlPostData)
		if err != nil {
			http.Error(w, "Erro no parser do YAML", http.StatusInternalServerError)
			return
		}

		slug := strings.TrimPrefix(strings.TrimSuffix(file, ".yaml"), "posts\\")
		postCards = append(postCards, PostCard{
			Title:   yamlPostData.Title,
			Excerpt: yamlPostData.Excerpt,
			Slug:    slug,
		})

		// Pass data to app config YAMLPostDataMap
		m.App.YAMLPostDataMap[slug] = yamlPostData
	}
	render.RenderTemplate(w, "blog", map[string]any{
		"Title":     "Blog",
		"PostCards": postCards,
	})
}
func (m *Repository) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "services", map[string]string{"Title": "Serviços"})
}

func (m *Repository) PostHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/posts/")
	// filePath := "posts/" + slug + ".md"
	postData := m.App.YAMLPostDataMap[slug]
	mdContent := postData.Content

	content, err := convertMarkdownToHTML(mdContent)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	post := Post{
		Title:   postData.Title,
		Date:    postData.Created,
		Content: template.HTML(content),
	}

	render.RenderTemplate(w, "post", post)
}

// ConvertMarkdownToHTML recebe uma string em formato .md e converte para HTML
func convertMarkdownToHTML(mdContent string) (string, error) {
	// Converte para HTML
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(mdContent), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
