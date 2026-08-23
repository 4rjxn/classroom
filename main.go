package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/classroom-cli/internal/auth"
	"github.com/classroom-cli/internal/ui"
	"github.com/classroom-cli/internal/utils"
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
