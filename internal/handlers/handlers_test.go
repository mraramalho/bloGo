package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/mraramalho/bloGo/internal/config"
)

var app = config.NewApp()
var repo *Repository

func TestMain(m *testing.M) {
	basePath, _ := os.Getwd()
	projectRoot := filepath.Join(basePath, filepath.ToSlash("../../"))

	// Muda o diretório de execução
	if err := os.Chdir(projectRoot); err != nil {
		log.Fatalf("Erro ao mudar diretório para o root do projeto: %v", err)
	}

	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("[Warning] Failed to load .env file: %v", err)
	}

	app.Session = scs.New()
	app.Session.Lifetime = 24 * time.Hour
	app.Session.Cookie.Secure = false

	repo = NewRepo(app)

	NewHandlers(repo)

	os.Exit(m.Run())
}

func TestWebhookHandler(t *testing.T) {

	payload := map[string]string{
		"ref": "refs/heads/main",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	secret := []byte(os.Getenv("GITHUB_WEBHOOK_SECRET"))
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hub-Signature-256", signature)

	if !validateSignature(os.Getenv("GITHUB_WEBHOOK_SECRET"), body, signature) {
		t.Error("erro na validação")
	}

	rr := httptest.NewRecorder()
	repo.WebHookHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("esperava status 200 OK, recebeu %d", rr.Code)
	}

}

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := app.Session.Load(req.Context(), "")
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	repo.HomeHandler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
