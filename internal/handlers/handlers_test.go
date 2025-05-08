package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mraramalho/bloGo/internal/config"
)

var app = config.NewApp()

func TestWebhookHandler(t *testing.T) {

	repo := NewRepo(app)

	NewHandlers(repo)

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
