package ui

import (
	"regexp"
	"strings"
	"testing"

	"github.com/4rjxn/classroom/internal/models"
	"github.com/4rjxn/classroom/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var hexColor = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

func themeColors(t Theme) map[string]string {
	return map[string]string{
		"Primary": t.Primary, "PrimaryDim": t.PrimaryDim, "Secondary": t.Secondary,
		"Accent": t.Accent, "Purple": t.Purple, "Cyan": t.Cyan, "Danger": t.Danger,
		"Text": t.Text, "Subtle": t.Subtle, "Muted": t.Muted, "TextOnAccent": t.TextOnAccent,
		"BgApp": t.BgApp, "BgCard": t.BgCard, "BgSelected": t.BgSelected,
		"BgHeader": t.BgHeader, "BgTabInactive": t.BgTabInactive, "BgSearch": t.BgSearch,
		"BgModal": t.BgModal, "Border": t.Border, "BorderFocus": t.BorderFocus,
		"ModalTitleFg": t.ModalTitleFg, "ModalTitleBg": t.ModalTitleBg,
		"BadgeDueFg": t.BadgeDueFg, "BadgeDueBg": t.BadgeDueBg,
		"BadgePointsFg": t.BadgePointsFg, "BadgePointsBg": t.BadgePointsBg,
		"BadgeAttachFg": t.BadgeAttachFg, "BadgeAttachBg": t.BadgeAttachBg,
		"BadgeActiveFg": t.BadgeActiveFg, "BadgeActiveBg": t.BadgeActiveBg,
	}
}

func TestThemesAreCompleteAndUnique(t *testing.T) {
	if len(themes) < 5 {
		t.Fatalf("expected a decent theme lineup, got %d", len(themes))
	}

	seen := make(map[string]bool, len(themes))
	for _, th := range themes {
		if th.Name == "" || th.Label == "" {
			t.Errorf("theme is missing a name or label: %+v", th)
		}
		if seen[th.Name] {
			t.Errorf("duplicate theme name %q", th.Name)
		}
		seen[th.Name] = true

		for field, value := range themeColors(th) {
			if value == "" {
				t.Errorf("theme %q: field %s is empty", th.Name, field)
				continue
			}
			if !hexColor.MatchString(value) {
				t.Errorf("theme %q: field %s = %q is not a #RRGGBB color", th.Name, field, value)
			}
		}
	}
}

func TestThemeLookupIsCaseInsensitive(t *testing.T) {
	if _, ok := themeByName("  NoRd  "); !ok {
		t.Error("expected case/whitespace insensitive lookup to match nord")
	}
	if th, ok := themeByName("Dracula"); !ok || th.Name != "dracula" {
		t.Error("expected label lookup to resolve to the dracula theme")
	}
	if _, ok := themeByName("does-not-exist"); ok {
		t.Error("expected unknown theme to miss")
	}
	if _, ok := themeByName(""); ok {
		t.Error("expected empty name to miss")
	}
}

func TestApplyThemeRebuildsStyles(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })
	defer applyTheme(defaultTheme())

	applyTheme(defaultTheme())
	before := colorPrimary
	emerald := appTitleStyle.Render("x")

	nord, ok := themeByName("nord")
	if !ok {
		t.Fatal("nord theme missing")
	}
	applyTheme(nord)

	if string(colorPrimary) != nord.Primary {
		t.Errorf("colorPrimary = %q, want %q", colorPrimary, nord.Primary)
	}
	if activeTheme.Name != "nord" {
		t.Errorf("activeTheme = %q, want nord", activeTheme.Name)
	}
	if after := appTitleStyle.Render("x"); after == emerald {
		t.Error("expected rendered styles to change after switching themes")
	}
	if before == colorPrimary {
		t.Error("expected the palette to actually change")
	}
}

func TestCycleThemeWraps(t *testing.T) {
	defer applyTheme(defaultTheme())

	applyTheme(themes[len(themes)-1])
	if got := cycleTheme(); got.Name != themes[0].Name {
		t.Errorf("cycle from last theme = %q, want %q", got.Name, themes[0].Name)
	}

	applyTheme(themes[0])
	if got := cycleTheme(); got.Name != themes[1].Name {
		t.Errorf("cycle from first theme = %q, want %q", got.Name, themes[1].Name)
	}
}

func TestResolveInitialThemePrecedence(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	if got := resolveInitialTheme(models.Config{}); got.Name != defaultTheme().Name {
		t.Errorf("with nothing set, got %q want %q", got.Name, defaultTheme().Name)
	}

	if err := utils.SaveTheme("dracula"); err != nil {
		t.Fatalf("SaveTheme: %v", err)
	}
	if got := resolveInitialTheme(models.Config{}); got.Name != "dracula" {
		t.Errorf("saved preference ignored, got %q", got.Name)
	}

	if got := resolveInitialTheme(models.Config{Theme: "gruvbox-dark"}); got.Name != "gruvbox-dark" {
		t.Errorf("config theme should win over saved preference, got %q", got.Name)
	}

	if got := resolveInitialTheme(models.Config{Theme: "nope"}); got.Name != "dracula" {
		t.Errorf("invalid config theme should fall back to preference, got %q", got.Name)
	}
}

func TestThemePickerKeepsLayout(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })
	defer applyTheme(defaultTheme())

	m := NewUiStateModel("token", models.Config{})
	m.width, m.height = 120, 45
	m.showThemes = true
	m.themeCursor = 0

	out := m.View()
	if strings.TrimSpace(out) == "" {
		t.Fatal("theme picker rendered nothing")
	}

	for i, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		if w := lipgloss.Width(line); w != 120 {
			t.Fatalf("line %d has width %d, want 120 (layout drift)", i, w)
		}
	}

	if !strings.Contains(out, themes[0].Label) {
		t.Errorf("picker should list %q", themes[0].Label)
	}
}

func TestThemeSwitchingKeybinding(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	defer applyTheme(defaultTheme())

	m := NewUiStateModel("token", models.Config{})
	m.width, m.height = 100, 30
	m.loading = false

	updated, _ := m.Update(teaKey("t"))
	next := updated.(UiStateModel)

	if next.actionsThemeName() == defaultTheme().Name {
		t.Error("pressing t should move to a different theme")
	}
	if got := utils.ReadTheme(); got != next.actionsThemeName() {
		t.Errorf("theme %q was applied but %q was persisted", next.actionsThemeName(), got)
	}

	updated, _ = next.Update(teaKey("T"))
	opened := updated.(UiStateModel)
	if !opened.showThemes {
		t.Fatal("pressing T should open the theme picker")
	}
	if opened.themeBeforePicker != next.actionsThemeName() {
		t.Errorf("picker should remember %q, got %q", next.actionsThemeName(), opened.themeBeforePicker)
	}
}

// actionsThemeName exposes the active palette id for assertions.
func (m UiStateModel) actionsThemeName() string { return activeTheme.Name }

func TestThemePickerEscRevertsPreview(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	defer applyTheme(defaultTheme())

	applyTheme(defaultTheme())
	m := NewUiStateModel("token", models.Config{})
	m.width, m.height = 100, 30

	updated, _ := m.Update(teaKey("T"))
	opened := updated.(UiStateModel)

	updated, _ = opened.Update(teaKey("j"))
	previewed := updated.(UiStateModel)
	if previewed.actionsThemeName() == defaultTheme().Name {
		t.Fatal("expected j to preview the next theme")
	}

	updated, _ = previewed.Update(teaKey("esc"))
	reverted := updated.(UiStateModel)
	if reverted.showThemes {
		t.Error("esc should close the picker")
	}
	if reverted.actionsThemeName() != defaultTheme().Name {
		t.Errorf("esc should revert the preview, still on %q", reverted.actionsThemeName())
	}
}

// teaKey builds a key message from a key name for tests.
func teaKey(name string) tea.KeyMsg {
	switch name {
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(name)}
	}
}

func TestThemeSwatchesDoNotClearRowBackground(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })
	defer applyTheme(defaultTheme())

	dracula, _ := themeByName("dracula")
	swatches := themeSwatches(dracula)

	if got := strings.Count(swatches, "●"); got != 6 {
		t.Fatalf("expected 6 dots, got %d", got)
	}
	if strings.Contains(swatches, "\x1b[0m") {
		t.Error("swatches must not emit a full reset: it clears the enclosing row background")
	}
	if got := strings.Count(swatches, "\x1b[38;2;"); got != 6 {
		t.Errorf("expected 6 foreground sequences, got %d in %q", got, swatches)
	}
	if !strings.HasSuffix(swatches, "\x1b[39m") {
		t.Errorf("swatches should reset the foreground when done, got %q", swatches)
	}
}

func TestThemeSwatchesArePlainWithoutColor(t *testing.T) {
	lipgloss.SetColorProfile(termenv.Ascii)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })
	defer applyTheme(defaultTheme())

	got := themeSwatches(themes[0])
	if got != "●●●●●●" {
		t.Errorf("without color support swatches should be plain dots, got %q", got)
	}
}

func TestSelectedPickerRowKeepsBackgroundUnderSwatches(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })
	defer applyTheme(defaultTheme())

	m := NewUiStateModel("token", models.Config{})
	m.width, m.height = 120, 45
	m.showThemes = true
	m.themeCursor = 3

	var row string
	for _, line := range strings.Split(m.View(), "\n") {
		if strings.Contains(line, themes[3].Label) {
			row = line
			break
		}
	}
	if row == "" {
		t.Fatalf("could not find the selected row for %q", themes[3].Label)
	}

	first, last := strings.Index(row, "●"), strings.LastIndex(row, "●")
	if first < 0 || last <= first {
		t.Fatalf("expected a run of swatch dots in %q", row)
	}
	if hole := row[first:last]; strings.Contains(hole, "\x1b[0m") {
		t.Errorf("row background is cleared between the swatch dots: %q", hole)
	}
}

func TestPickerSwatchColumnIsAligned(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })
	defer applyTheme(defaultTheme())

	m := NewUiStateModel("token", models.Config{})
	m.width, m.height = 120, 45
	m.showThemes = true
	m.themeCursor = 0

	// The current theme's row carries the extra checkmark slot.
	columns := map[int][]string{}
	for _, name := range []string{"emerald", "rose-pine", "catppuccin-latte", "gruvbox-light"} {
		th, _ := themeByName(name)
		applyTheme(th)

		m.themeCursor = 0
		for _, line := range strings.Split(m.View(), "\n") {
			// One dot per accent; the preview's tab strip has a dot too.
			if strings.Count(line, "●") != len(themeSwatchAccents(themes[0])) {
				continue
			}
			dot := strings.Index(line, "●")
			// Every row must start its swatches at the same column.
			columns[lipgloss.Width(line[:dot])] = append(columns[lipgloss.Width(line[:dot])], name)
		}
	}

	if len(columns) != 1 {
		t.Errorf("swatch column drifts across rows/themes: %v", columns)
	}
}

func TestThemeLabelsFitAlignedField(t *testing.T) {
	for _, th := range themes {
		if w := lipgloss.Width(th.Label); w > themeLabelWidth {
			t.Errorf("theme %q label is %d wide, exceeds the aligned field %d", th.Name, w, themeLabelWidth)
		}
	}
}

func TestPickerStylesPaintTheirBackground(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })
	defer applyTheme(defaultTheme())

	// Without a background, lipgloss leaves the row's padding unpainted.
	for name, style := range map[string]lipgloss.Style{
		"pickerRowStyle":    pickerRowStyle,
		"pickerRowSelStyle": pickerRowSelStyle,
		"previewCardStyle":  previewCardStyle,
	} {
		if style.GetBackground() == nil {
			t.Errorf("%s has no background: its padding would show the terminal background", name)
		}
	}
}

func TestCanvasPaintDoesNotResetEnclosingStyle(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })

	c := newCanvas("#FFFFFF", "#101010")
	out := c.paint("chip", "#FF0000", "#00FF00")

	if strings.Contains(out, "\x1b[0m") {
		t.Errorf("paint must not emit a reset, got %q", out)
	}
	if !strings.HasSuffix(out, sgrSeq("#FFFFFF", false)+sgrSeq("#101010", true)) {
		t.Errorf("paint must restore the canvas colors, got %q", out)
	}
	if !strings.HasPrefix(out, sgrSeq("#FF0000", false)+sgrSeq("#00FF00", true)) {
		t.Errorf("paint must apply the chip colors, got %q", out)
	}
}

func TestCanvasFillPadsWithBackground(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })

	c := newCanvas("#FFFFFF", "#101010")
	filled := c.fill(c.paint("chip", "#FF0000", "#00FF00"), 12)

	if got := lipgloss.Width(filled); got != 12 {
		t.Errorf("filled width = %d, want 12", got)
	}
	if strings.Contains(filled, "\x1b[0m") {
		t.Errorf("fill must not emit a reset, got %q", filled)
	}
	if !strings.Contains(filled, sgrSeq("#101010", true)) {
		t.Error("fill should pad using the canvas background")
	}
}

func TestCanvasPaintIsPlainWithoutColor(t *testing.T) {
	lipgloss.SetColorProfile(termenv.Ascii)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })

	c := newCanvas("#FFFFFF", "#101010")
	if got := c.paint("chip", "#FF0000", "#00FF00"); got != "chip" {
		t.Errorf("without color support paint should be plain text, got %q", got)
	}
	if got := c.fill("chip", 8); got != "chip    " {
		t.Errorf("without color support fill should pad with plain spaces, got %q", got)
	}
}
