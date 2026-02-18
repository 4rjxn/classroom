package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/classroom-cli/internal/domain"
	"github.com/classroom-cli/internal/models"
)

type ViewState int

const (
	menueView ViewState = iota
	courseView
	materialsView
)

type UiStateModel struct {
	State     ViewState
	Token     string
	courses   []models.CourseModel
	materials []models.CourseWorkMaterial
	cursor    int
}

func (m UiStateModel) Init() tea.Cmd {
	return nil
}

func (m UiStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "j":
			if m.cursor < len(m.courses)-1 {
				m.cursor++
			}
		case "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "esc":
			if m.State > 0 {
				m.State -= 1
			}
		case "enter":
			switch m.State {
			case menueView:
				m.State = courseView
				m.courses = domain.ListCourses(m.Token).Courses
			case courseView:
				m.State = materialsView
				m.materials = domain.ListMaterialsInCourse(m.Token, m.courses[m.cursor].Id).Materials
			}
		}
	}
	return m, nil

}

func (m UiStateModel) View() string {
	switch m.State {
	case menueView:
		return ">>>"
	case courseView:
		var courseList strings.Builder
		for index, course := range m.courses {
			courseList.WriteString("  ")
			if index == m.cursor {
				courseList.WriteString("* ")
			}
			courseList.WriteString(course.Name)
			courseList.WriteByte('\n')
		}
		return courseList.String()
	case materialsView:
		var materialList strings.Builder
		for index, materials := range m.materials {
			materialList.WriteString("  ")
			if index == m.cursor {
				materialList.WriteString("* ")
			}
			materialList.WriteString(materials.Title)
			materialList.WriteByte('\n')
		}
		return materialList.String()

	}
	return ""
}
