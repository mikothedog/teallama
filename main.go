package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"teallama/ollama"
	"teallama/ui"
)

func main() {
	// Create client without model (will be selected in UI)
	client := ollama.New("http://localhost:11434")

	initialModel := ui.New(client)

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
