package ui

import "github.com/charmbracelet/lipgloss"

// Use ANSI 0-15 colors to inherit terminal theme
var (
	// Primary accents using terminal's blue/magenta
	Primary   = lipgloss.Color("4") // Maps to theme's blue
	Secondary = lipgloss.Color("5") // Maps to theme's magenta/purple
	Accent    = lipgloss.Color("6") // Maps to theme's cyan/teal

	Red       = lipgloss.Color("1") // Maps to theme's red
	Turquoise = lipgloss.Color("2") // Maps to theme's turquoise
	Yellow    = lipgloss.Color("3") // Maps to theme's yellow

	// Text colors using terminal's foreground
	Text       = lipgloss.Color("7")  // Normal white/fg
	TextDim    = lipgloss.Color("8")  // Bright black (gray)
	TextBright = lipgloss.Color("15") // Bright white (bold)

	// Chat bubbles - use ANSI colors for backgrounds
	UserStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true).
			MarginLeft(2)

	UserBubble = lipgloss.NewStyle().
			Foreground(TextBright).
			Background(lipgloss.Color("0")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(0, 1).
			MarginLeft(2).
			MarginBottom(1)

	AssistantStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true).
			MarginLeft(2)

	AssistantBubble = lipgloss.NewStyle().
			Foreground(Text).
			Background(lipgloss.Color("0")). // Uses terminal background (black/dark)
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(0, 1).
			MarginLeft(2).
			MarginBottom(1)

	LoadingStyle = lipgloss.NewStyle().
			Foreground(Accent). // Cyan from theme
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Blink(true).
			MarginLeft(2)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(TextDim) // Gray from theme

	ChatBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary). // Blue from theme
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Align(lipgloss.Center).
			Padding(1, 0)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(TextBright).
				Background(Primary). // Blue highlight
				Padding(0, 1, 0, 0).
				Bold(true)

	ItemStyle = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 1, 0, 0)

	CursorStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)
)
