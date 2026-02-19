package ui

import "github.com/charmbracelet/lipgloss"

var (
	UserColor      = lipgloss.Color("#5C6BC0")
	AssistantColor = lipgloss.Color("#7D56F4")
	BorderColor    = lipgloss.Color("#7D56F4")

	UserStyle = lipgloss.NewStyle().
			Foreground(UserColor).
			Bold(true).
			MarginLeft(2)

	AssistantStyle = lipgloss.NewStyle().
			Foreground(AssistantColor).
			Bold(true).
			MarginLeft(2)

	UserBubble = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(UserColor).
			Padding(0, 1).
			MarginLeft(2).
			MarginBottom(1)

	AssistantBubble = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(lipgloss.Color("#424242")).
			Padding(0, 1).
			MarginLeft(2).
			MarginBottom(1)

	LoadingStyle = lipgloss.NewStyle().
			Foreground(AssistantColor).
			Blink(true).
			MarginLeft(2)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	ChatBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)
)
