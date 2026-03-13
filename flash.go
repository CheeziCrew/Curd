package curd

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// FlashClearMsg is sent when a flash message should be cleared.
type FlashClearMsg struct{}

// FlashModel displays a timed notification message that auto-clears.
type FlashModel struct {
	message string
	isError bool
	styles  StyleSet
}

// NewFlashModel creates a flash notification tied to a palette.
func NewFlashModel(p Palette) FlashModel {
	return FlashModel{styles: p.Styles()}
}

// Show sets a success flash message and returns a cmd to clear it after duration.
func (m FlashModel) Show(message string, duration time.Duration) (FlashModel, tea.Cmd) {
	m.message = message
	m.isError = false
	return m, tea.Tick(duration, func(time.Time) tea.Msg { return FlashClearMsg{} })
}

// ShowError sets an error flash message and returns a cmd to clear it after duration.
func (m FlashModel) ShowError(message string, duration time.Duration) (FlashModel, tea.Cmd) {
	m.message = message
	m.isError = true
	return m, tea.Tick(duration, func(time.Time) tea.Msg { return FlashClearMsg{} })
}

// Clear removes the current flash message.
func (m FlashModel) Clear() FlashModel {
	m.message = ""
	return m
}

// View renders the flash message. Returns empty string if no message is active.
func (m FlashModel) View() string {
	if m.message == "" {
		return ""
	}
	if m.isError {
		return lipgloss.NewStyle().
			Foreground(m.styles.FailStyle.GetForeground()).
			Bold(true).
			PaddingLeft(2).
			Render("✗ " + m.message)
	}
	return lipgloss.NewStyle().
		Foreground(m.styles.SuccessStyle.GetForeground()).
		Bold(true).
		PaddingLeft(2).
		Render("✓ " + m.message)
}

// Active returns true if a flash message is currently displayed.
func (m FlashModel) Active() bool {
	return m.message != ""
}
