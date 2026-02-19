package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"teallama/ollama" // moduleName/packagePath
	"teallama/ui"
)

func main() {
	// 1. Create dependencies (the "backend")
	client := ollama.New(
		"http://localhost:11434/api/generate",
		"deepseek-r1:14b",
	)

	// 2. Create UI with injected dependencies
	initialModel := ui.New(client)

	// 3. Run Bubble Tea
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
