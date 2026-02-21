package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikothedog/teallama/internal/ollama"
	"github.com/mikothedog/teallama/internal/ui"
)

func main() {
	// Create Ollama client (model selected in UI)
	client := ollama.New("http://localhost:11434")

	// Create Bubble Tea program
	program := tea.NewProgram(
		ui.New(client),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
