package main

import (
	"fmt"
	"os"

	"github.com/4rjxn/classroom/internal/auth"
	"github.com/4rjxn/classroom/internal/ui"
	"github.com/4rjxn/classroom/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config, loadedPath, err := utils.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	_ = loadedPath // successfully discovered config

	token, err := auth.OfflineGeneration(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Authentication error: %v\n", err)
		os.Exit(1)
	}

	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: received empty authorization token.")
		os.Exit(1)
	}

	appModel := ui.NewUiStateModel(token, config)
	p := tea.NewProgram(appModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}
}
