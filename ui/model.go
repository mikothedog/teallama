package ui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"teallama/models"
	"teallama/ollama"
)

// Screen states
type screen int

const (
	screenSelect screen = iota
	screenChat
)

// Message types for Bubble Tea
type (
	modelsLoadedMsg []string
	streamStartMsg  struct{ reader io.ReadCloser }
	streamChunkMsg  struct{ content string }
	streamDoneMsg   struct{}
	errMsg          struct{ error }
)

type Model struct {
	// State
	screen          screen
	availableModels []string
	cursor          int
	client          *ollama.Client

	// Chat state
	viewport viewport.Model
	textarea textarea.Model
	messages []models.Message
	width    int
	height   int
	loading  bool
	stream   io.ReadCloser // Active stream connection
}

func New(client *ollama.Client) Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
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
		m.loadModels(), // Fetch available models on start
	)
}

// Command to fetch models from Ollama
func (m Model) loadModels() tea.Cmd {
	return func() tea.Msg {
		models, err := m.client.ListModels()
		if err != nil {
			return errMsg{err}
		}
		return modelsLoadedMsg(models)
	}
}

// Command to read next chunk from stream
func (m Model) readChunk() tea.Msg {
	if m.stream == nil {
		return streamDoneMsg{}
	}

	scanner := bufio.NewScanner(m.stream)
	if scanner.Scan() {
		line := scanner.Bytes()

		var resp struct {
			Response string `json:"response"`
			Done     bool   `json:"done"`
		}

		if err := json.Unmarshal(line, &resp); err != nil {
			return m.readChunk() // Skip malformed lines
		}

		if resp.Done {
			m.stream.Close()
			m.stream = nil
			return streamDoneMsg{}
		}

		return streamChunkMsg{content: resp.Response}
	}

	// EOF or error
	if m.stream != nil {
		m.stream.Close()
		m.stream = nil
	}
	return streamDoneMsg{}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	// Always update sub-components
	if m.screen == screenChat {
		m.textarea, tiCmd = m.textarea.Update(msg)
		m.viewport, vpCmd = m.viewport.Update(msg)
	} else {
		// On selection screen, just process viewport updates (for window resize)
		// but NOT textarea updates
		_, vpCmd = m.viewport.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate heights based on screen
		if m.screen == screenChat {
			headerHeight := 6 // Logo height
			inputHeight := 4
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - headerHeight - inputHeight - 4
			m.textarea.SetWidth(msg.Width - 4)
			m.viewport.SetContent(m.renderMessages())
		}
		return m, nil

	case modelsLoadedMsg:
		m.availableModels = msg
		return m, nil

	case errMsg:
		// Show error in current screen
		m.availableModels = []string{"Error: " + msg.Error(), "Press Ctrl+C to exit"}
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// Selection screen controls
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
				if len(m.availableModels) > 0 {
					selected := m.availableModels[m.cursor]
					m.client.SetModel(selected)
					m.screen = screenChat

					m.textarea.Reset()
					m.textarea.Focus()

					return m, tea.Batch(tea.WindowSize(), textarea.Blink)
				}
			}
			return m, nil
		}

		// Chat screen controls
		if m.screen == screenChat {
			if msg.Type == tea.KeyEnter && !m.loading {
				content := m.textarea.Value()
				if content == "" {
					return m, nil
				}

				// Add user message
				m.messages = append(m.messages, models.Message{
					Content: content,
					IsUser:  true,
				})

				m.viewport.SetContent(m.renderMessages())
				m.viewport.GotoBottom()
				m.textarea.Reset()
				m.loading = true

				// Start streaming
				return m, m.startStream(content)
			}
		}

	case streamStartMsg:
		m.stream = msg.reader
		// Add empty assistant message to be filled
		m.messages = append(m.messages, models.Message{
			Content: "",
			IsUser:  false,
		})
		return m, m.readChunk

	case streamChunkMsg:
		// Append to last message (assistant)
		lastIdx := len(m.messages) - 1
		m.messages[lastIdx].Content += msg.content
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
		// Continue reading next chunk
		return m, m.readChunk

	case streamDoneMsg:
		m.loading = false
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) startStream(prompt string) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.StreamRequest(prompt, m.messages[:len(m.messages)-1])
		if err != nil {
			return errMsg{err}
		}
		return streamStartMsg{reader: resp.Body}
	}
}

func (m Model) renderMessages() string {
	if len(m.messages) == 0 {
		return "Welcome! Start chatting..."
	}

	var b strings.Builder
	for _, msg := range m.messages {
		if msg.IsUser {
			b.WriteString(UserStyle.Render("You") + "\n")
			wrapped := wordWrap(msg.Content, m.viewport.Width-8)
			b.WriteString(UserBubble.Render(wrapped) + "\n\n")
		} else {
			b.WriteString(AssistantStyle.Render("Assistant") + "\n")
			if msg.Content == "" && m.loading {
				b.WriteString(LoadingStyle.Render("Thinking...") + "\n\n")
			} else {
				wrapped := wordWrap(msg.Content, m.viewport.Width-8)
				b.WriteString(AssistantBubble.Render(wrapped) + "\n\n")
			}
		}
	}
	return b.String()
}

func (m Model) View() string {
	if m.screen == screenSelect {
		return m.viewSelect()
	}
	return m.viewChat()
}

func (m Model) viewSelect() string {
	var b strings.Builder

	title := `
  ████████╗███████╗ █████╗ ██╗     ██╗      █████╗ ███╗   ███╗ █████╗ 
  ╚══██╔══╝██╔════╝██╔══██╗██║     ██║     ██╔══██╗████╗ ████║██╔══██╗
     ██║   █████╗  ███████║██║     ██║     ███████║██╔████╔██║███████║
     ██║   ██╔══╝  ██╔══██║██║     ██║     ██╔══██║██║╚██╔╝██║██╔══██║
     ██║   ███████╗██║  ██║███████╗███████╗██║  ██║██║ ╚═╝ ██║██║  ██║
     ╚═╝   ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝`

	b.WriteString(HeaderStyle.Render(title) + "\n")
	b.WriteString(SubtleStyle.Render("  Choose your AI model\n"))
	b.WriteString(SubtleStyle.Render("\n"))

	if len(m.availableModels) == 0 {
		b.WriteString("  Loading models...")
		return b.String()
	}

	for i, model := range m.availableModels {
		model = strings.TrimSpace(model)
		isSelected := m.cursor == i

		var prefix string
		if isSelected {
			prefix = "> " // Cursor + space
		} else {
			prefix = "  " // Two spaces to match cursor width
		}

		var styledModel string
		if isSelected {
			styledModel = SelectedItemStyle.Render(model)
		} else {
			styledModel = ItemStyle.Render(model)
		}

		// If selected, style the cursor too
		if isSelected {
			b.WriteString(CursorStyle.Render(">") + " " + styledModel + "\n")
		} else {
			b.WriteString(prefix + styledModel + "\n")
		}
	}

	b.WriteString("\n" + SubtleStyle.Render("  ↑/↓ or k/j to navigate • Enter to select • Ctrl+C to quit"))
	return b.String()
}

func (m Model) viewChat() string {
	// Header
	logo := "TEALLAMA"
	header := HeaderStyle.Render(logo)

	// Model indicator
	indicator := SubtleStyle.Render(fmt.Sprintf("Model: %s", m.client.Model))

	// Chat area
	chatArea := ChatBorder.Render(m.viewport.View())

	// Help text
	help := SubtleStyle.Render("Enter to send • Ctrl+C to quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		indicator,
		chatArea,
		"",
		m.textarea.View(),
		help,
	)
}

func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}
	var result strings.Builder
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		for len(line) > width {
			result.WriteString(line[:width] + "\n")
			line = line[width:]
		}
		result.WriteString(line + "\n")
	}
	return strings.TrimSuffix(result.String(), "\n")
}
