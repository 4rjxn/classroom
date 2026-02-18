package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/classroom-cli/internal/domain"
	"github.com/classroom-cli/internal/models"
)

type ViewState int

const (
	menueView ViewState = iota
	courseView
	courseDetailsView
)

type UiStateModel struct {
	State         ViewState
	Token         string
	courses       []models.CourseModel
	materials     []models.CourseWorkMaterial
	announcements []models.CourseAnnouncement
	cursor        int
	isMaterial    bool
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
				m.State = courseDetailsView
				m.materials = domain.ListMaterialsInCourse(m.Token, m.courses[m.cursor].Id).Materials
				m.announcements = domain.ListAnnouncementsInCourse(m.Token, m.courses[m.cursor].Id).Announcements
			}
		case "tab":
			switch m.State {
			case courseDetailsView:
				m.isMaterial = !m.isMaterial
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
		var courseList []string
		for _, course := range m.courses {
			courseList = append(courseList, course.Name+" "+course.Sub)
		}
		return list.New(courseList).Enumerator(func(items list.Items, index int) string {
			if index == m.cursor {
				return ">"
			}
			return ""
		}).String()

	case courseDetailsView:
		if m.isMaterial {
			parentList := list.New()
			for _, materials := range m.materials {
				parentList.Item(materials.Title + "\n  " + materials.Description)
			}
			parentList.Enumerator(func(items list.Items, index int) string {
				if index == m.cursor {
					return ">"
				}
				return ""
			})
			return parentList.String()
		} else {
			parentList := list.New()
			for _, announs := range m.announcements {
				parentList.Item(announs.Text)
			}
			parentList.Enumerator(func(items list.Items, index int) string {
				if index == m.cursor {
					return ">"
				}
				return ""
			})
			return parentList.String()

		}
	}
	return ""
}
