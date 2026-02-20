package ollama

import (
	"bytes"
	"encoding/json"
	"net/http"

	"teallama/models"
)

type Client struct {
	URL   string
	Model string // Selected model
	http  *http.Client
}

func New(url string) *Client {
	return &Client{
		URL:  url,
		http: &http.Client{Timeout: 0}, // No timeout for streaming
	}
}

func (c *Client) SetModel(model string) {
	c.Model = model
}

// ListModels fetches available models from Ollama
func (c *Client) ListModels() ([]string, error) {
	resp, err := c.http.Get(c.URL + "/api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	names := make([]string, len(result.Models))
	for i, m := range result.Models {
		names[i] = m.Name
	}
	return names, nil
}

// StreamRequest initiates a streaming generation and returns the response body
func (c *Client) StreamRequest(prompt string, history []models.Message) (*http.Response, error) {
	type request struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}

	// Build context from history
	var context string
	for _, msg := range history {
		role := "Assistant"
		if msg.IsUser {
			role = "User"
		}
		context += role + ": " + msg.Content + "\n"
	}
	context += "User: " + prompt + "\nAssistant: "

	body, _ := json.Marshal(request{
		Model:  c.Model,
		Prompt: context,
		Stream: true,
	})

	resp, err := c.http.Post(c.URL+"/api/generate", "application/json", bytes.NewReader(body))
	return resp, err
}
