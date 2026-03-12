package curd

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// Hint represents a single key-description pair in a hint bar.
type Hint struct {
	Key  string // Displayed key (e.g. "esc", "j/k", "enter")
	Desc string // What the key does (e.g. "back", "move", "select")
}

// RenderHintBar renders a styled hint bar from a list of hints.
// Keys are rendered in the palette's accent color, descriptions in normal text,
// separated by dimmed dots. The bar has a top border for visual separation.
func RenderHintBar(styles StyleSet, hints []Hint) string {
	keyStyle := lipgloss.NewStyle().Foreground(styles.AccentStyle.GetForeground()).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(styles.Normal.GetForeground())
	sepStyle := lipgloss.NewStyle().Foreground(styles.Dim.GetForeground())
	barStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(styles.Dim.GetForeground()).
		MarginTop(1).
		PaddingLeft(1)

	var parts []string
	for _, h := range hints {
		parts = append(parts, keyStyle.Render(h.Key)+" "+descStyle.Render(h.Desc))
	}

	return barStyle.Render(strings.Join(parts, sepStyle.Render("  ·  ")))
}
