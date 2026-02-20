package ollama

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mikothedog/teallama/internal/models"
)

type Client struct {
	URL   string
	Model string
	http  *http.Client
}

func New(url string) *Client {
	return &Client{
		URL:  url,
		http: &http.Client{},
	}
}

func (c *Client) SetModel(model string) {
	c.Model = model
}

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

// THIS IS THE METHOD - returns io.ReadCloser
func (c *Client) StreamRequest(prompt string, history []models.Message) (io.ReadCloser, error) {
	type request struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}

	context := models.HistoryToString(history)
	fullPrompt := context + "User: " + prompt + "\nAssistant: "

	body, _ := json.Marshal(request{
		Model:  c.Model,
		Prompt: fullPrompt,
		Stream: true,
	})

	resp, err := c.http.Post(c.URL+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return resp.Body, nil // Return the body (io.ReadCloser), not resp
}
