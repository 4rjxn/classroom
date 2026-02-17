package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/classroom-cli/internal/auth"
	"github.com/classroom-cli/internal/domain"
	"github.com/classroom-cli/internal/models"

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
	token := auth.GenerateToken(config)
	var action string
	var param string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("== hello there ==\n")
	for {
		fmt.Print("Available actions:\n")
		fmt.Print(" list courses [list]\n")
		fmt.Print(" list materials [choose <id>]\n")
		fmt.Print(" quit [q]\n")
		fmt.Print(">> ")
		scanner.Scan()
		parts := strings.Fields(scanner.Text())
		if len(parts) > 0 {
			action = parts[0]
		}
		if len(parts) > 1 {
			param = parts[1]
		}
		if action == "q" {
			return
		}
		if action == "list" {
			fmt.Println("======================================")
			domain.ListCourses(token)
			fmt.Println("======================================")
		} else if action == "choose" && param != "" {
			fmt.Println("======================================")
			fmt.Println("----Materials----")
			domain.ListMaterialsInCourse(token, param)
			fmt.Println("======================================")
			fmt.Println("----Announcements----")
			domain.ListAnnouncementsInCourse(token, param)
			fmt.Println("======================================")
		}
	}
}
