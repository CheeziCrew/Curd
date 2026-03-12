package curd

import (
	"fmt"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ScanFunc is the scanning strategy. Apps provide their own.
// It runs in a goroutine and returns discovered repos.
type ScanFunc func(rootPath string) ([]RepoInfo, error)

// reposScanResultMsg is sent after scanning completes.
type reposScanResultMsg struct {
	repos []RepoInfo
}

// RepoSelectConfig configures the repo selector.
type RepoSelectConfig struct {
	Palette      Palette
	RootPath     string
	Caller       string // which screen requested the selection
	ParentOffset int    // lines consumed by parent screen
	TermHeight   int
	Scanner      ScanFunc
	SingleSelect bool // enter picks cursor item directly, no toggle/select-all
}

// RepoSelectModel lets the user pick repos from discovered subdirectories.
type RepoSelectModel struct {
	repos        []RepoInfo
	cursor       int
	selected     map[int]bool
	caller       string
	loading      bool
	spinner      spinner.Model
	rootPath     string
	termHeight   int
	winOffset    int
	parentOffset int
	styles       StyleSet
	scanFunc     ScanFunc
	singleSelect bool
}

func measureRepoChrome(_ StyleSet, scrollable bool) int {
	header := lipgloss.NewStyle().Bold(true).Render("Select Repos") + "\n\n"
	chrome := lipgloss.Height(header)
	if scrollable {
		chrome += 2
	}
	chrome += 3 // hint bar: margin-top + border + content
	return chrome
}

// NewRepoSelectModel creates a repo selector.
func NewRepoSelectModel(cfg RepoSelectConfig) RepoSelectModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(cfg.Palette.Accent)

	return RepoSelectModel{
		selected:     make(map[int]bool),
		caller:       cfg.Caller,
		loading:      true,
		spinner:      s,
		rootPath:     cfg.RootPath,
		parentOffset: cfg.ParentOffset,
		termHeight:   cfg.TermHeight,
		styles:       cfg.Palette.Styles(),
		scanFunc:     cfg.Scanner,
		singleSelect: cfg.SingleSelect,
	}
}

func (m RepoSelectModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.scanRepos())
}

func (m RepoSelectModel) scanRepos() tea.Cmd {
	root := m.rootPath
	scan := m.scanFunc
	return func() tea.Msg {
		if root == "" {
			root = "."
		}
		repos, err := scan(root)
		if err != nil {
			return reposScanResultMsg{repos: nil}
		}
		return reposScanResultMsg{repos: repos}
	}
}

func (m *RepoSelectModel) visibleRepoCount() int {
	scrollable := len(m.repos) > 10
	chrome := measureRepoChrome(m.styles, scrollable)
	wh := m.termHeight - m.parentOffset - chrome - 1
	if wh < 5 {
		wh = 5
	}
	return wh
}

func (m *RepoSelectModel) ensureCursorVisible() {
	wh := m.visibleRepoCount()
	if m.cursor < m.winOffset {
		m.winOffset = m.cursor
	}
	if m.cursor >= m.winOffset+wh {
		m.winOffset = m.cursor - wh + 1
	}
}

// SelectedPaths returns paths of all currently selected repos.
func (m RepoSelectModel) SelectedPaths() []string {
	var paths []string
	for i, r := range m.repos {
		if m.selected[i] {
			paths = append(paths, r.Path)
		}
	}
	return paths
}

func (m RepoSelectModel) Update(msg tea.Msg) (RepoSelectModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termHeight = msg.Height
		m.ensureCursorVisible()

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case reposScanResultMsg:
		return m.handleScanResult(msg), nil

	case tea.MouseWheelMsg:
		m.handleMouseWheel(msg)

	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	}
	return m, nil
}

func (m *RepoSelectModel) handleScanResult(msg reposScanResultMsg) RepoSelectModel {
	m.loading = false
	m.repos = msg.repos
	for i, r := range m.repos {
		if r.IsDirty {
			m.selected[i] = true
		}
	}
	return *m
}

func (m *RepoSelectModel) handleMouseWheel(msg tea.MouseWheelMsg) {
	if m.loading || len(m.repos) == 0 {
		return
	}
	mouse := msg.Mouse()
	if mouse.Button == tea.MouseWheelUp {
		m.moveCursorUp()
	} else if mouse.Button == tea.MouseWheelDown {
		m.moveCursorDown()
	}
}

func (m *RepoSelectModel) moveCursorUp() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.repos) - 1
		wh := m.visibleRepoCount()
		m.winOffset = len(m.repos) - wh
		if m.winOffset < 0 {
			m.winOffset = 0
		}
	}
	m.ensureCursorVisible()
}

func (m *RepoSelectModel) moveCursorDown() {
	m.cursor++
	if m.cursor >= len(m.repos) {
		m.cursor = 0
		m.winOffset = 0
	}
	m.ensureCursorVisible()
}

func (m RepoSelectModel) handleKeyPress(msg tea.KeyPressMsg) (RepoSelectModel, tea.Cmd) {
	if m.loading {
		return m, nil
	}

	switch {
	case IsUp(msg):
		m.moveCursorUp()
	case IsDown(msg):
		m.moveCursorDown()
	case IsToggle(msg):
		m.handleToggle()
	case IsSelectAll(msg):
		m.handleSelectAll()
	case IsEnter(msg):
		return m.handleEnter()
	case IsBack(msg):
		return m, func() tea.Msg { return BackToMenuMsg{} }
	}
	return m, nil
}

func (m *RepoSelectModel) handleToggle() {
	if !m.singleSelect {
		m.selected[m.cursor] = !m.selected[m.cursor]
	}
}

func (m *RepoSelectModel) handleSelectAll() {
	if m.singleSelect {
		return
	}
	allSelected := true
	for i := range m.repos {
		if !m.selected[i] {
			allSelected = false
			break
		}
	}
	if allSelected {
		m.selected = make(map[int]bool)
	} else {
		for i := range m.repos {
			m.selected[i] = true
		}
	}
}

func (m RepoSelectModel) handleEnter() (RepoSelectModel, tea.Cmd) {
	caller := m.caller
	if m.singleSelect {
		if m.cursor < len(m.repos) {
			path := m.repos[m.cursor].Path
			return m, func() tea.Msg {
				return RepoSelectDoneMsg{Paths: []string{path}, Caller: caller}
			}
		}
	} else {
		var paths []string
		for i, r := range m.repos {
			if m.selected[i] {
				paths = append(paths, r.Path)
			}
		}
		return m, func() tea.Msg {
			return RepoSelectDoneMsg{Paths: paths, Caller: caller}
		}
	}
	return m, nil
}

func (m RepoSelectModel) View() string {
	st := m.styles

	if m.loading {
		content := fmt.Sprintf("%s Scanning repos…", m.spinner.View())
		return st.InputBox.Render(content)
	}

	if len(m.repos) == 0 {
		return st.InputBox.Render(st.Dim.Render("No git repositories found."))
	}

	wh := m.visibleRepoCount()
	scrollable := len(m.repos) > wh

	visibleStart := m.winOffset
	visibleEnd := visibleStart + wh
	if visibleEnd > len(m.repos) {
		visibleEnd = len(m.repos)
	}
	if visibleStart > len(m.repos) {
		visibleStart = 0
	}

	var s string
	s += m.viewHeader(st)
	s += m.viewScrollUp(st, visibleStart, scrollable)
	s += m.viewRepoList(st, visibleStart, visibleEnd)
	s += m.viewScrollDown(st, visibleEnd, scrollable)
	s += m.viewHintBar(st)

	return s
}

func (m RepoSelectModel) viewHeader(st StyleSet) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(st.Title.GetForeground())
	if m.singleSelect {
		counter := st.Dim.Render(fmt.Sprintf("%d repos", len(m.repos)))
		return titleStyle.Render("Select Repo") + "  " + counter + "\n\n"
	}
	selCount := 0
	for i := range m.repos {
		if m.selected[i] {
			selCount++
		}
	}
	counter := st.Dim.Render(fmt.Sprintf("%d/%d selected", selCount, len(m.repos)))
	return titleStyle.Render("Select Repos") + "  " + counter + "\n\n"
}

func (m RepoSelectModel) viewScrollUp(st StyleSet, visibleStart int, scrollable bool) string {
	if visibleStart > 0 {
		return st.Help.Render(fmt.Sprintf("  ↑ %d more above", visibleStart)) + "\n"
	}
	if scrollable {
		return "\n"
	}
	return ""
}

func (m RepoSelectModel) viewRepoList(st StyleSet, visibleStart, visibleEnd int) string {
	var s string
	for i := visibleStart; i < visibleEnd; i++ {
		s += m.viewRepoLine(st, i) + "\n"
	}
	return s
}

func (m RepoSelectModel) viewRepoLine(st StyleSet, i int) string {
	r := m.repos[i]
	isCursor := i == m.cursor
	isSelected := m.selected[i]

	prefix := m.repoPrefix(st, isCursor, isSelected)
	name := m.repoName(st, r.Name, isCursor, isSelected)
	info := m.repoInfo(st, r, isCursor, isSelected)
	branchInfo := m.repoBranch(st, r)

	line := fmt.Sprintf("%s  %s%s%s", prefix, name, info, branchInfo)
	if isCursor {
		return st.RepoActiveItem.Render(line)
	}
	return st.RepoInactiveItem.Render(line)
}

func (m RepoSelectModel) repoPrefix(st StyleSet, isCursor, isSelected bool) string {
	if m.singleSelect {
		if isCursor {
			return st.CursorMark.Render("▸")
		}
		return " "
	}
	if isSelected {
		return st.CheckStyle.Render("●")
	}
	return st.UncheckStyle.Render("○")
}

func (m RepoSelectModel) repoName(st StyleSet, name string, isCursor, isSelected bool) string {
	if isCursor {
		return st.RepoCursorName.Render(name)
	}
	if isSelected {
		return st.RepoSelectedName.Render(name)
	}
	return st.RepoUnselectedName.Render(name)
}

func (m RepoSelectModel) repoInfo(st StyleSet, r RepoInfo, isCursor, isSelected bool) string {
	totalChanges := r.Modified + r.Added + r.Deleted + r.Untracked
	if totalChanges > 0 {
		return st.DirtyStyle.Render(fmt.Sprintf(" %dΔ", totalChanges))
	}
	if isSelected || isCursor {
		return st.CleanMark.Render(" ✓")
	}
	return ""
}

func (m RepoSelectModel) repoBranch(st StyleSet, r RepoInfo) string {
	if r.Branch != "" && r.Branch != r.DefaultBranch {
		return st.BranchMark.Render(fmt.Sprintf(" (%s)", r.Branch))
	}
	return ""
}

func (m RepoSelectModel) viewScrollDown(st StyleSet, visibleEnd int, scrollable bool) string {
	remaining := len(m.repos) - visibleEnd
	if remaining > 0 {
		return st.Help.Render(fmt.Sprintf("  ↓ %d more below", remaining)) + "\n"
	}
	if scrollable {
		return "\n"
	}
	return ""
}

func (m RepoSelectModel) viewHintBar(st StyleSet) string {
	if m.singleSelect {
		return RenderHintBar(st, []Hint{
			{Key: "enter", Desc: "select"},
			{Key: "esc", Desc: "back"},
		})
	}
	return RenderHintBar(st, []Hint{
		{Key: "space", Desc: "toggle"},
		{Key: "ctrl+a", Desc: "all"},
		{Key: "enter", Desc: "confirm"},
		{Key: "esc", Desc: "back"},
	})
}
