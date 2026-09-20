package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Colors of the active theme, rebuilt by buildStyles.
var (
	colorPrimary    lipgloss.Color
	colorPrimaryDim lipgloss.Color
	colorSecondary  lipgloss.Color
	colorAccent     lipgloss.Color
	colorPurple     lipgloss.Color
	colorCyan       lipgloss.Color
	colorDanger     lipgloss.Color
	colorMuted      lipgloss.Color
	colorSubtle     lipgloss.Color
	colorBg         lipgloss.Color
	colorBgSelected lipgloss.Color
	colorWhite      lipgloss.Color
	colorDark       lipgloss.Color
	colorBgApp      lipgloss.Color
)

// Styles are rebuilt by buildStyles on every theme change.
var (
	appTitleStyle         lipgloss.Style
	headerBarStyle        lipgloss.Style
	breadcrumbStyle       lipgloss.Style
	breadcrumbActiveStyle lipgloss.Style

	tabActiveStyle   lipgloss.Style
	tabInactiveStyle lipgloss.Style
	tabCountStyle    lipgloss.Style

	cardStyle       lipgloss.Style
	activeCardStyle lipgloss.Style
	detailPaneStyle lipgloss.Style

	itemSelectedStyle   lipgloss.Style
	itemUnselectedStyle lipgloss.Style
	itemSubtextStyle    lipgloss.Style

	badgeDueStyle         lipgloss.Style
	badgePointsStyle      lipgloss.Style
	badgeAttachmentsStyle lipgloss.Style
	badgeActiveStyle      lipgloss.Style

	detailTitleStyle    lipgloss.Style
	metaLabelStyle      lipgloss.Style
	metaValueStyle      lipgloss.Style
	descriptionStyle    lipgloss.Style
	attachmentItemStyle lipgloss.Style

	statusBarStyle    lipgloss.Style
	toastSuccessStyle lipgloss.Style
	toastErrorStyle   lipgloss.Style
	helpKeyStyle      lipgloss.Style
	helpDescStyle     lipgloss.Style
	modalBoxStyle     lipgloss.Style
	modalTitleStyle   lipgloss.Style
	inlineHintStyle   lipgloss.Style
	spinnerStyle      lipgloss.Style
	searchBoxStyle    lipgloss.Style
	emptyStateStyle   lipgloss.Style
	tabBarStyle       lipgloss.Style
	headerHintStyle   lipgloss.Style
	modalHintStyle    lipgloss.Style
	pickerRowStyle    lipgloss.Style
	pickerRowSelStyle lipgloss.Style
	previewCardStyle  lipgloss.Style
)

// buildStyles rebuilds colors and styles from activeTheme.
func buildStyles() {
	colorPrimary = lipgloss.Color(activeTheme.Primary)
	colorPrimaryDim = lipgloss.Color(activeTheme.PrimaryDim)
	colorSecondary = lipgloss.Color(activeTheme.Secondary)
	colorAccent = lipgloss.Color(activeTheme.Accent)
	colorPurple = lipgloss.Color(activeTheme.Purple)
	colorCyan = lipgloss.Color(activeTheme.Cyan)
	colorDanger = lipgloss.Color(activeTheme.Danger)
	colorMuted = lipgloss.Color(activeTheme.Muted)
	colorSubtle = lipgloss.Color(activeTheme.Subtle)
	colorBg = lipgloss.Color(activeTheme.BgCard)
	colorBgSelected = lipgloss.Color(activeTheme.BgSelected)
	colorWhite = lipgloss.Color(activeTheme.Text)
	colorDark = lipgloss.Color(activeTheme.TextOnAccent)
	colorBgApp = lipgloss.Color(activeTheme.BgApp)

	// Header & App Bar
	appTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorDark).
		Background(colorPrimary).
		Padding(0, 1).
		MarginRight(1)

	headerBarStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(activeTheme.Text)).
		Background(lipgloss.Color(activeTheme.BgHeader)).
		Padding(0, 1)

	breadcrumbStyle = lipgloss.NewStyle().
		Foreground(colorSubtle)

	breadcrumbActiveStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary)

	// Tabs
	tabActiveStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorDark).
		Background(colorPrimary).
		Padding(0, 1).
		MarginRight(1)

	tabInactiveStyle = lipgloss.NewStyle().
		Foreground(colorSubtle).
		Background(lipgloss.Color(activeTheme.BgTabInactive)).
		Padding(0, 1).
		MarginRight(1)

	tabCountStyle = lipgloss.NewStyle().
		Foreground(colorMuted)

	// Cards & Panes
	cardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(activeTheme.Border)).
		Padding(0, 1)

	activeCardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(0, 1)

	detailPaneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(activeTheme.BorderFocus)).
		Padding(0, 1)

	// List Items
	itemSelectedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary).
		Background(colorBgSelected).
		Padding(0, 1)

	itemUnselectedStyle = lipgloss.NewStyle().
		Foreground(colorWhite).
		Padding(0, 1)

	itemSubtextStyle = lipgloss.NewStyle().
		Foreground(colorMuted)

	// Badges
	badgeDueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(activeTheme.BadgeDueFg)).
		Background(lipgloss.Color(activeTheme.BadgeDueBg)).
		Bold(true).
		Padding(0, 1)

	badgePointsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(activeTheme.BadgePointsFg)).
		Background(lipgloss.Color(activeTheme.BadgePointsBg)).
		Bold(true).
		Padding(0, 1)

	badgeAttachmentsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(activeTheme.BadgeAttachFg)).
		Background(lipgloss.Color(activeTheme.BadgeAttachBg)).
		Bold(true).
		Padding(0, 1)

	badgeActiveStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(activeTheme.BadgeActiveFg)).
		Background(lipgloss.Color(activeTheme.BadgeActiveBg)).
		Bold(true).
		Padding(0, 1)

	// Detail View Sections
	detailTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		MarginBottom(1)

	metaLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorSubtle)

	metaValueStyle = lipgloss.NewStyle().
		Foreground(colorWhite)

	descriptionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(activeTheme.Text))

	attachmentItemStyle = lipgloss.NewStyle().
		Foreground(colorSecondary).
		Underline(true)

	// Footer & Status Bar
	statusBarStyle = lipgloss.NewStyle().
		Foreground(colorSubtle).
		Background(lipgloss.Color(activeTheme.BgApp)).
		Padding(0, 1)

	toastSuccessStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorDark).
		Background(colorPrimary).
		Padding(0, 1)

	toastErrorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Background(colorDanger).
		Padding(0, 1)

	helpKeyStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary)

	helpDescStyle = lipgloss.NewStyle().
		Foreground(colorMuted)

	// Modal / Overlay
	modalBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(colorPrimary).
		Background(lipgloss.Color(activeTheme.BgModal)).
		Padding(1, 2)

	modalTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(activeTheme.ModalTitleFg)).
		Background(lipgloss.Color(activeTheme.ModalTitleBg)).
		Padding(0, 2)

	// Shared
	inlineHintStyle = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Italic(true)

	spinnerStyle = lipgloss.NewStyle().Foreground(colorPrimary)

	searchBoxStyle = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Background(lipgloss.Color(activeTheme.BgSearch)).
		Padding(0, 1)

	emptyStateStyle = lipgloss.NewStyle().
		Foreground(colorSubtle)

	tabBarStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(activeTheme.BgApp))

	headerHintStyle = lipgloss.NewStyle().
		Foreground(colorSubtle)

	modalHintStyle = lipgloss.NewStyle().
		Foreground(colorSubtle).
		Italic(true)

	// Picker rows need their own background, otherwise lipgloss leaves the
	// padding cell beside the swatches unpainted.
	pickerRowStyle = lipgloss.NewStyle().
		Foreground(colorWhite).
		Background(lipgloss.Color(activeTheme.BgModal)).
		Padding(0, 1)

	pickerRowSelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary).
		Background(colorBgSelected).
		Padding(0, 1)

	previewCardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(activeTheme.Border)).
		Background(lipgloss.Color(activeTheme.BgCard))
}

// Paint helpers
// colorsEnabled reports whether the terminal can render color.
func colorsEnabled() bool {
	return lipgloss.ColorProfile() != termenv.Ascii
}

// sgrSeq returns the SGR sequence for a color, or "" without color support.
func sgrSeq(hex string, background bool) string {
	profile := lipgloss.ColorProfile()
	if profile == termenv.Ascii {
		return ""
	}
	c := profile.Color(hex)
	if c == nil {
		return ""
	}
	seq := profile.Convert(c).Sequence(background)
	if seq == "" {
		return ""
	}
	return termenv.CSI + seq + "m"
}

// canvas paints text with explicit colors, then restores the parent colors.
// No reset is emitted, so the enclosing block's background survives.
type canvas struct {
	fg string
	bg string
}

func newCanvas(fg, bg string) canvas {
	return canvas{fg: fg, bg: bg}
}

// paint draws text in fg/bg, then returns to the canvas colors.
func (c canvas) paint(text, fg, bg string) string {
	return sgrSeq(fg, false) + sgrSeq(bg, true) + text + sgrSeq(c.fg, false) + sgrSeq(c.bg, true)
}

// fill pads text to width with canvas-colored spaces.
func (c canvas) fill(text string, width int) string {
	if pad := width - lipgloss.Width(text); pad > 0 {
		return text + sgrSeq(c.bg, true) + strings.Repeat(" ", pad)
	}
	return text
}
