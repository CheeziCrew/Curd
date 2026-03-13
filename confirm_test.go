package curd

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func testConfirmConfig() ConfirmConfig {
	return ConfirmConfig{
		Question: "Drop uncommitted changes?",
		Caller:   "cleanup",
		Palette:  SwissgitPalette,
	}
}

func TestNewConfirmModel(t *testing.T) {
	m := NewConfirmModel(testConfirmConfig())
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1 (default to No)", m.cursor)
	}
}

func TestConfirmUpdate_Toggle(t *testing.T) {
	m := NewConfirmModel(testConfirmConfig())
	// Default is No (cursor=1), toggle should go to Yes (cursor=0)
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 after toggle", m.cursor)
	}
	// Toggle back to No
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1 after second toggle", m.cursor)
	}
}

func TestConfirmUpdate_Enter(t *testing.T) {
	tests := []struct {
		name      string
		cursor    int
		wantYes   bool
	}{
		{"enter on Yes", 0, true},
		{"enter on No", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewConfirmModel(testConfirmConfig())
			m.cursor = tt.cursor
			_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			if cmd == nil {
				t.Fatal("expected cmd, got nil")
			}
			msg := cmd()
			cm, ok := msg.(ConfirmMsg)
			if !ok {
				t.Fatalf("expected ConfirmMsg, got %T", msg)
			}
			if cm.Confirmed != tt.wantYes {
				t.Errorf("Confirmed = %v, want %v", cm.Confirmed, tt.wantYes)
			}
			if cm.Caller != "cleanup" {
				t.Errorf("Caller = %q, want %q", cm.Caller, "cleanup")
			}
		})
	}
}

func TestConfirmUpdate_QuickKeys(t *testing.T) {
	tests := []struct {
		name    string
		key     rune
		wantYes bool
	}{
		{"y key confirms", 'y', true},
		{"n key denies", 'n', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewConfirmModel(testConfirmConfig())
			_, cmd := m.Update(tea.KeyPressMsg{Code: tt.key})
			if cmd == nil {
				t.Fatal("expected cmd, got nil")
			}
			msg := cmd()
			cm, ok := msg.(ConfirmMsg)
			if !ok {
				t.Fatalf("expected ConfirmMsg, got %T", msg)
			}
			if cm.Confirmed != tt.wantYes {
				t.Errorf("Confirmed = %v, want %v", cm.Confirmed, tt.wantYes)
			}
		})
	}
}

func TestConfirmUpdate_Esc(t *testing.T) {
	m := NewConfirmModel(testConfirmConfig())
	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if cmd == nil {
		t.Fatal("expected cmd, got nil")
	}
	msg := cmd()
	if _, ok := msg.(BackToMenuMsg); !ok {
		t.Fatalf("expected BackToMenuMsg, got %T", msg)
	}
}

func TestConfirmView(t *testing.T) {
	m := NewConfirmModel(testConfirmConfig())
	v := m.View()
	if v == "" {
		t.Error("View() returned empty string")
	}
}

func TestConfirmInit(t *testing.T) {
	m := NewConfirmModel(testConfirmConfig())
	cmd := m.Init()
	if cmd != nil {
		t.Error("Init() should return nil cmd")
	}
}
