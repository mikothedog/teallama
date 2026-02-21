package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewSelect() string {
	var b strings.Builder

	b.WriteString(HeaderStyle.Render(Logo) + "\n")
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
	// logo := "TEALLAMA"
	header := HeaderStyle.Render(Logo)
	indicator := SubtleStyle.Render(fmt.Sprintf("Model: %s", m.client.Model))
	chatArea := ChatBorder.Render(m.viewport.View())
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
