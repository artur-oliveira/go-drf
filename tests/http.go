package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

// RequestOptions define os parâmetros para MakeRequest
type RequestOptions struct {
	Method string
	URL    string
	Token  string      // Token JWT (Bearer)
	Body   interface{} // DTO para ser serializado em JSON
}

// MakeRequest é o helper central para todas as chamadas de API
// Ele agora recebe o *fiber.App como argumento.
func MakeRequest(t *testing.T, app *fiber.App, opts RequestOptions) (*http.Response, string) {
	var bodyReader io.Reader
	if opts.Body != nil {
		bodyBytes, err := json.Marshal(opts.Body)
		if err != nil {
			t.Fatalf("Falha ao serializar corpo da requisição: %v", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req := httptest.NewRequest(opts.Method, opts.URL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if opts.Token != "" {
		req.Header.Set("Authorization", "Bearer "+opts.Token)
	}

	// Usa o app.Test() para fazer a requisição em memória
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Falha ao executar requisição %s %s: %v", opts.Method, opts.URL, err)
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Falha ao ler corpo da resposta: %v", err)
	}

	return resp, string(respBodyBytes)
}
