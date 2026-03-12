package curd

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// KeyMap defines shared keybindings across all swiss* apps.
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Quit   key.Binding
	Toggle key.Binding
	All    key.Binding
	Tab    key.Binding
}

// DefaultKeyMap provides standard keybindings.
var DefaultKeyMap = KeyMap{
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Enter:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Toggle: key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "toggle")),
	All:    key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("ctrl+a", "toggle all")),
	Tab:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
}

// Convenience helpers — keeps Update methods clean.

func IsUp(msg tea.KeyPressMsg) bool {
	return key.Matches(msg, DefaultKeyMap.Up)
}

func IsDown(msg tea.KeyPressMsg) bool {
	return key.Matches(msg, DefaultKeyMap.Down)
}

func IsEnter(msg tea.KeyPressMsg) bool {
	return key.Matches(msg, DefaultKeyMap.Enter)
}

func IsBack(msg tea.KeyPressMsg) bool {
	return key.Matches(msg, DefaultKeyMap.Back)
}

func IsToggle(msg tea.KeyPressMsg) bool {
	return key.Matches(msg, DefaultKeyMap.Toggle)
}

func IsSelectAll(msg tea.KeyPressMsg) bool {
	return key.Matches(msg, DefaultKeyMap.All)
}
