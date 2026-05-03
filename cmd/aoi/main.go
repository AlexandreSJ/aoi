package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/AlexandreSJ/aoi/internal/ui"
)

func main() {
	p := tea.NewProgram(
		ui.NewApp(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}
