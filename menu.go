package curd

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/harmonica"
)

// MenuItem defines a single menu entry.
type MenuItem struct {
	Icon    string
	Name    string
	Command string
	Desc    string
	Tag     string // optional tag like "[mvn]" shown between name and desc
}

// MenuConfig configures the menu appearance.
type MenuConfig struct {
	Banner  []string // ASCII art lines for the logo
	Tagline string   // shown below logo
	Items   []MenuItem
	Palette Palette
}

// MenuModel is the shared main menu component.
type MenuModel struct {
	config  MenuConfig
	styles  StyleSet
	cursor  int
	width   int
	height  int

	// Harmonica spring for smooth cursor animation.
	spring  harmonica.Spring
	springY float64
	springV float64
	targetY float64
}

// NewMenuModel creates a menu from config.
func NewMenuModel(cfg MenuConfig) MenuModel {
	return MenuModel{
		config:  cfg,
		styles:  cfg.Palette.Styles(),
		spring:  harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.8),
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		switch {
		case IsUp(msg):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.config.Items) - 1
			}
			m.targetY = float64(m.cursor)
		case IsDown(msg):
			m.cursor++
			if m.cursor >= len(m.config.Items) {
				m.cursor = 0
			}
			m.targetY = float64(m.cursor)
		case IsEnter(msg):
			return m, func() tea.Msg {
				return MenuSelectionMsg{Command: m.config.Items[m.cursor].Command}
			}
		}

	case tea.MouseWheelMsg:
		mouse := msg.Mouse()
		if mouse.Button == tea.MouseWheelUp {
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.config.Items) - 1
			}
		} else if mouse.Button == tea.MouseWheelDown {
			m.cursor++
			if m.cursor >= len(m.config.Items) {
				m.cursor = 0
			}
		}
		m.targetY = float64(m.cursor)
	}

	// Tick spring animation.
	m.springY, m.springV = m.spring.Update(m.springY, m.springV, m.targetY)

	return m, nil
}

func (m MenuModel) View() string {
	var s strings.Builder
	st := m.styles

	// Logo banner with gradient
	var logoLines []string
	for i, line := range m.config.Banner {
		ci := i % len(m.config.Palette.LogoGradient)
		logoLines = append(logoLines, lipgloss.NewStyle().Foreground(m.config.Palette.LogoGradient[ci]).Bold(true).Render(line))
	}
	logoContent := strings.Join(logoLines, "\n")
	logoContent += "\n" + st.Tagline.Render("  "+m.config.Tagline) + "  " + st.Version.Render("⚙")
	s.WriteString(st.LogoBox.Render(logoContent))
	s.WriteString("\n\n")

	// Menu items
	for i, item := range m.config.Items {
		tag := ""
		if item.Tag != "" {
			tag = st.Dim.Render(" ["+item.Tag+"]") + " "
		}

		if i == m.cursor {
			line := fmt.Sprintf("%s  %s %s %s", item.Icon, st.MenuActiveName.Render(item.Name), tag, st.MenuActiveDesc.Render(item.Desc))
			s.WriteString(st.MenuActiveItem.Render(line))
		} else {
			line := fmt.Sprintf("%s  %s %s %s", item.Icon, st.MenuInactiveName.Render(item.Name), tag, st.MenuInactiveDesc.Render(item.Desc))
			s.WriteString(st.MenuInactiveItem.Render(line))
		}
		s.WriteString("\n")
	}

	s.WriteString(RenderHintBar(st, []Hint{
		{Key: "j/k", Desc: "move"},
		{Key: "enter", Desc: "select"},
		{Key: "q", Desc: "quit"},
	}))

	return s.String()
}
