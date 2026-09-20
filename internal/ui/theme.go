package ui

import (
	"strings"

	"github.com/4rjxn/classroom/internal/models"
	"github.com/4rjxn/classroom/internal/utils"
	"github.com/charmbracelet/lipgloss"
)

// Theme is a color palette for the TUI.
type Theme struct {
	Name  string // id used in config.toml and prefs
	Label string // shown in the picker

	// Accents
	Primary    string
	PrimaryDim string
	Secondary  string
	Accent     string
	Purple     string
	Cyan       string
	Danger     string

	// Text
	Text         string
	Subtle       string
	Muted        string
	TextOnAccent string

	// Surfaces
	BgApp         string
	BgCard        string
	BgSelected    string
	BgHeader      string
	BgTabInactive string
	BgSearch      string
	BgModal       string

	// Lines
	Border      string
	BorderFocus string

	// Modal title bar
	ModalTitleFg string
	ModalTitleBg string

	// Badges
	BadgeDueFg    string
	BadgeDueBg    string
	BadgePointsFg string
	BadgePointsBg string
	BadgeAttachFg string
	BadgeAttachBg string
	BadgeActiveFg string
	BadgeActiveBg string
}

// themes holds the built-in palettes. Order is picker order; themes[0] is the
// default. Append a Theme here to add one.
var themes = []Theme{
	{
		Name:  "emerald",
		Label: "Emerald (default)",

		Primary: "#10B981", PrimaryDim: "#047857", Secondary: "#3B82F6",
		Accent: "#F59E0B", Purple: "#8B5CF6", Cyan: "#06B6D4", Danger: "#EF4444",

		Text: "#F9FAFB", Subtle: "#9CA3AF", Muted: "#6B7280", TextOnAccent: "#111827",

		BgApp: "#111827", BgCard: "#1F2937", BgSelected: "#2D3748", BgHeader: "#1E293B",
		BgTabInactive: "#374151", BgSearch: "#1F2937", BgModal: "#1E293B",

		Border: "#374151", BorderFocus: "#4B5563",

		ModalTitleFg: "#F9FAFB", ModalTitleBg: "#047857",

		BadgeDueFg: "#78350F", BadgeDueBg: "#FDE68A",
		BadgePointsFg: "#1E3A8A", BadgePointsBg: "#BFDBFE",
		BadgeAttachFg: "#4C1D95", BadgeAttachBg: "#DDD6FE",
		BadgeActiveFg: "#064E3B", BadgeActiveBg: "#A7F3D0",
	},
	{
		Name:  "catppuccin-mocha",
		Label: "Catppuccin Mocha",

		Primary: "#A6E3A1", PrimaryDim: "#3F6F4A", Secondary: "#89B4FA",
		Accent: "#F9E2AF", Purple: "#CBA6F7", Cyan: "#94E2D5", Danger: "#F38BA8",

		Text: "#CDD6F4", Subtle: "#A6ADC8", Muted: "#6C7086", TextOnAccent: "#1E1E2E",

		BgApp: "#11111B", BgCard: "#1E1E2E", BgSelected: "#313244", BgHeader: "#1E1E2E",
		BgTabInactive: "#313244", BgSearch: "#1E1E2E", BgModal: "#1E1E2E",

		Border: "#313244", BorderFocus: "#585B70",

		ModalTitleFg: "#CDD6F4", ModalTitleBg: "#3F6F4A",

		BadgeDueFg: "#1E1E2E", BadgeDueBg: "#F9E2AF",
		BadgePointsFg: "#1E1E2E", BadgePointsBg: "#89B4FA",
		BadgeAttachFg: "#1E1E2E", BadgeAttachBg: "#CBA6F7",
		BadgeActiveFg: "#1E1E2E", BadgeActiveBg: "#A6E3A1",
	},
	{
		Name:  "catppuccin-latte",
		Label: "Catppuccin Latte (light)",

		Primary: "#40A02B", PrimaryDim: "#2F7A21", Secondary: "#1E66F5",
		Accent: "#DF8E1D", Purple: "#8839EF", Cyan: "#179299", Danger: "#D20F39",

		Text: "#4C4F69", Subtle: "#6C6F85", Muted: "#9CA0B0", TextOnAccent: "#FFFFFF",

		BgApp: "#DCE0E8", BgCard: "#EFF1F5", BgSelected: "#CCD0DA", BgHeader: "#E6E9EF",
		BgTabInactive: "#CCD0DA", BgSearch: "#EFF1F5", BgModal: "#EFF1F5",

		Border: "#BCC0CC", BorderFocus: "#9CA0B0",

		ModalTitleFg: "#FFFFFF", ModalTitleBg: "#2F7A21",

		BadgeDueFg: "#6B4B00", BadgeDueBg: "#FFE9A8",
		BadgePointsFg: "#10306B", BadgePointsBg: "#C8D8FF",
		BadgeAttachFg: "#4B1D95", BadgeAttachBg: "#E0D4FF",
		BadgeActiveFg: "#0A4A1E", BadgeActiveBg: "#C8F0C8",
	},
	{
		Name:  "dracula",
		Label: "Dracula",

		Primary: "#50FA7B", PrimaryDim: "#2F7A45", Secondary: "#8BE9FD",
		Accent: "#FFB86C", Purple: "#BD93F9", Cyan: "#8BE9FD", Danger: "#FF5555",

		Text: "#F8F8F2", Subtle: "#B0B3C8", Muted: "#6272A4", TextOnAccent: "#282A36",

		BgApp: "#21222C", BgCard: "#282A36", BgSelected: "#44475A", BgHeader: "#282A36",
		BgTabInactive: "#44475A", BgSearch: "#282A36", BgModal: "#282A36",

		Border: "#44475A", BorderFocus: "#6272A4",

		ModalTitleFg: "#F8F8F2", ModalTitleBg: "#2F7A45",

		BadgeDueFg: "#282A36", BadgeDueBg: "#F1FA8C",
		BadgePointsFg: "#282A36", BadgePointsBg: "#8BE9FD",
		BadgeAttachFg: "#282A36", BadgeAttachBg: "#BD93F9",
		BadgeActiveFg: "#282A36", BadgeActiveBg: "#50FA7B",
	},
	{
		Name:  "nord",
		Label: "Nord",

		Primary: "#88C0D0", PrimaryDim: "#4A7285", Secondary: "#81A1C1",
		Accent: "#EBCB8B", Purple: "#B48EAD", Cyan: "#8FBCBB", Danger: "#BF616A",

		Text: "#ECEFF4", Subtle: "#D8DEE9", Muted: "#7B88A1", TextOnAccent: "#2E3440",

		BgApp: "#2E3440", BgCard: "#3B4252", BgSelected: "#434C5E", BgHeader: "#3B4252",
		BgTabInactive: "#434C5E", BgSearch: "#3B4252", BgModal: "#3B4252",

		Border: "#434C5E", BorderFocus: "#4C566A",

		ModalTitleFg: "#ECEFF4", ModalTitleBg: "#4A7285",

		BadgeDueFg: "#2E3440", BadgeDueBg: "#EBCB8B",
		BadgePointsFg: "#2E3440", BadgePointsBg: "#81A1C1",
		BadgeAttachFg: "#2E3440", BadgeAttachBg: "#B48EAD",
		BadgeActiveFg: "#2E3440", BadgeActiveBg: "#A3BE8C",
	},
	{
		Name:  "gruvbox-dark",
		Label: "Gruvbox Dark",

		Primary: "#B8BB26", PrimaryDim: "#6B6D18", Secondary: "#83A598",
		Accent: "#FABD2F", Purple: "#D3869B", Cyan: "#8EC07C", Danger: "#FB4934",

		Text: "#EBDBB2", Subtle: "#D5C4A1", Muted: "#A89984", TextOnAccent: "#282828",

		BgApp: "#1D2021", BgCard: "#282828", BgSelected: "#3C3836", BgHeader: "#282828",
		BgTabInactive: "#3C3836", BgSearch: "#282828", BgModal: "#282828",

		Border: "#3C3836", BorderFocus: "#504945",

		ModalTitleFg: "#EBDBB2", ModalTitleBg: "#6B6D18",

		BadgeDueFg: "#282828", BadgeDueBg: "#FABD2F",
		BadgePointsFg: "#282828", BadgePointsBg: "#83A598",
		BadgeAttachFg: "#282828", BadgeAttachBg: "#D3869B",
		BadgeActiveFg: "#282828", BadgeActiveBg: "#B8BB26",
	},
	{
		Name:  "gruvbox-light",
		Label: "Gruvbox Light (light)",

		Primary: "#79740E", PrimaryDim: "#574F0A", Secondary: "#076678",
		Accent: "#B57614", Purple: "#8F3F71", Cyan: "#427B58", Danger: "#9D0006",

		Text: "#3C3836", Subtle: "#504945", Muted: "#7C6F64", TextOnAccent: "#FBF1C7",

		BgApp: "#F2E5BC", BgCard: "#FBF1C7", BgSelected: "#EBDBB2", BgHeader: "#EBDBB2",
		BgTabInactive: "#D5C4A1", BgSearch: "#FBF1C7", BgModal: "#FBF1C7",

		Border: "#D5C4A1", BorderFocus: "#BDAE93",

		ModalTitleFg: "#FBF1C7", ModalTitleBg: "#574F0A",

		BadgeDueFg: "#6B4F00", BadgeDueBg: "#F2DFA0",
		BadgePointsFg: "#0B3D4D", BadgePointsBg: "#BCD8DD",
		BadgeAttachFg: "#5C2A48", BadgeAttachBg: "#E3C6D8",
		BadgeActiveFg: "#2F3D0A", BadgeActiveBg: "#D8D9A0",
	},
	{
		Name:  "tokyo-night",
		Label: "Tokyo Night",

		Primary: "#7AA2F7", PrimaryDim: "#3D5A9E", Secondary: "#7DCFFF",
		Accent: "#E0AF68", Purple: "#BB9AF7", Cyan: "#7DCFFF", Danger: "#F7768E",

		Text: "#C0CAF5", Subtle: "#A9B1D6", Muted: "#565F89", TextOnAccent: "#1A1B26",

		BgApp: "#16161E", BgCard: "#1A1B26", BgSelected: "#292E42", BgHeader: "#1A1B26",
		BgTabInactive: "#292E42", BgSearch: "#1A1B26", BgModal: "#1A1B26",

		Border: "#292E42", BorderFocus: "#414868",

		ModalTitleFg: "#C0CAF5", ModalTitleBg: "#3D5A9E",

		BadgeDueFg: "#1A1B26", BadgeDueBg: "#E0AF68",
		BadgePointsFg: "#1A1B26", BadgePointsBg: "#7AA2F7",
		BadgeAttachFg: "#1A1B26", BadgeAttachBg: "#BB9AF7",
		BadgeActiveFg: "#1A1B26", BadgeActiveBg: "#9ECE6A",
	},
	{
		Name:  "one-dark",
		Label: "One Dark",

		Primary: "#61AFEF", PrimaryDim: "#2C5C8A", Secondary: "#56B6C2",
		Accent: "#E5C07B", Purple: "#C678DD", Cyan: "#56B6C2", Danger: "#E06C75",

		Text: "#ABB2BF", Subtle: "#B6BDC9", Muted: "#5C6370", TextOnAccent: "#282C34",

		BgApp: "#21252B", BgCard: "#282C34", BgSelected: "#3E4451", BgHeader: "#282C34",
		BgTabInactive: "#3E4451", BgSearch: "#282C34", BgModal: "#282C34",

		Border: "#3E4451", BorderFocus: "#4B5263",

		ModalTitleFg: "#ABB2BF", ModalTitleBg: "#2C5C8A",

		BadgeDueFg: "#282C34", BadgeDueBg: "#E5C07B",
		BadgePointsFg: "#282C34", BadgePointsBg: "#61AFEF",
		BadgeAttachFg: "#282C34", BadgeAttachBg: "#C678DD",
		BadgeActiveFg: "#282C34", BadgeActiveBg: "#98C379",
	},
	{
		Name:  "rose-pine",
		Label: "Rosé Pine",

		Primary: "#C4A7E7", PrimaryDim: "#5B4A7A", Secondary: "#9CCFD8",
		Accent: "#F6C177", Purple: "#C4A7E7", Cyan: "#9CCFD8", Danger: "#EB6F92",

		Text: "#E0DEF4", Subtle: "#908CAA", Muted: "#6E6A86", TextOnAccent: "#191724",

		BgApp: "#191724", BgCard: "#1F1D2E", BgSelected: "#26233A", BgHeader: "#1F1D2E",
		BgTabInactive: "#26233A", BgSearch: "#1F1D2E", BgModal: "#1F1D2E",

		Border: "#26233A", BorderFocus: "#403D52",

		ModalTitleFg: "#E0DEF4", ModalTitleBg: "#5B4A7A",

		BadgeDueFg: "#191724", BadgeDueBg: "#F6C177",
		BadgePointsFg: "#191724", BadgePointsBg: "#9CCFD8",
		BadgeAttachFg: "#191724", BadgeAttachBg: "#C4A7E7",
		BadgeActiveFg: "#E0DEF4", BadgeActiveBg: "#31748F",
	},
	{
		Name:  "solarized-dark",
		Label: "Solarized Dark",

		Primary: "#2AA198", PrimaryDim: "#1A6B64", Secondary: "#268BD2",
		Accent: "#B58900", Purple: "#6C71C4", Cyan: "#2AA198", Danger: "#DC322F",

		Text: "#93A1A1", Subtle: "#839496", Muted: "#586E75", TextOnAccent: "#002B36",

		BgApp: "#002B36", BgCard: "#073642", BgSelected: "#0E4B5A", BgHeader: "#073642",
		BgTabInactive: "#0E4B5A", BgSearch: "#073642", BgModal: "#073642",

		Border: "#12505F", BorderFocus: "#586E75",

		ModalTitleFg: "#93A1A1", ModalTitleBg: "#1A6B64",

		BadgeDueFg: "#002B36", BadgeDueBg: "#B58900",
		BadgePointsFg: "#002B36", BadgePointsBg: "#268BD2",
		BadgeAttachFg: "#002B36", BadgeAttachBg: "#6C71C4",
		BadgeActiveFg: "#002B36", BadgeActiveBg: "#859900",
	},
	{
		Name:  "synthwave",
		Label: "Synthwave",

		Primary: "#FF7EDB", PrimaryDim: "#8A2F75", Secondary: "#36F9F6",
		Accent: "#FEDE5D", Purple: "#B893CE", Cyan: "#36F9F6", Danger: "#FE4450",

		Text: "#F8F8F2", Subtle: "#CDBDE0", Muted: "#7A6A8F", TextOnAccent: "#241B2F",

		BgApp: "#1A1127", BgCard: "#241B2F", BgSelected: "#342B45", BgHeader: "#241B2F",
		BgTabInactive: "#342B45", BgSearch: "#241B2F", BgModal: "#241B2F",

		Border: "#3B2F4F", BorderFocus: "#5B4A75",

		ModalTitleFg: "#F8F8F2", ModalTitleBg: "#8A2F75",

		BadgeDueFg: "#241B2F", BadgeDueBg: "#FEDE5D",
		BadgePointsFg: "#241B2F", BadgePointsBg: "#36F9F6",
		BadgeAttachFg: "#241B2F", BadgeAttachBg: "#B893CE",
		BadgeActiveFg: "#241B2F", BadgeActiveBg: "#72F1B8",
	},
	{
		Name:  "mono",
		Label: "Monochrome",

		Primary: "#E0E0E0", PrimaryDim: "#4A4A4A", Secondary: "#B8B8B8",
		Accent: "#CFCFCF", Purple: "#C0C0C0", Cyan: "#D0D0D0", Danger: "#F0F0F0",

		Text: "#E0E0E0", Subtle: "#A8A8A8", Muted: "#7A7A7A", TextOnAccent: "#1C1C1C",

		BgApp: "#141414", BgCard: "#1C1C1C", BgSelected: "#2E2E2E", BgHeader: "#1C1C1C",
		BgTabInactive: "#2E2E2E", BgSearch: "#1C1C1C", BgModal: "#1C1C1C",

		Border: "#303030", BorderFocus: "#4A4A4A",

		ModalTitleFg: "#E0E0E0", ModalTitleBg: "#4A4A4A",

		BadgeDueFg: "#1C1C1C", BadgeDueBg: "#D0D0D0",
		BadgePointsFg: "#1C1C1C", BadgePointsBg: "#B8B8B8",
		BadgeAttachFg: "#1C1C1C", BadgeAttachBg: "#C8C8C8",
		BadgeActiveFg: "#1C1C1C", BadgeActiveBg: "#DCDCDC",
	},
}

// activeTheme is the palette currently applied.
var activeTheme Theme

// themeLabelWidth keeps the picker's swatch column aligned.
var themeLabelWidth = func() int {
	w := 0
	for _, t := range themes {
		if n := lipgloss.Width(t.Label); n > w {
			w = n
		}
	}
	return w
}()

func init() {
	applyTheme(themes[0])
}

// themeByName resolves an id or label, case-insensitively.
func themeByName(name string) (Theme, bool) {
	want := strings.ToLower(strings.TrimSpace(name))
	if want == "" {
		return Theme{}, false
	}
	for _, t := range themes {
		if strings.ToLower(t.Name) == want || strings.ToLower(t.Label) == want {
			return t, true
		}
	}
	return Theme{}, false
}

// defaultTheme returns the first palette.
func defaultTheme() Theme {
	return themes[0]
}

// currentThemeIndex is the position of the active theme.
func currentThemeIndex() int {
	if activeTheme.Name == "" {
		return 0
	}
	for i, t := range themes {
		if t.Name == activeTheme.Name {
			return i
		}
	}
	return 0
}

// cycleTheme returns the next palette, wrapping.
func cycleTheme() Theme {
	return themes[(currentThemeIndex()+1)%len(themes)]
}

// applyTheme applies a palette and rebuilds the styles.
func applyTheme(t Theme) {
	activeTheme = t
	buildStyles()
}

// resolveInitialTheme prefers config, then the saved pick, then the default.
func resolveInitialTheme(cfg models.Config) Theme {
	if t, ok := themeByName(cfg.Theme); ok {
		return t
	}
	if t, ok := themeByName(utils.ReadTheme()); ok {
		return t
	}
	return defaultTheme()
}
