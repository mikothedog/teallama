package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mikothedog/teallama/internal/commands"
	"github.com/mikothedog/teallama/internal/models"
	"github.com/mikothedog/teallama/internal/ollama"
)

// Screen states
type screen int

const (
	screenSelect screen = iota
	screenChat
)

// Model implements tea.Model
type Model struct {
	screen          screen
	availableModels []string
	cursor          int
	client          *ollama.Client
	viewport        viewport.Model
	textarea        textarea.Model
	messages        []models.Message
	width           int
	height          int
	loading         bool
	stream          io.ReadCloser
}

func New(client *ollama.Client) Model {
	ta := textarea.New()
	ta.Placeholder = "What are you pondering?"
	ta.ShowLineNumbers = false
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.KeyMap.InsertNewline.SetEnabled(false)
	ta.Focus()

	vp := viewport.New(80, 20)

	return Model{
		screen:   screenSelect,
		client:   client,
		textarea: ta,
		viewport: vp,
		messages: []models.Message{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		commands.LoadModels(m.client), // Call from commands package
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	if m.screen == screenChat {
		m.textarea, tiCmd = m.textarea.Update(msg)
	}
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.screen == screenChat {
			headerHeight := 6
			inputHeight := 4
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - headerHeight - inputHeight - 4
			m.textarea.SetWidth(msg.Width - 4)
			m.viewport.SetContent(m.renderMessages())
		}
		return m, nil

	// Handle messages from commands package
	case commands.ModelsLoadedMsg:
		m.availableModels = msg.Models
		return m, nil

	case commands.ErrMsg:
		m.loading = false
		if m.stream != nil {
			m.stream.Close()
			m.stream = nil
		}
		if m.screen == screenChat {
			m.messages = append(m.messages, models.Message{
				Content: fmt.Sprintf("Error: %v", msg.Err),
				IsUser:  false,
			})
			m.viewport.SetContent(m.renderMessages())
		} else {
			m.availableModels = []string{"Error loading models"}
		}
		return m, nil

	case commands.StreamStartMsg:
		m.stream = msg.Reader
		m.messages = append(m.messages, models.Message{
			Content: "",
			IsUser:  false,
		})
		return m, commands.ReadStreamChunk(m.stream)

	case commands.StreamChunkMsg:
		if len(m.messages) > 0 {
			lastIdx := len(m.messages) - 1
			m.messages[lastIdx].Content += msg.Content
			m.viewport.SetContent(m.renderMessages())
			m.viewport.GotoBottom()
		}
		return m, commands.ReadStreamChunk(m.stream)

	case commands.StreamDoneMsg:
		m.loading = false
		m.stream = nil
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		if m.screen == screenSelect {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.availableModels)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.availableModels) > 0 && m.cursor < len(m.availableModels) {
					selected := strings.TrimSpace(m.availableModels[m.cursor])
					m.client.SetModel(selected)
					m.screen = screenChat
					m.textarea.Reset()
					m.textarea.Focus()
					return m, tea.Batch(
						tea.WindowSize(),
						textarea.Blink,
					)
				}
			}
			return m, nil
		}

		if m.screen == screenChat && msg.Type == tea.KeyEnter && !m.loading {
			content := m.textarea.Value()
			if content == "" {
				return m, nil
			}

			m.messages = append(m.messages, models.Message{
				Content: content,
				IsUser:  true,
			})

			m.viewport.SetContent(m.renderMessages())
			m.viewport.GotoBottom()
			m.textarea.Reset()
			m.loading = true

			history := m.messages[:len(m.messages)-1]

			return m, tea.Batch(
				tiCmd,
				commands.StreamRequest(m.client, content, history), // Call from commands
			)
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	if m.screen == screenSelect {
		return m.viewSelect()
	}
	return m.viewChat()
}
