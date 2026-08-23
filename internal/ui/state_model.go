package ui

import (
	"fmt"
	"strings"

	"github.com/4rjxn/classroom/internal/models"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewState int

const (
	viewCourses ViewState = iota
	viewCourseDetail
)

type TabType int

const (
	tabAssignments TabType = iota
	tabMaterials
	tabAnnouncements
	tabInfo
)

type CourseDetailData struct {
	CourseWork    []models.CourseWorkModel
	Materials     []models.CourseWorkMaterial
	Announcements []models.CourseAnnouncement
	Loaded        bool
	Loading       bool
}

type UiStateModel struct {
	State             ViewState
	Token             string
	Config            models.Config
	courses           []models.CourseModel
	filteredCourses   []models.CourseModel
	selectedCourseIdx int
	courseOffset      int
	activeTab         TabType
	tabCursors        [4]int
	tabOffsets        [4]int
	courseData        map[string]*CourseDetailData
	currentCourse     *models.CourseModel
	spinner           spinner.Model
	viewport          viewport.Model
	searchInput       textinput.Model
	searching         bool
	loading           bool
	loadingMsg        string
	showHelp          bool
	showPicker        bool
	pickerCursor      int
	pickerAttachments []models.Attachment
	statusMsg         string
	statusIsErr       bool
	width             int
	height            int
	err               error
}

func NewUiStateModel(token string, cfg models.Config) UiStateModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorPrimary)

	ti := textinput.New()
	ti.Placeholder = "Type to search courses... (Esc to cancel)"
	ti.Prompt = "🔍 "
	ti.CharLimit = 50

	vp := viewport.New(40, 15)

	return UiStateModel{
		State:       viewCourses,
		Token:       token,
		Config:      cfg,
		activeTab:   tabAssignments,
		courseData:  make(map[string]*CourseDetailData),
		spinner:     s,
		viewport:    vp,
		searchInput: ti,
		loading:     true,
		loadingMsg:  "Connecting to Google Classroom...",
		width:       80,
		height:      24,
	}
}

func (m UiStateModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchCoursesCmd(m.Token),
	)
}

func (m UiStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateViewportDimensions()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case statusMsg:
		m.statusMsg = msg.Text
		m.statusIsErr = msg.IsError

	case clearStatusMsg:
		m.statusMsg = ""
		m.statusIsErr = false

	case urlOpenedMsg:
		if msg.Err != nil {
			cmds = append(cmds, flashStatusCmd(fmt.Sprintf("Failed to open URL: %v", msg.Err), true))
		} else {
			cmds = append(cmds, flashStatusCmd("✓ Opened in browser", false))
		}

	case coursesLoadedMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			m.statusMsg = fmt.Sprintf("Error fetching courses: %v", msg.Err)
			m.statusIsErr = true
		} else {
			m.courses = msg.Courses
			m.applyCourseFilter()
			m.selectedCourseIdx = 0
			m.courseOffset = 0
			m.statusMsg = fmt.Sprintf("Loaded %d courses", len(m.courses))
			m.statusIsErr = false
		}

	case courseDetailsLoadedMsg:
		m.loading = false
		data, exists := m.courseData[msg.CourseID]
		if !exists {
			data = &CourseDetailData{}
			m.courseData[msg.CourseID] = data
		}
		data.Loading = false
		data.Loaded = true
		data.CourseWork = msg.CourseWork
		data.Materials = msg.Materials
		data.Announcements = msg.Announcements

		if msg.Err != nil {
			cmds = append(cmds, flashStatusCmd(fmt.Sprintf("Warning: %v", msg.Err), true))
		} else {
			cmds = append(cmds, flashStatusCmd("✓ Course details updated", false))
		}
		m.updateDetailViewport()

	case tea.KeyMsg:
		// 1. Help Modal Interceptions
		if m.showHelp {
			switch msg.String() {
			case "q", "esc", "?", "enter":
				m.showHelp = false
			}
			return m, nil
		}

		// 2. Attachment Picker Modal Interceptions
		if m.showPicker {
			switch msg.String() {
			case "esc", "q":
				m.showPicker = false
			case "j", "down":
				if m.pickerCursor < len(m.pickerAttachments)-1 {
					m.pickerCursor++
				}
			case "k", "up":
				if m.pickerCursor > 0 {
					m.pickerCursor--
				}
			case "enter":
				if len(m.pickerAttachments) > 0 && m.pickerCursor < len(m.pickerAttachments) {
					target := m.pickerAttachments[m.pickerCursor].URL
					m.showPicker = false
					cmds = append(cmds, openBrowserCmd(target))
				}
			}
			return m, tea.Batch(cmds...)
		}

		// 3. Search Mode Interceptions
		if m.searching {
			switch msg.String() {
			case "esc":
				m.searching = false
				m.searchInput.Blur()
			case "enter":
				m.searching = false
				m.searchInput.Blur()
			default:
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				cmds = append(cmds, cmd)
				m.applyCourseFilter()
			}
			return m, tea.Batch(cmds...)
		}

		// 4. Global Hotkeys
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		}

		// 5. View-Specific Navigation
		switch m.State {
		case viewCourses:
			cmd := m.handleCoursesKeys(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case viewCourseDetail:
			cmd := m.handleCourseDetailKeys(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *UiStateModel) handleCoursesKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q":
		return tea.Quit
	case "j", "down":
		if len(m.filteredCourses) > 0 && m.selectedCourseIdx < len(m.filteredCourses)-1 {
			m.selectedCourseIdx++
		}
	case "k", "up":
		if m.selectedCourseIdx > 0 {
			m.selectedCourseIdx--
		}
	case "g", "home":
		m.selectedCourseIdx = 0
		m.courseOffset = 0
	case "G", "end":
		if len(m.filteredCourses) > 0 {
			m.selectedCourseIdx = len(m.filteredCourses) - 1
		}
	case "/":
		m.searching = true
		m.searchInput.Focus()
		return textinput.Blink
	case "r", "ctrl+r":
		m.loading = true
		m.loadingMsg = "Refreshing courses..."
		return fetchCoursesCmd(m.Token)
	case "o":
		if len(m.filteredCourses) > 0 && m.selectedCourseIdx < len(m.filteredCourses) {
			course := m.filteredCourses[m.selectedCourseIdx]
			if course.AlternateLink != "" {
				return openBrowserCmd(course.AlternateLink)
			}
		}
	case "enter":
		if len(m.filteredCourses) > 0 && m.selectedCourseIdx < len(m.filteredCourses) {
			course := m.filteredCourses[m.selectedCourseIdx]
			m.currentCourse = &course
			m.State = viewCourseDetail
			m.activeTab = tabAssignments
			m.tabCursors = [4]int{0, 0, 0, 0}
			m.tabOffsets = [4]int{0, 0, 0, 0}

			data, exists := m.courseData[course.Id]
			if !exists || !data.Loaded {
				m.loading = true
				m.loadingMsg = fmt.Sprintf("Loading %s...", course.Name)
				return fetchCourseDataCmd(m.Token, course.Id)
			}
			m.updateDetailViewport()
		}
	}
	return nil
}

func (m *UiStateModel) handleCourseDetailKeys(msg tea.KeyMsg) tea.Cmd {
	data := m.getCurrentCourseData()

	switch msg.String() {
	case "esc", "backspace":
		m.State = viewCourses
		return nil

	case "q":
		return tea.Quit

	case "1":
		m.activeTab = tabAssignments
		m.updateDetailViewport()
	case "2":
		m.activeTab = tabMaterials
		m.updateDetailViewport()
	case "3":
		m.activeTab = tabAnnouncements
		m.updateDetailViewport()
	case "4":
		m.activeTab = tabInfo
		m.updateDetailViewport()

	case "tab", "l", "right":
		m.activeTab = (m.activeTab + 1) % 4
		m.updateDetailViewport()

	case "shift+tab", "h", "left":
		if m.activeTab == 0 {
			m.activeTab = 3
		} else {
			m.activeTab--
		}
		m.updateDetailViewport()

	case "j", "down":
		count := m.getCurrentTabItemCount(data)
		if count > 0 && m.tabCursors[m.activeTab] < count-1 {
			m.tabCursors[m.activeTab]++
			m.updateDetailViewport()
		}

	case "k", "up":
		if m.tabCursors[m.activeTab] > 0 {
			m.tabCursors[m.activeTab]--
			m.updateDetailViewport()
		}

	case "g", "home":
		m.tabCursors[m.activeTab] = 0
		m.tabOffsets[m.activeTab] = 0
		m.updateDetailViewport()

	case "G", "end":
		count := m.getCurrentTabItemCount(data)
		if count > 0 {
			m.tabCursors[m.activeTab] = count - 1
			m.updateDetailViewport()
		}

	case "d":
		m.viewport.HalfViewDown()
	case "u":
		m.viewport.HalfViewUp()

	case "r", "ctrl+r":
		if m.currentCourse != nil {
			m.loading = true
			m.loadingMsg = fmt.Sprintf("Refreshing %s...", m.currentCourse.Name)
			return fetchCourseDataCmd(m.Token, m.currentCourse.Id)
		}

	case "o":
		return m.openCurrentItemWebLink()

	case "a", "enter":
		return m.handleAttachmentAction()
	}

	return nil
}

func (m *UiStateModel) getCurrentCourseData() *CourseDetailData {
	if m.currentCourse == nil {
		return &CourseDetailData{}
	}
	if data, ok := m.courseData[m.currentCourse.Id]; ok {
		return data
	}
	return &CourseDetailData{}
}

func (m *UiStateModel) getCurrentTabItemCount(data *CourseDetailData) int {
	if data == nil {
		return 0
	}
	switch m.activeTab {
	case tabAssignments:
		return len(data.CourseWork)
	case tabMaterials:
		return len(data.Materials)
	case tabAnnouncements:
		return len(data.Announcements)
	case tabInfo:
		return 1
	}
	return 0
}

func (m *UiStateModel) openCurrentItemWebLink() tea.Cmd {
	if m.currentCourse == nil {
		return nil
	}
	data := m.getCurrentCourseData()
	cursor := m.tabCursors[m.activeTab]

	switch m.activeTab {
	case tabAssignments:
		if len(data.CourseWork) > cursor {
			link := data.CourseWork[cursor].AlternateLink
			if link != "" {
				return openBrowserCmd(link)
			}
		}
	case tabMaterials:
		if len(data.Materials) > cursor {
			link := data.Materials[cursor].AlternateLink
			if link != "" {
				return openBrowserCmd(link)
			}
		}
	case tabAnnouncements:
		if len(data.Announcements) > cursor {
			link := data.Announcements[cursor].AlternateLink
			if link != "" {
				return openBrowserCmd(link)
			}
		}
	case tabInfo:
		if m.currentCourse.AlternateLink != "" {
			return openBrowserCmd(m.currentCourse.AlternateLink)
		}
	}

	return openBrowserCmd(m.currentCourse.AlternateLink)
}

func (m *UiStateModel) handleAttachmentAction() tea.Cmd {
	data := m.getCurrentCourseData()
	cursor := m.tabCursors[m.activeTab]
	var atts []models.Attachment

	switch m.activeTab {
	case tabAssignments:
		if len(data.CourseWork) > cursor {
			atts = data.CourseWork[cursor].GetAttachments()
		}
	case tabMaterials:
		if len(data.Materials) > cursor {
			atts = data.Materials[cursor].GetAttachments()
		}
	case tabAnnouncements:
		if len(data.Announcements) > cursor {
			atts = data.Announcements[cursor].GetAttachments()
		}
	case tabInfo:
		if m.currentCourse != nil && m.currentCourse.AlternateLink != "" {
			return openBrowserCmd(m.currentCourse.AlternateLink)
		}
	}

	if len(atts) == 0 {
		return m.openCurrentItemWebLink()
	} else if len(atts) == 1 {
		return openBrowserCmd(atts[0].URL)
	} else {
		m.pickerAttachments = atts
		m.pickerCursor = 0
		m.showPicker = true
		return nil
	}
}

func (m *UiStateModel) applyCourseFilter() {
	query := strings.TrimSpace(strings.ToLower(m.searchInput.Value()))
	if query == "" {
		m.filteredCourses = m.courses
	} else {
		var filtered []models.CourseModel
		for _, c := range m.courses {
			target := strings.ToLower(c.Name + " " + c.Sub + " " + c.Section + " " + c.Room)
			if strings.Contains(target, query) {
				filtered = append(filtered, c)
			}
		}
		m.filteredCourses = filtered
	}
	m.selectedCourseIdx = 0
	m.courseOffset = 0
}

func (m *UiStateModel) updateViewportDimensions() {
	colWidth := (m.width - 4) / 2
	if colWidth < 25 {
		colWidth = 25
	}
	vpHeight := m.height - 5
	if vpHeight < 4 {
		vpHeight = 4
	}
	m.viewport.Width = colWidth - 4
	m.viewport.Height = vpHeight - 2
}

func (m *UiStateModel) updateDetailViewport() {
	m.updateViewportDimensions()
	data := m.getCurrentCourseData()
	cursor := m.tabCursors[m.activeTab]

	var content strings.Builder
	width := m.viewport.Width

	switch m.activeTab {
	case tabAssignments:
		if len(data.CourseWork) == 0 {
			content.WriteString("No assignments found in this course.")
		} else if cursor < len(data.CourseWork) {
			cw := data.CourseWork[cursor]
			content.WriteString(detailTitleStyle.Render(cw.Title) + "\n\n")
			content.WriteString(metaLabelStyle.Render("Due Date:   ") + metaValueStyle.Render(cw.DueString()) + "\n")
			if cw.MaxPoints > 0 {
				content.WriteString(metaLabelStyle.Render("Points:     ") + metaValueStyle.Render(fmt.Sprintf("%.0f pts", cw.MaxPoints)) + "\n")
			}
			if cw.State != "" {
				content.WriteString(metaLabelStyle.Render("State:      ") + metaValueStyle.Render(cw.State) + "\n")
			}
			if cw.WorkType != "" {
				content.WriteString(metaLabelStyle.Render("Type:       ") + metaValueStyle.Render(cw.WorkType) + "\n")
			}
			content.WriteString("\n" + strings.Repeat("─", max(5, width-2)) + "\n\n")

			if cw.Description != "" {
				content.WriteString(metaLabelStyle.Render("Description:") + "\n")
				content.WriteString(descriptionStyle.Render(cw.Description) + "\n\n")
			}

			atts := cw.GetAttachments()
			if len(atts) > 0 {
				content.WriteString(metaLabelStyle.Render(fmt.Sprintf("Attachments (%d):", len(atts))) + "\n")
				for i, a := range atts {
					content.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, a.String()))
				}
				content.WriteString("\n" + lipgloss.NewStyle().Foreground(colorPrimary).Italic(true).Render("Press 'a' or Enter to open attachments") + "\n")
			}
		}

	case tabMaterials:
		if len(data.Materials) == 0 {
			content.WriteString("No materials posted in this course.")
		} else if cursor < len(data.Materials) {
			mat := data.Materials[cursor]
			content.WriteString(detailTitleStyle.Render(mat.Title) + "\n\n")
			if mat.State != "" {
				content.WriteString(metaLabelStyle.Render("State: ") + metaValueStyle.Render(mat.State) + "\n")
			}
			content.WriteString("\n" + strings.Repeat("─", max(5, width-2)) + "\n\n")

			if mat.Description != "" {
				content.WriteString(metaLabelStyle.Render("Description:") + "\n")
				content.WriteString(descriptionStyle.Render(mat.Description) + "\n\n")
			}

			atts := mat.GetAttachments()
			if len(atts) > 0 {
				content.WriteString(metaLabelStyle.Render(fmt.Sprintf("Attachments (%d):", len(atts))) + "\n")
				for i, a := range atts {
					content.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, a.String()))
				}
				content.WriteString("\n" + lipgloss.NewStyle().Foreground(colorPrimary).Italic(true).Render("Press 'a' or Enter to open attachments") + "\n")
			}
		}

	case tabAnnouncements:
		if len(data.Announcements) == 0 {
			content.WriteString("No announcements in this course.")
		} else if cursor < len(data.Announcements) {
			ann := data.Announcements[cursor]
			content.WriteString(detailTitleStyle.Render("📢 Course Announcement") + "\n\n")
			if ann.CreationTime != "" {
				content.WriteString(metaLabelStyle.Render("Posted: ") + metaValueStyle.Render(ann.CreationTime) + "\n")
			}
			content.WriteString("\n" + strings.Repeat("─", max(5, width-2)) + "\n\n")

			if ann.Text != "" {
				content.WriteString(descriptionStyle.Render(ann.Text) + "\n\n")
			}

			atts := ann.GetAttachments()
			if len(atts) > 0 {
				content.WriteString(metaLabelStyle.Render(fmt.Sprintf("Attachments (%d):", len(atts))) + "\n")
				for i, a := range atts {
					content.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, a.String()))
				}
				content.WriteString("\n" + lipgloss.NewStyle().Foreground(colorPrimary).Italic(true).Render("Press 'a' or Enter to open attachments") + "\n")
			}
		}

	case tabInfo:
		if m.currentCourse != nil {
			c := m.currentCourse
			content.WriteString(detailTitleStyle.Render(c.Name) + "\n\n")
			if c.Section != "" {
				content.WriteString(metaLabelStyle.Render("Section:     ") + metaValueStyle.Render(c.Section) + "\n")
			}
			if c.Sub != "" {
				content.WriteString(metaLabelStyle.Render("Subject:     ") + metaValueStyle.Render(c.Sub) + "\n")
			}
			if c.Room != "" {
				content.WriteString(metaLabelStyle.Render("Room:        ") + metaValueStyle.Render(c.Room) + "\n")
			}
			if c.EnrollmentCode != "" {
				content.WriteString(metaLabelStyle.Render("Join Code:   ") + badgeActiveStyle.Render(c.EnrollmentCode) + "\n")
			}
			if c.CourseState != "" {
				content.WriteString(metaLabelStyle.Render("Status:      ") + metaValueStyle.Render(c.CourseState) + "\n")
			}
			if c.DescriptionHeading != "" {
				content.WriteString(metaLabelStyle.Render("Heading:     ") + metaValueStyle.Render(c.DescriptionHeading) + "\n")
			}
			if c.Description != "" {
				content.WriteString("\n" + metaLabelStyle.Render("Description:") + "\n" + descriptionStyle.Render(c.Description) + "\n")
			}
			if c.AlternateLink != "" {
				content.WriteString("\n" + metaLabelStyle.Render("Classroom Link:") + "\n" + attachmentItemStyle.Render(c.AlternateLink) + "\n")
				content.WriteString("\n" + lipgloss.NewStyle().Foreground(colorPrimary).Italic(true).Render("Press 'o' or Enter to open in browser") + "\n")
			}
		}
	}

	m.viewport.SetContent(content.String())
	m.viewport.GotoTop()
}

func (m UiStateModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	if m.showHelp {
		return renderHelpModal(m.width, m.height)
	}

	var viewStr string
	switch m.State {
	case viewCourses:
		viewStr = m.renderCoursesView()
	case viewCourseDetail:
		viewStr = m.renderCourseDetailView()
	}

	if m.showPicker {
		pickerView := m.renderAttachmentPicker()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, pickerView)
	}

	return viewStr
}

func (m UiStateModel) renderHeader() string {
	title := appTitleStyle.Render("📚 CLASSROOM")
	var breadcrumbs string

	if m.State == viewCourses {
		breadcrumbs = breadcrumbActiveStyle.Render("Courses")
	} else if m.currentCourse != nil {
		tabName := "Assignments"
		switch m.activeTab {
		case tabMaterials:
			tabName = "Materials"
		case tabAnnouncements:
			tabName = "Announcements"
		case tabInfo:
			tabName = "Course Info"
		}
		breadcrumbs = breadcrumbStyle.Render("Courses > ") +
			breadcrumbStyle.Render(m.currentCourse.Name+" > ") +
			breadcrumbActiveStyle.Render(tabName)
	}

	leftHeader := lipgloss.JoinHorizontal(lipgloss.Center, title, breadcrumbs)

	var rightHeader string
	if m.loading {
		rightHeader = lipgloss.NewStyle().Foreground(colorPrimary).Render(m.spinner.View() + " " + m.loadingMsg)
	} else {
		rightHeader = lipgloss.NewStyle().Foreground(colorSubtle).Render("? for help")
	}

	space := m.width - lipgloss.Width(leftHeader) - lipgloss.Width(rightHeader) - 2
	if space < 0 {
		space = 0
	}

	header := lipgloss.JoinHorizontal(lipgloss.Center, leftHeader, strings.Repeat(" ", space), rightHeader)
	return headerBarStyle.Width(m.width).Render(header)
}

func (m UiStateModel) renderFooter() string {
	var statusPart string
	if m.statusMsg != "" {
		if m.statusIsErr {
			statusPart = toastErrorStyle.Render(" " + m.statusMsg + " ")
		} else {
			statusPart = toastSuccessStyle.Render(" " + m.statusMsg + " ")
		}
	}

	var helpPart string
	if m.State == viewCourses {
		helpPart = helpKeyStyle.Render("↑/↓/j/k") + helpDescStyle.Render(" navigate  ") +
			helpKeyStyle.Render("Enter") + helpDescStyle.Render(" open  ") +
			helpKeyStyle.Render("o") + helpDescStyle.Render(" web  ") +
			helpKeyStyle.Render("/") + helpDescStyle.Render(" search  ") +
			helpKeyStyle.Render("r") + helpDescStyle.Render(" refresh  ") +
			helpKeyStyle.Render("q") + helpDescStyle.Render(" quit")
	} else {
		helpPart = helpKeyStyle.Render("1-4/Tab") + helpDescStyle.Render(" tabs  ") +
			helpKeyStyle.Render("↑/↓") + helpDescStyle.Render(" items  ") +
			helpKeyStyle.Render("d/u") + helpDescStyle.Render(" scroll  ") +
			helpKeyStyle.Render("a/Enter") + helpDescStyle.Render(" attach  ") +
			helpKeyStyle.Render("o") + helpDescStyle.Render(" web  ") +
			helpKeyStyle.Render("Esc") + helpDescStyle.Render(" back")
	}

	space := m.width - lipgloss.Width(statusPart) - lipgloss.Width(helpPart) - 2
	if space < 0 {
		space = 0
	}

	footer := lipgloss.JoinHorizontal(lipgloss.Center, statusPart, strings.Repeat(" ", space), helpPart)
	return statusBarStyle.Width(m.width).Render(footer)
}

func (m *UiStateModel) renderCoursesView() string {
	header := m.renderHeader()
	footer := m.renderFooter()

	// Height budget calculation
	headerHeight := 3
	footerHeight := 1
	searchHeight := 0

	var searchBox string
	if m.searching || m.searchInput.Value() != "" {
		searchHeight = 1
		searchBox = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Background(lipgloss.Color("#1F2937")).
			Padding(0, 1).
			Width(m.width).
			Render(m.searchInput.View())
	}

	cardHeight := m.height - headerHeight - footerHeight - searchHeight
	if cardHeight < 6 {
		cardHeight = 6
	}
	innerHeight := cardHeight - 2 // subtracting top & bottom borders

	colWidth := (m.width - 4) / 2
	if colWidth < 25 {
		colWidth = 25
	}

	if len(m.filteredCourses) == 0 {
		emptyMsg := "No courses found."
		if m.loading {
			emptyMsg = m.spinner.View() + " Loading courses from Google Classroom..."
		} else if m.searchInput.Value() != "" {
			emptyMsg = fmt.Sprintf("No courses matching '%s'. Press Esc to clear filter.", m.searchInput.Value())
		}
		emptyBox := lipgloss.NewStyle().
			Foreground(colorSubtle).
			Width(m.width-4).
			Height(cardHeight).
			Padding(2, 4).
			Render(emptyMsg)

		if searchHeight > 0 {
			return lipgloss.JoinVertical(lipgloss.Left, header, searchBox, emptyBox, footer)
		}
		return lipgloss.JoinVertical(lipgloss.Left, header, emptyBox, footer)
	}

	// Each course item takes 2 lines (Line 1: Title, Line 2: Details)
	itemsPerPage := innerHeight / 2
	if itemsPerPage < 1 {
		itemsPerPage = 1
	}

	// Adjust scroll window
	if m.selectedCourseIdx < m.courseOffset {
		m.courseOffset = m.selectedCourseIdx
	}
	if m.selectedCourseIdx >= m.courseOffset+itemsPerPage {
		m.courseOffset = m.selectedCourseIdx - itemsPerPage + 1
	}
	if m.courseOffset < 0 {
		m.courseOffset = 0
	}

	startIdx := m.courseOffset
	endIdx := min(len(m.filteredCourses), startIdx+itemsPerPage)

	var leftList strings.Builder
	for i := startIdx; i < endIdx; i++ {
		c := m.filteredCourses[i]
		isSelected := (i == m.selectedCourseIdx)

		curStr := "  "
		if isSelected {
			curStr = "❯ "
		}

		title := c.Name
		maxTextWidth := colWidth - 6
		if maxTextWidth < 10 {
			maxTextWidth = 10
		}
		if len(title) > maxTextWidth {
			title = title[:maxTextWidth-3] + "..."
		}

		sub := c.Sub
		if sub == "" {
			sub = c.Section
		}
		if sub == "" {
			sub = "Google Classroom"
		}
		if len(sub) > maxTextWidth {
			sub = sub[:maxTextWidth-3] + "..."
		}

		var row string
		if isSelected {
			row = itemSelectedStyle.Width(colWidth - 2).Render(fmt.Sprintf("%s%s\n   %s", curStr, title, sub))
		} else {
			row = itemUnselectedStyle.Width(colWidth - 2).Render(fmt.Sprintf("%s%s\n   %s", curStr, title, itemSubtextStyle.Render(sub)))
		}

		leftList.WriteString(row)
		if i < endIdx-1 {
			leftList.WriteString("\n")
		}
	}

	leftCard := cardStyle.Width(colWidth).Height(cardHeight).Render(leftList.String())

	// Right info card for selected course
	var rightCard string
	if m.selectedCourseIdx < len(m.filteredCourses) {
		c := m.filteredCourses[m.selectedCourseIdx]
		var rightContent strings.Builder
		rightContent.WriteString(detailTitleStyle.Render(c.Name) + "\n\n")
		if c.Section != "" {
			rightContent.WriteString(metaLabelStyle.Render("Section: ") + metaValueStyle.Render(c.Section) + "\n")
		}
		if c.Sub != "" {
			rightContent.WriteString(metaLabelStyle.Render("Subject: ") + metaValueStyle.Render(c.Sub) + "\n")
		}
		if c.Room != "" {
			rightContent.WriteString(metaLabelStyle.Render("Room:    ") + metaValueStyle.Render(c.Room) + "\n")
		}
		if c.EnrollmentCode != "" {
			rightContent.WriteString(metaLabelStyle.Render("Code:    ") + badgeActiveStyle.Render(c.EnrollmentCode) + "\n")
		}
		if c.CourseState != "" {
			rightContent.WriteString(metaLabelStyle.Render("State:   ") + metaValueStyle.Render(c.CourseState) + "\n")
		}

		sepLen := colWidth - 6
		if sepLen > 0 {
			rightContent.WriteString("\n" + strings.Repeat("─", sepLen) + "\n\n")
		}
		rightContent.WriteString(lipgloss.NewStyle().Foreground(colorPrimary).Render("Press Enter to open course dashboard\nPress 'o' to open in Classroom web\nPress '/' to search/filter"))

		rightCard = cardStyle.Width(colWidth).Height(cardHeight).Render(rightContent.String())
	}

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftCard, " ", rightCard)

	if searchHeight > 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, searchBox, body, footer)
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m *UiStateModel) renderCourseDetailView() string {
	header := m.renderHeader()
	footer := m.renderFooter()

	data := m.getCurrentCourseData()

	// Tab Bar
	tabAssignmentsStr := "[1] Assignments"
	if len(data.CourseWork) > 0 {
		tabAssignmentsStr = fmt.Sprintf("[1] Assignments (%d)", len(data.CourseWork))
	}
	tabMaterialsStr := "[2] Materials"
	if len(data.Materials) > 0 {
		tabMaterialsStr = fmt.Sprintf("[2] Materials (%d)", len(data.Materials))
	}
	tabAnnounceStr := "[3] Announcements"
	if len(data.Announcements) > 0 {
		tabAnnounceStr = fmt.Sprintf("[3] Announcements (%d)", len(data.Announcements))
	}
	tabInfoStr := "[4] Course Info"

	tabs := []struct {
		name string
		tab  TabType
	}{
		{tabAssignmentsStr, tabAssignments},
		{tabMaterialsStr, tabMaterials},
		{tabAnnounceStr, tabAnnouncements},
		{tabInfoStr, tabInfo},
	}

	var tabBar strings.Builder
	for _, t := range tabs {
		if m.activeTab == t.tab {
			tabBar.WriteString(tabActiveStyle.Render("● " + t.name))
		} else {
			tabBar.WriteString(tabInactiveStyle.Render("○ " + t.name))
		}
	}
	tabBarLine := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color("#111827")).Render(tabBar.String())

	// Height budget calculation
	headerHeight := 3
	tabBarHeight := 1
	footerHeight := 1
	cardHeight := m.height - headerHeight - tabBarHeight - footerHeight
	if cardHeight < 6 {
		cardHeight = 6
	}
	innerHeight := cardHeight - 2

	colWidth := (m.width - 4) / 2
	if colWidth < 25 {
		colWidth = 25
	}

	// Update viewport dimensions
	m.viewport.Width = colWidth - 4
	m.viewport.Height = innerHeight

	// Calculate pagination for left list
	itemsPerPage := innerHeight / 2
	if itemsPerPage < 1 {
		itemsPerPage = 1
	}

	cursor := m.tabCursors[m.activeTab]
	offset := m.tabOffsets[m.activeTab]

	if cursor < offset {
		offset = cursor
	}
	if cursor >= offset+itemsPerPage {
		offset = cursor - itemsPerPage + 1
	}
	if offset < 0 {
		offset = 0
	}
	m.tabOffsets[m.activeTab] = offset

	var leftList strings.Builder

	switch m.activeTab {
	case tabAssignments:
		if len(data.CourseWork) == 0 {
			if m.loading {
				leftList.WriteString(m.spinner.View() + " Loading coursework...")
			} else {
				leftList.WriteString("No assignments found.")
			}
		} else {
			startIdx := offset
			endIdx := min(len(data.CourseWork), startIdx+itemsPerPage)

			for i := startIdx; i < endIdx; i++ {
				cw := data.CourseWork[i]
				isSelected := (i == cursor)
				curStr := "  "
				if isSelected {
					curStr = "❯ "
				}
				title := cw.Title
				maxTextWidth := colWidth - 6
				if maxTextWidth < 10 {
					maxTextWidth = 10
				}
				if len(title) > maxTextWidth {
					title = title[:maxTextWidth-3] + "..."
				}

				var badges []string
				if !cw.DueDate.IsZero() {
					badges = append(badges, badgeDueStyle.Render("Due: "+cw.DueDate.String()))
				}
				if cw.MaxPoints > 0 {
					badges = append(badges, badgePointsStyle.Render(fmt.Sprintf("%.0f pts", cw.MaxPoints)))
				}
				if len(cw.Materials) > 0 {
					badges = append(badges, badgeAttachmentsStyle.Render(fmt.Sprintf("📎 %d", len(cw.Materials))))
				}

				badgeLine := strings.Join(badges, " ")

				if isSelected {
					leftList.WriteString(itemSelectedStyle.Width(colWidth - 2).Render(fmt.Sprintf("%s%s\n   %s", curStr, title, badgeLine)))
				} else {
					leftList.WriteString(itemUnselectedStyle.Width(colWidth - 2).Render(fmt.Sprintf("%s%s\n   %s", curStr, title, badgeLine)))
				}
				if i < endIdx-1 {
					leftList.WriteString("\n")
				}
			}
		}

	case tabMaterials:
		if len(data.Materials) == 0 {
			if m.loading {
				leftList.WriteString(m.spinner.View() + " Loading materials...")
			} else {
				leftList.WriteString("No materials found.")
			}
		} else {
			startIdx := offset
			endIdx := min(len(data.Materials), startIdx+itemsPerPage)

			for i := startIdx; i < endIdx; i++ {
				mat := data.Materials[i]
				isSelected := (i == cursor)
				curStr := "  "
				if isSelected {
					curStr = "❯ "
				}
				title := mat.Title
				maxTextWidth := colWidth - 6
				if maxTextWidth < 10 {
					maxTextWidth = 10
				}
				if len(title) > maxTextWidth {
					title = title[:maxTextWidth-3] + "..."
				}

				var badges []string
				if len(mat.Materials) > 0 {
					badges = append(badges, badgeAttachmentsStyle.Render(fmt.Sprintf("📎 %d files", len(mat.Materials))))
				}
				badgeLine := strings.Join(badges, " ")

				if isSelected {
					leftList.WriteString(itemSelectedStyle.Width(colWidth - 2).Render(fmt.Sprintf("%s%s\n   %s", curStr, title, badgeLine)))
				} else {
					leftList.WriteString(itemUnselectedStyle.Width(colWidth - 2).Render(fmt.Sprintf("%s%s\n   %s", curStr, title, badgeLine)))
				}
				if i < endIdx-1 {
					leftList.WriteString("\n")
				}
			}
		}

	case tabAnnouncements:
		if len(data.Announcements) == 0 {
			if m.loading {
				leftList.WriteString(m.spinner.View() + " Loading announcements...")
			} else {
				leftList.WriteString("No announcements found.")
			}
		} else {
			startIdx := offset
			endIdx := min(len(data.Announcements), startIdx+itemsPerPage)

			for i := startIdx; i < endIdx; i++ {
				ann := data.Announcements[i]
				isSelected := (i == cursor)
				curStr := "  "
				if isSelected {
					curStr = "❯ "
				}
				snippet := strings.ReplaceAll(ann.Text, "\n", " ")
				maxTextWidth := colWidth - 6
				if maxTextWidth < 10 {
					maxTextWidth = 10
				}
				if len(snippet) > maxTextWidth {
					snippet = snippet[:maxTextWidth-3] + "..."
				}

				var badges []string
				if len(ann.Materials) > 0 {
					badges = append(badges, badgeAttachmentsStyle.Render(fmt.Sprintf("📎 %d", len(ann.Materials))))
				}
				badgeLine := strings.Join(badges, " ")

				if isSelected {
					leftList.WriteString(itemSelectedStyle.Width(colWidth - 2).Render(fmt.Sprintf("%s%s\n   %s", curStr, snippet, badgeLine)))
				} else {
					leftList.WriteString(itemUnselectedStyle.Width(colWidth - 2).Render(fmt.Sprintf("%s%s\n   %s", curStr, snippet, badgeLine)))
				}
				if i < endIdx-1 {
					leftList.WriteString("\n")
				}
			}
		}

	case tabInfo:
		leftList.WriteString(itemSelectedStyle.Width(colWidth - 2).Render("❯ Course Details & Web Links"))
	}

	leftPane := cardStyle.Width(colWidth).Height(cardHeight).Render(leftList.String())
	rightPane := detailPaneStyle.Width(colWidth).Height(cardHeight).Render(m.viewport.View())

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, " ", rightPane)
	return lipgloss.JoinVertical(lipgloss.Left, header, tabBarLine, body, footer)
}

func (m UiStateModel) renderAttachmentPicker() string {
	var content strings.Builder
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Background(colorPrimaryDim).
		Padding(0, 2).
		Render("📎 Select Attachment to Open in Browser")

	content.WriteString(title + "\n\n")

	for i, att := range m.pickerAttachments {
		cur := "  "
		if i == m.pickerCursor {
			cur = "❯ "
			row := itemSelectedStyle.Render(fmt.Sprintf("%s[%d] %s", cur, i+1, att.String()))
			content.WriteString(row + "\n")
		} else {
			row := itemUnselectedStyle.Render(fmt.Sprintf("%s[%d] %s", cur, i+1, att.String()))
			content.WriteString(row + "\n")
		}
	}

	content.WriteString("\n" + lipgloss.NewStyle().Foreground(colorSubtle).Italic(true).Render("Press Enter to open • Esc to cancel"))

	return modalBoxStyle.Render(content.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
