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

type RequestOptions struct {
	Method string
	URL    string
	Token  string
	Body   interface{}
}

func MakeRequest(t *testing.T, app *fiber.App, opts RequestOptions) (*http.Response, string) {
	var bodyReader io.Reader
	var bodyString string
	if opts.Body != nil {
		bodyBytes, err := json.Marshal(opts.Body)
		if err != nil {
			t.Fatalf("Falha ao serializar corpo da requisição: %v", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
		bodyString = string(bodyBytes)
	}

	req := httptest.NewRequest(opts.Method, opts.URL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if opts.Token != "" {
		req.Header.Set("Authorization", "Bearer "+opts.Token)
	}
	t.Logf("\nRequest\nPath: %s\nBody: %s", req.URL.String(), bodyString)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Falha ao executar requisição %s %s: %v", opts.Method, opts.URL, err)
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Falha ao ler corpo da resposta: %v", err)
	}

	var respBodyString = string(respBodyBytes)

	t.Logf("\nResponse\nStatus: %d\nPath: %s\nResponse: %s", resp.StatusCode, req.URL.String(), respBodyString)

	return resp, respBodyString
}
