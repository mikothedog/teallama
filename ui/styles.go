package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Primary   = lipgloss.Color("#7D56F4")
	Secondary = lipgloss.Color("#5C6BC0")
	Text      = lipgloss.Color("#E0E0E0")
	Subtle    = lipgloss.Color("241")

	// Chat styles (existing)
	UserStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true).
			MarginLeft(2)

	UserBubble = lipgloss.NewStyle().
			Foreground(Text).
			Background(Secondary).
			Padding(0, 1).
			MarginLeft(2).
			MarginBottom(1)

	AssistantStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			MarginLeft(2)

	AssistantBubble = lipgloss.NewStyle().
			Foreground(Text).
			Background(lipgloss.Color("#424242")).
			Padding(0, 1).
			MarginLeft(2).
			MarginBottom(1)

	LoadingStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Blink(true).
			MarginLeft(2)

	SubtleStyle = lipgloss.NewStyle().Foreground(Subtle)

	ChatBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Align(lipgloss.Center).
			Padding(1, 0)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Text).
				Background(Primary).
				Padding(0, 1, 0, 0).
				Bold(true)

	ItemStyle = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 1, 0, 0)

	CursorStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)
)
