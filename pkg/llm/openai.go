package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
	defaultModel   = "gpt-5.5"
)

type SQLSuggestion struct {
	SQL         string `json:"sql"`
	Explanation string `json:"explanation"`
}

type Client struct {
	APIKey     string
	BaseURL    string
	Model      string
	HTTPClient *http.Client
}

func NewClientFromEnv() Client {
	baseURL := strings.TrimRight(os.Getenv("OPENAI_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = defaultModel
	}

	return Client{
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		BaseURL: baseURL,
		Model:   model,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c Client) GenerateSQL(ctx context.Context, dialect, schema, question string) (SQLSuggestion, error) {
	if c.APIKey == "" {
		return SQLSuggestion{}, errors.New("OPENAI_API_KEY is not set")
	}

	body := map[string]any{
		"model": c.Model,
		"input": []map[string]string{
			{
				"role": "developer",
				"content": "Generate one read-only SQL query for the supplied database schema. " +
					"Never invent tables or columns. Do not generate INSERT, UPDATE, DELETE, DROP, ALTER, " +
					"TRUNCATE, CREATE, GRANT, or REVOKE statements.",
			},
			{
				"role": "user",
				"content": fmt.Sprintf(
					"SQL dialect: %s\n\nDatabase schema:\n%s\n\nQuestion:\n%s",
					dialect,
					schema,
					question,
				),
			},
		},
		"text": map[string]any{
			"verbosity": "low",
			"format": map[string]any{
				"type":   "json_schema",
				"name":   "sql_suggestion",
				"strict": true,
				"schema": map[string]any{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]any{
						"sql":         map[string]string{"type": "string"},
						"explanation": map[string]string{"type": "string"},
					},
					"required": []string{"sql", "explanation"},
				},
			},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return SQLSuggestion{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.BaseURL+"/responses",
		bytes.NewReader(payload),
	)
	if err != nil {
		return SQLSuggestion{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return SQLSuggestion{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return SQLSuggestion{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return SQLSuggestion{}, fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, apiError(responseBody))
	}

	var response struct {
		Output []struct {
			Content []struct {
				Type    string `json:"type"`
				Text    string `json:"text"`
				Refusal string `json:"refusal"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return SQLSuggestion{}, err
	}

	for _, output := range response.Output {
		for _, content := range output.Content {
			if content.Refusal != "" {
				return SQLSuggestion{}, errors.New(content.Refusal)
			}
			if content.Type != "output_text" || content.Text == "" {
				continue
			}

			var suggestion SQLSuggestion
			if err := json.Unmarshal([]byte(content.Text), &suggestion); err != nil {
				return SQLSuggestion{}, fmt.Errorf("decode structured response: %w", err)
			}
			suggestion.SQL = strings.TrimSpace(suggestion.SQL)
			if suggestion.SQL == "" {
				return SQLSuggestion{}, errors.New("model returned an empty SQL query")
			}
			return suggestion, nil
		}
	}

	return SQLSuggestion{}, errors.New("model response did not contain SQL")
}

func apiError(body []byte) string {
	var response struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &response) == nil && response.Error.Message != "" {
		return response.Error.Message
	}
	return strings.TrimSpace(string(body))
}
