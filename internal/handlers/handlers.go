package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/justinas/nosurf"
	"github.com/mraramalho/bloGo/internal/config"
	"github.com/mraramalho/bloGo/internal/render"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v2"
)

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
	if err := render.RenderTemplate(w, "index", map[string]string{"Title": "Home"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	if err := render.RenderTemplate(w, "services", map[string]string{"Title": "Serviços"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *Repository) ContactHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["Title"] = "Contato"

	// Verificar se há um parâmetro de sucesso na URL
	if r.URL.Query().Get("success") == "true" {
		data["Success"] = true
	}

	// Adicionar token CSRF ao contexto
	data["CSRFToken"] = nosurf.Token(r)

	if r.Method == http.MethodGet {
		if err := render.RenderTemplate(w, "contact", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		// Validar os dados do formulário
		if name == "" || email == "" || message == "" {
			data["Error"] = "Por favor, preencha todos os campos"
			data["Name"] = name
			data["Email"] = email
			data["Message"] = message
			if err := render.RenderTemplate(w, "contact", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		// Lógica para enviar o email
		err := sendEmail(name, email, message)
		if err != nil {
			// Handle the error by showing an error message to the user
			data["Error"] = "Não foi possível enviar o email. Por favor, tente novamente mais tarde."
			data["Name"] = name
			data["Email"] = email
			data["Message"] = message
			if err = render.RenderTemplate(w, "contact", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		// Redirect on success
		http.Redirect(w, r, "/contact?success=true", http.StatusSeeOther)
	}
}

// BlogHandler handles the blog logic for parsing post data, feeds it into a
// Map for post data and renders blog cards.
func (m *Repository) BlogHandler(w http.ResponseWriter, r *http.Request) {

	files, err := filepath.Glob("posts/*.yaml")
	if err != nil {
		http.Error(w, "posts directory reading error", http.StatusInternalServerError)
		return
	}
	var posts []config.Post
	for _, file := range files {
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, "YAML directory reading error", http.StatusInternalServerError)
			return
		}

		// Crie uma nova instância para cada arquivo
		postData := &config.Post{}
		err = yaml.Unmarshal(yamlFile, postData)
		if err != nil {
			http.Error(w, "YAML parser error", http.StatusInternalServerError)
			return
		}

		file = filepath.ToSlash(file)
		slug := strings.TrimSuffix(filepath.Base(file), ".yaml")
		postData.Slug = slug

		content, err := convertMarkdownToHTML(postData.MDContent)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		postData.HTMLContent = template.HTML(content)

		// Pass data to app config PostDataMap
		m.App.PostDataMap[slug] = postData
		posts = append(posts, *postData)

	}
	if err := render.RenderTemplate(w, "blog", map[string]any{
		"Title": "Blog",
		"Posts": posts,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// func (m *Repository) ServiceHandler(w http.ResponseWriter, r *http.Request) {
// 	render.RenderTemplate(w, "services", map[string]string{"Title": "Serviços"})
// }

func (m *Repository) PostHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/posts/")
	postData := m.App.PostDataMap[slug]

	if err := render.RenderTemplate(w, "post", postData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ConvertMarkdownToHTML receives a string with markdown content and returns
// a string with HTML content
func convertMarkdownToHTML(mdContent string) (string, error) {
	// Converte para HTML
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(mdContent), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WebHookHandler handles the webhook logic for the GitHub repository.
func (m *Repository) WebHookHandler(w http.ResponseWriter, r *http.Request) {

	githubSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")

	if githubSecret == "" {
		http.Error(w, "invalid config", http.StatusInternalServerError)
		return
	}

	// Verifica se o método da requisição é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lê o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "body reading error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Restaura o corpo para o próximo uso (decodificação)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Valida a assinatura
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" || !validateSignature(githubSecret, body, signature) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	var payload struct {
		Ref string `json:"ref"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "JSON Decoding Error", http.StatusBadRequest)
		return
	}

	if payload.Ref == "refs/heads/main" {
		// Executa o comando git pull
		cmd := exec.Command("git", "pull")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("git pull error: %v\nOutput: %s", err, string(output))
			http.Error(w, "failed to pull", http.StatusInternalServerError)
			return
		}
		log.Println("Posts successfully updated.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Ignorado"))

}

func validateSignature(secret string, body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	expectedSig := "sha256=" + hex.EncodeToString(expectedMAC)
	return hmac.Equal([]byte(expectedSig), []byte(signature))
}
