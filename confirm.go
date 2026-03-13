package curd

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ConfirmMsg is sent when the user answers a yes/no confirmation prompt.
type ConfirmMsg struct {
	Confirmed bool
	Caller    string // optional context for routing in the parent
}

// ConfirmConfig configures the confirmation dialog.
type ConfirmConfig struct {
	Question string  // e.g. "Drop uncommitted changes?"
	Caller   string  // passed through to ConfirmMsg
	Palette  Palette
}

// ConfirmModel is a yes/no confirmation dialog.
type ConfirmModel struct {
	config  ConfirmConfig
	styles  StyleSet
	cursor  int // 0 = yes, 1 = no
}

// NewConfirmModel creates a confirmation dialog.
func NewConfirmModel(cfg ConfirmConfig) ConfirmModel {
	return ConfirmModel{
		config: cfg,
		styles: cfg.Palette.Styles(),
		cursor: 1, // default to "no" for safety
	}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case IsUp(msg), IsDown(msg), IsToggle(msg):
			m.cursor = 1 - m.cursor
		case msg.String() == "y":
			return m, m.confirm(true)
		case msg.String() == "n":
			return m, m.confirm(false)
		case IsEnter(msg):
			return m, m.confirm(m.cursor == 0)
		case IsBack(msg):
			return m, func() tea.Msg { return BackToMenuMsg{} }
		}
	}
	return m, nil
}

func (m ConfirmModel) confirm(yes bool) tea.Cmd {
	return func() tea.Msg {
		return ConfirmMsg{Confirmed: yes, Caller: m.config.Caller}
	}
}

func (m ConfirmModel) View() string {
	var b strings.Builder
	st := m.styles

	b.WriteString(st.Title.Render(m.config.Question))
	b.WriteString("\n\n")

	yes := "  Yes"
	no := "  No"
	if m.cursor == 0 {
		yes = renderOption(st, true, "Yes")
		no = renderOption(st, false, "No")
	} else {
		yes = renderOption(st, false, "Yes")
		no = renderOption(st, true, "No")
	}
	b.WriteString(yes + "\n")
	b.WriteString(no + "\n")

	b.WriteString(RenderHintBar(st, []Hint{
		{Key: "j/k", Desc: "move"},
		{Key: "y/n", Desc: "quick pick"},
		{Key: "enter", Desc: "confirm"},
		{Key: "esc", Desc: "back"},
	}))

	return b.String()
}

func renderOption(st StyleSet, active bool, label string) string {
	if active {
		return fmt.Sprintf("%s %s",
			st.CursorMark.Render("●"),
			lipgloss.NewStyle().Foreground(st.AccentStyle.GetForeground()).Bold(true).Render(label),
		)
	}
	return fmt.Sprintf("%s %s",
		st.UncheckStyle.Render("○"),
		st.Dim.Render(label),
	)
}
