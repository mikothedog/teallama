package ollama

import (
	"bytes"
	"encoding/json"
	"net/http"
	"teallama/models" // Import your local module
	"time"
)

// Client talks to Ollama HTTP API.
type Client struct {
	URL   string
	Model string
	http  *http.Client
}

// New creates a client. Exported (capital N).
func New(url, model string) *Client {
	return &Client{
		URL:   url,
		Model: model,
		http:  &http.Client{Timeout: 120 * time.Second},
	}
}

// Generate sends prompt and returns response. Blocking call.
// Returns (string, error) - we'll wrap it in a tea.Cmd later.
func (c *Client) Generate(prompt string, history []models.Message) (string, error) {
	type request struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}

	fullPrompt := models.HistoryToString(history) + "User: " + prompt + "\nAssistant: "

	body, _ := json.Marshal(request{
		Model:  c.Model,
		Prompt: fullPrompt,
		Stream: false, // Simpler for now
	})

	resp, err := c.http.Post(c.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Response string `json:"response"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Response, nil
}
