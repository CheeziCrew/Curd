package curd

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// FilterListConfig configures a filterable list.
type FilterListConfig struct {
	Items    []list.Item
	Delegate list.ItemDelegate
	Title    string
	Palette  Palette
	Width    int
	Height   int
}

// FilterListModel wraps bubbles/list with palette-aware styling.
type FilterListModel struct {
	List   list.Model // exposed so apps can call List methods directly
	styles StyleSet
}

// NewFilterListModel creates a palette-styled filterable list.
func NewFilterListModel(cfg FilterListConfig) FilterListModel {
	st := cfg.Palette.Styles()

	l := list.New(cfg.Items, cfg.Delegate, cfg.Width, cfg.Height)
	l.Title = cfg.Title
	l.SetShowHelp(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	// Apply palette styles to list chrome.
	ls := list.DefaultStyles(true)
	ls.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(cfg.Palette.AccentBright).
		MarginBottom(1)
	ls.Filter.Focused.Prompt = lipgloss.NewStyle().
		Foreground(cfg.Palette.AccentBright).
		Bold(true)
	ls.StatusBar = lipgloss.NewStyle().
		Foreground(ColorGray)
	ls.NoItems = lipgloss.NewStyle().
		Foreground(ColorGray).
		Italic(true)
	l.Styles = ls

	return FilterListModel{
		List:   l,
		styles: st,
	}
}

func (m FilterListModel) Init() tea.Cmd {
	return nil
}

func (m FilterListModel) Update(msg tea.Msg) (FilterListModel, tea.Cmd) {
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m FilterListModel) View() string {
	v := m.List.View()
	hints := []Hint{
		{Key: "/", Desc: "filter"},
		{Key: "j/k", Desc: "move"},
		{Key: "enter", Desc: "select"},
		{Key: "esc", Desc: "back"},
	}
	return v + "\n" + RenderHintBar(m.styles, hints)
}

// SelectedItem returns the currently highlighted item, or nil if the list is empty.
func (m FilterListModel) SelectedItem() list.Item {
	return m.List.SelectedItem()
}

// SetSize updates the list dimensions.
func (m *FilterListModel) SetSize(w, h int) {
	m.List.SetSize(w, h)
}

// Filtering returns true if the user is currently typing a filter query.
func (m FilterListModel) Filtering() bool {
	return m.List.FilterState() == list.Filtering
}
