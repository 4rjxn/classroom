package ui_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestLipglossBorderDimensions(t *testing.T) {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Width(20).
		Height(10)

	rendered := style.Render("hello")
	lines := strings.Split(rendered, "\n")
	w := lipgloss.Width(rendered)
	h := lipgloss.Height(rendered)

	t.Logf("Width(20).Height(10) rendered dimensions: W=%d, H=%d (lines=%d)", w, h, len(lines))
}
