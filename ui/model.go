package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss" // ADD THIS LINE
	"teallama/models"
	"teallama/ollama"
)

// Logo and other styles
const logo = `
   /\_/\
  ( o.o )
  > ^ <
 Teallama
`

var logoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#88C0D0")). // A nice blue-ish color
	Align(lipgloss.Center).
	Width(80) // Initial width, will be adjusted on WindowSizeMsg

type Model struct {
	viewport viewport.Model
	textarea textarea.Model
	client   *ollama.Client
	messages []models.Message
	width    int
	height   int
	loading  bool
	err      error
}

// Internal message types (lowercase = private)
type responseMsg string
type errMsg error

func New(client *ollama.Client) Model {

	ta := textarea.New()
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 2000

	// CRITICAL FIX: Disable line numbers (you see "1" in screenshot)
	ta.ShowLineNumbers = false
	ta.SetHeight(3)
	ta.SetWidth(80)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.KeyMap.InsertNewline.SetEnabled(false)

	// Setup viewport with default size (will adjust on WindowSizeMsg)
	vp := viewport.New(80, 20)
	vp.SetContent("Welcome to Teallama!\nType a message and press Enter.")

	return Model{
		textarea: ta,
		viewport: vp,
		client:   client,
		messages: []models.Message{},
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	// Update sub-components
	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Reserve space: viewport takes most, textarea takes 3 lines, rest for margins
		vpHeight := msg.Height - 8
		if vpHeight < 10 {
			vpHeight = 10
		}

		m.viewport.Width = msg.Width - 4
		m.viewport.Height = vpHeight
		m.textarea.SetWidth(msg.Width - 4)

		// Refresh content with new dimensions
		m.viewport.SetContent(m.renderMessages())
		logoStyle = logoStyle.Width(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.loading {
				return m, nil
			}

			content := m.textarea.Value()
			if content == "" {
				return m, nil
			}

			// Add user message
			m.messages = append(m.messages, models.Message{
				Content: content,
				IsUser:  true,
			})

			// CRITICAL: Update viewport immediately so user sees their message
			m.viewport.SetContent(m.renderMessages())
			m.viewport.GotoBottom()

			m.textarea.Reset()
			m.loading = true

			// Fetch response asynchronously
			return m, tea.Batch(
				tiCmd,
				m.fetchResponse(content),
			)
		}

	case responseMsg:
		// Assistant response received
		m.loading = false
		m.messages = append(m.messages, models.Message{
			Content: string(msg),
			IsUser:  false,
		})
		// CRITICAL: Refresh viewport with new message
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
		return m, nil

	case errMsg:
		m.loading = false
		m.err = msg
		m.messages = append(m.messages, models.Message{
			Content: fmt.Sprintf("Error: %v", msg),
			IsUser:  false,
		})
		m.viewport.SetContent(m.renderMessages())
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

// fetchResponse creates a tea.Cmd that calls Ollama
func (m Model) fetchResponse(prompt string) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.Generate(prompt, m.messages[:len(m.messages)-1])
		if err != nil {
			return errMsg(err)
		}
		return responseMsg(resp)
	}
}

// renderMessages builds the chat content string
func (m Model) renderMessages() string {
	if len(m.messages) == 0 {
		return "Welcome! Start chatting..."
	}

	var b strings.Builder

	for _, msg := range m.messages {
		if msg.IsUser {
			b.WriteString(UserStyle.Render("You") + "\n")
			b.WriteString(UserBubble.Render(msg.Content) + "\n\n")
		} else {
			b.WriteString(AssistantStyle.Render("Assistant") + "\n")
			// Word wrap long responses
			wrapped := wordWrap(msg.Content, m.viewport.Width-8)
			b.WriteString(AssistantBubble.Render(wrapped) + "\n\n")
		}
	}

	if m.loading {
		b.WriteString(LoadingStyle.Render("Thinking..."))
	}

	return b.String()
}

func (m Model) View() string {
	if m.err != nil && len(m.messages) == 0 {
		return fmt.Sprintf("Error: %v\n\nPress Ctrl+C to exit.", m.err)
	}

	// CRITICAL FIX: Actually render the viewport!
	// Previously you were ignoring m.viewport and just concatenating strings
	chatArea := ChatBorder.Render(m.viewport.View())

	helpText := SubtleStyle.Render("Enter to send • Ctrl+C to quit")
	status := SubtleStyle.Render(fmt.Sprintf("Model: deepseek-r1:14b | Messages: %d", len(m.messages)))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		logoStyle.Render(logo),
		chatArea,
		"",
		m.textarea.View(),
		helpText,
		status,
	)
}

// Helper to prevent text overflow
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
