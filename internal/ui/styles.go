package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorPrimary    = lipgloss.Color("#10B981") // Emerald Green
	colorPrimaryDim = lipgloss.Color("#047857")
	colorSecondary  = lipgloss.Color("#3B82F6") // Google Blue
	colorAccent     = lipgloss.Color("#F59E0B") // Amber / Due date
	colorPurple     = lipgloss.Color("#8B5CF6") // Materials purple
	colorCyan       = lipgloss.Color("#06B6D4") // Announcements cyan
	colorDanger     = lipgloss.Color("#EF4444") // Red
	colorMuted      = lipgloss.Color("#6B7280") // Gray
	colorSubtle     = lipgloss.Color("#9CA3AF") // Light Gray
	colorBg         = lipgloss.Color("#1F2937") // Dark card bg
	colorBgSelected = lipgloss.Color("#2D3748") // Selection bg
	colorWhite      = lipgloss.Color("#F9FAFB")
	colorDark       = lipgloss.Color("#111827")

	// Header & App Bar
	appTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorDark).
			Background(colorPrimary).
			Padding(0, 1).
			MarginRight(1)

	headerBarStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite).
			Background(lipgloss.Color("#1E293B")).
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
				Background(lipgloss.Color("#374151")).
				Padding(0, 1).
				MarginRight(1)

	tabCountStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	// Cards & Panes
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#374151")).
			Padding(0, 1)

	activeCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)

	detailPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4B5563")).
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
			Foreground(lipgloss.Color("#78350F")).
			Background(lipgloss.Color("#FDE68A")).
			Bold(true).
			Padding(0, 1)

	badgePointsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1E3A8A")).
				Background(lipgloss.Color("#BFDBFE")).
				Bold(true).
				Padding(0, 1)

	badgeAttachmentsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4C1D95")).
				Background(lipgloss.Color("#DDD6FE")).
				Bold(true).
				Padding(0, 1)

	badgeActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#064E3B")).
				Background(lipgloss.Color("#A7F3D0")).
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
				Foreground(lipgloss.Color("#E5E7EB"))

	attachmentItemStyle = lipgloss.NewStyle().
				Foreground(colorSecondary).
				Underline(true)

	// Footer & Status Bar
	statusBarStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			Background(lipgloss.Color("#111827")).
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
			Background(lipgloss.Color("#1E293B")).
			Padding(1, 2)
)
