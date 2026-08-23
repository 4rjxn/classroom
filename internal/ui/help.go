package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderHelpModal(width int, height int) string {
	helpSections := []struct {
		title string
		keys  [][2]string
	}{
		{
			title: "Navigation",
			keys: [][2]string{
				{"↑ / k", "Move cursor up"},
				{"↓ / j", "Move cursor down"},
				{"g / Home", "Jump to top"},
				{"G / End", "Jump to bottom"},
				{"Enter", "Select course / View details"},
				{"Esc", "Go back / Exit search / Close modal"},
			},
		},
		{
			title: "Course Tabs & Detail Pane",
			keys: [][2]string{
				{"Tab / → / l", "Next tab"},
				{"Shift+Tab / ← / h", "Previous tab"},
				{"1 / 2 / 3 / 4", "Jump to Assignments / Materials / Announcements / Info"},
				{"d / u", "Scroll detail pane down / up"},
			},
		},
		{
			title: "Actions & Web Links",
			keys: [][2]string{
				{"o", "Open item / course in Google Classroom web"},
				{"a / Enter", "Open attachments or open attachment picker"},
				{"/ ", "Search and filter courses (in Courses view)"},
				{"r / Ctrl+r", "Refresh data from Google Classroom"},
				{"? ", "Toggle this help modal"},
				{"q / Ctrl+c", "Quit Classroom CLI"},
			},
		},
	}

	var content strings.Builder
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Background(colorPrimaryDim).
		Padding(0, 2).
		Render("⌨  Classroom CLI Keyboard Shortcuts")

	content.WriteString(title + "\n\n")

	for _, sec := range helpSections {
		secTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Render("▸ " + sec.title)
		content.WriteString(secTitle + "\n")

		for _, k := range sec.keys {
			key := helpKeyStyle.Width(18).Render("  " + k[0])
			desc := helpDescStyle.Render(k[1])
			content.WriteString(key + " " + desc + "\n")
		}
		content.WriteString("\n")
	}

	hint := lipgloss.NewStyle().
		Foreground(colorSubtle).
		Italic(true).
		Render("Press Esc or ? to return")
	content.WriteString(hint)

	box := modalBoxStyle.Render(content.String())

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}
