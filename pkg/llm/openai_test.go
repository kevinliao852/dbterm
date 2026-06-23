package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateSQL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/responses" {
			t.Fatalf("expected /responses, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("missing bearer token")
		}

		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request["model"] != "test-model" {
			t.Fatalf("expected test model, got %v", request["model"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"output": [{
				"content": [{
					"type": "output_text",
					"text": "{\"sql\":\"SELECT COUNT(*) FROM users\",\"explanation\":\"Counts users.\"}"
				}]
			}]
		}`))
	}))
	defer server.Close()

	client := Client{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		Model:      "test-model",
		HTTPClient: server.Client(),
	}
	suggestion, err := client.GenerateSQL(
		context.Background(),
		"sqlite",
		"CREATE TABLE users (id INTEGER);",
		"How many users are there?",
	)
	if err != nil {
		t.Fatalf("generate SQL: %v", err)
	}
	if suggestion.SQL != "SELECT COUNT(*) FROM users" {
		t.Fatalf("unexpected SQL: %q", suggestion.SQL)
	}
}

func TestGenerateSQLRequiresAPIKey(t *testing.T) {
	_, err := (Client{}).GenerateSQL(context.Background(), "sqlite", "", "question")
	if err == nil {
		t.Fatal("expected missing API key error")
	}
}
