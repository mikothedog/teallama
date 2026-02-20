package commands

import (
	"bufio"
	"encoding/json"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikothedog/teallama/internal/models"
	"github.com/mikothedog/teallama/internal/ollama"
)

// Message types returned by commands
type (
	ModelsLoadedMsg struct{ Models []string }
	StreamStartMsg  struct{ Reader io.ReadCloser }
	StreamChunkMsg  struct{ Content string }
	StreamDoneMsg   struct{}
	ErrMsg          struct{ Err error }
)

// LoadModels fetches available models from Ollama
func LoadModels(client *ollama.Client) tea.Cmd {
	return func() tea.Msg {
		models, err := client.ListModels()
		if err != nil {
			return ErrMsg{Err: err}
		}
		return ModelsLoadedMsg{Models: models}
	}
}

// THIS IS THE COMMAND - wraps the client method in tea.Cmd
func StreamRequest(client *ollama.Client, prompt string, history []models.Message) tea.Cmd {
	return func() tea.Msg {
		reader, err := client.StreamRequest(prompt, history) // Calls the method above
		if err != nil {
			return ErrMsg{Err: err}
		}
		return StreamStartMsg{Reader: reader}
	}
}

// ReadStreamChunk reads the next token from the stream
func ReadStreamChunk(stream io.ReadCloser) tea.Cmd {
	return func() tea.Msg {
		scanner := bufio.NewScanner(stream)

		if scanner.Scan() {
			line := scanner.Bytes()

			var resp struct {
				Response string `json:"response"`
				Done     bool   `json:"done"`
			}

			if err := json.Unmarshal(line, &resp); err != nil {
				return ReadStreamChunk(stream)()
			}

			if resp.Done {
				stream.Close()
				return StreamDoneMsg{}
			}

			return StreamChunkMsg{Content: resp.Response}
		}

		stream.Close()
		return StreamDoneMsg{}
	}
}
