package curd

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

// NewStyledInput creates a textinput with palette-consistent styling.
// Width defaults to 60 characters.
func NewStyledInput(placeholder string, p Palette) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetWidth(60)
	s := textinput.DefaultStyles(true)
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(p.AccentBright)
	s.Focused.Text = lipgloss.NewStyle().Foreground(ColorFg)
	s.Blurred.Prompt = lipgloss.NewStyle().Foreground(ColorGray)
	s.Blurred.Text = lipgloss.NewStyle().Foreground(ColorGray)
	ti.SetStyles(s)
	return ti
}
