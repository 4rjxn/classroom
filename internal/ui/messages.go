package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/4rjxn/classroom/internal/domain"
	"github.com/4rjxn/classroom/internal/models"
	"github.com/pkg/browser"
)

// Messages
type coursesLoadedMsg struct {
	Courses []models.CourseModel
	Err     error
}

type courseDetailsLoadedMsg struct {
	CourseID      string
	CourseWork    []models.CourseWorkModel
	Materials     []models.CourseWorkMaterial
	Announcements []models.CourseAnnouncement
	Err           error
}

type urlOpenedMsg struct {
	URL string
	Err error
}

type statusMsg struct {
	Text    string
	IsError bool
}

type clearStatusMsg struct{}

// Commands
func fetchCoursesCmd(token string) tea.Cmd {
	return func() tea.Msg {
		resp, err := domain.ListCourses(token)
		if err != nil {
			return coursesLoadedMsg{Err: err}
		}
		return coursesLoadedMsg{Courses: resp.Courses}
	}
}

func fetchCourseDataCmd(token string, courseID string) tea.Cmd {
	return func() tea.Msg {
		type cwRes struct {
			cw  []models.CourseWorkModel
			err error
		}
		type matRes struct {
			mat []models.CourseWorkMaterial
			err error
		}
		type annRes struct {
			ann []models.CourseAnnouncement
			err error
		}

		cwChan := make(chan cwRes, 1)
		matChan := make(chan matRes, 1)
		annChan := make(chan annRes, 1)

		// Fetch coursework, materials, and announcements concurrently
		go func() {
			resp, err := domain.ListCourseWorkInCourse(token, courseID)
			cwChan <- cwRes{cw: resp.CourseWork, err: err}
		}()

		go func() {
			resp, err := domain.ListMaterialsInCourse(token, courseID)
			matChan <- matRes{mat: resp.Materials, err: err}
		}()

		go func() {
			resp, err := domain.ListAnnouncementsInCourse(token, courseID)
			annChan <- annRes{ann: resp.Announcements, err: err}
		}()

		r1 := <-cwChan
		r2 := <-matChan
		r3 := <-annChan

		var firstErr error
		if r1.err != nil {
			firstErr = r1.err
		} else if r2.err != nil {
			firstErr = r2.err
		} else if r3.err != nil {
			firstErr = r3.err
		}

		return courseDetailsLoadedMsg{
			CourseID:      courseID,
			CourseWork:    r1.cw,
			Materials:     r2.mat,
			Announcements: r3.ann,
			Err:           firstErr,
		}
	}
}

func openBrowserCmd(targetURL string) tea.Cmd {
	return func() tea.Msg {
		if targetURL == "" {
			return urlOpenedMsg{Err: fmt.Errorf("no URL provided to open")}
		}
		err := browser.OpenURL(targetURL)
		return urlOpenedMsg{URL: targetURL, Err: err}
	}
}

func flashStatusCmd(text string, isError bool) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return statusMsg{Text: text, IsError: isError}
		},
		tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
			return clearStatusMsg{}
		}),
	)
}
