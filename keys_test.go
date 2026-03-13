package curd

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// keyMsg builds a KeyPressMsg for a special key (Code-based).
func keyMsg(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code}
}

// textKeyMsg builds a KeyPressMsg for a printable character.
func textKeyMsg(s string) tea.KeyPressMsg {
	r := []rune(s)
	return tea.KeyPressMsg{Code: r[0], Text: s}
}

// modKeyMsg builds a KeyPressMsg with a modifier.
func modKeyMsg(code rune, mod tea.KeyMod) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code, Mod: mod}
}

func TestIsUp(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyPressMsg
		want bool
	}{
		{"arrow up", keyMsg(tea.KeyUp), true},
		{"k key", textKeyMsg("k"), true},
		{"j key", textKeyMsg("j"), false},
		{"enter key", keyMsg(tea.KeyEnter), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsUp(tc.msg); got != tc.want {
				t.Errorf("IsUp(%q) = %v, want %v", tc.msg.String(), got, tc.want)
			}
		})
	}
}

func TestIsDown(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyPressMsg
		want bool
	}{
		{"arrow down", keyMsg(tea.KeyDown), true},
		{"j key", textKeyMsg("j"), true},
		{"k key", textKeyMsg("k"), false},
		{"space key", keyMsg(tea.KeySpace), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsDown(tc.msg); got != tc.want {
				t.Errorf("IsDown(%q) = %v, want %v", tc.msg.String(), got, tc.want)
			}
		})
	}
}

func TestIsEnter(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyPressMsg
		want bool
	}{
		{"enter", keyMsg(tea.KeyEnter), true},
		{"escape", keyMsg(tea.KeyEscape), false},
		{"k key", textKeyMsg("k"), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsEnter(tc.msg); got != tc.want {
				t.Errorf("IsEnter(%q) = %v, want %v", tc.msg.String(), got, tc.want)
			}
		})
	}
}

func TestIsBack(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyPressMsg
		want bool
	}{
		{"escape", keyMsg(tea.KeyEscape), true},
		{"enter", keyMsg(tea.KeyEnter), false},
		{"q key", textKeyMsg("q"), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsBack(tc.msg); got != tc.want {
				t.Errorf("IsBack(%q) = %v, want %v", tc.msg.String(), got, tc.want)
			}
		})
	}
}

func TestIsToggle(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyPressMsg
		want bool
	}{
		{"space", keyMsg(tea.KeySpace), true},
		{"enter", keyMsg(tea.KeyEnter), false},
		{"a key", textKeyMsg("a"), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsToggle(tc.msg); got != tc.want {
				t.Errorf("IsToggle(%q) = %v, want %v", tc.msg.String(), got, tc.want)
			}
		})
	}
}

func TestIsSelectAll(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyPressMsg
		want bool
	}{
		{"ctrl+a", modKeyMsg('a', tea.ModCtrl), true},
		{"plain a", textKeyMsg("a"), false},
		{"space", keyMsg(tea.KeySpace), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsSelectAll(tc.msg); got != tc.want {
				t.Errorf("IsSelectAll(%q) = %v, want %v", tc.msg.String(), got, tc.want)
			}
		})
	}
}
