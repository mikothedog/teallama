package ui

import "github.com/charmbracelet/lipgloss"

// Use ANSI 0-15 colors to inherit terminal theme
var (
	Color1 = lipgloss.Color("1")
	Color2 = lipgloss.Color("2")
	Color3 = lipgloss.Color("3")

	Primary   = lipgloss.Color("4")
	Secondary = lipgloss.Color("5")
	Accent    = lipgloss.Color("6")

	Text       = lipgloss.Color("7")
	TextDim    = lipgloss.Color("8")
	TextBright = lipgloss.Color("15")

	UserStyle = lipgloss.NewStyle().
			Foreground(Color1).
			Bold(true).
			MarginLeft(2)

	UserBubble = lipgloss.NewStyle().
			Foreground(TextBright).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(0, 1).
			MarginLeft(2).
			MarginBottom(1)

	AssistantStyle = lipgloss.NewStyle().
			Foreground(Color3).
			Bold(true).
			MarginLeft(2)

	AssistantBubble = lipgloss.NewStyle().
			Foreground(Text).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(0, 1).
			MarginLeft(2).
			MarginBottom(1)

	LoadingStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Blink(true).
			MarginLeft(2)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(TextDim)

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
				Foreground(TextBright).
				Background(Primary).
				Padding(0, 1, 0, 0).
				Bold(true)

	ItemStyle = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 1, 0, 0)

	CursorStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	TextAreaFocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	TextAreaBlurred = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(TextDim).
			Padding(0, 1).
			Faint(true)

	TextAreaText = lipgloss.NewStyle().
			Background(lipgloss.NoColor{})

	TextAreaCursorLine = lipgloss.NewStyle().
				Background(lipgloss.NoColor{})
)
