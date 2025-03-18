package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
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

func main() {

	// Servir arquivos estáticos (CSS, JS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Página inicial
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index", map[string]string{"Title": "Home"})
	})

	// Pagina de Serviços
	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "services", map[string]string{"Title": "Serviços"})
	})

	// Página de Contato
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			renderTemplate(w, "contact", map[string]string{"Title": "Contato"})
		}

		if r.Method == http.MethodPost {
			name := r.FormValue("name")
			email := r.FormValue("email")
			message := r.FormValue("message")
			// Lógica para enviar o email
			err := sendEmail(name, email, message)
			if err != nil {
				// Handle the error by showing an error message to the user
				renderTemplate(w, "contact", map[string]interface{}{
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
	})

	// Página de Blog
	http.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
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
		renderTemplate(w, "blog", map[string]any{
			"Title": "Blog",
			"Posts": posts,
		})
	})

	// Página de um Post Específico
	http.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
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

		renderTemplate(w, "post", post)
	})

	// Inicia o servidor
	log.Println("Servidor rodando em http://localhost:8888")
	http.ListenAndServe(":8888", nil)
}

// renderTemplate renderiza os templates HTML
func renderTemplate(w http.ResponseWriter, tmpl string, data any) {

	tmplPath := "templates/" + tmpl + ".page.html"
	t, err := template.ParseFiles("templates/base.html", tmplPath)
	if err != nil {
		http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func sendEmail(name, email, message string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")

	// Create a new message
	msg := gomail.NewMessage()

	// Set email headers
	msg.SetHeader("From", email)
	msg.SetHeader("To", "aramalho.1991@gmail.com")
	msg.SetHeader("Subject", "Vi seu site e gostaria de entrar em contato")

	// Set email body
	msg.SetBody("text/plain", message)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, user, password)

	// Send the email
	if err := dialer.DialAndSend(msg); err != nil {
		log.Println("Error sending email:", err)
		return err
	}
	log.Println("Email sent successfully!")
	return nil
}
