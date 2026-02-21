package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/classroom-cli/internal/auth"
	"github.com/classroom-cli/internal/models"
	"github.com/classroom-cli/internal/ui"

	"github.com/BurntSushi/toml"
)

func ReadConfig(config *models.Config) {
	fileName := "./config.toml"
	_, err := os.Stat(fileName)
	if err != nil {
		panic("Config file croupted")
	}
	toml.DecodeFile(fileName, config)
}

func main() {
	var config models.Config
	ReadConfig(&config)
	token := auth.OffileGeneration(config)
	p := tea.NewProgram(ui.UiStateModel{Token: token, State: 0})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
