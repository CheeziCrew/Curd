package curd

import (
	"testing"
	"time"
)

func TestNewFlashModel(t *testing.T) {
	m := NewFlashModel(SwissgitPalette)
	if m.Active() {
		t.Error("new flash should not be active")
	}
	if m.View() != "" {
		t.Error("inactive flash should render empty string")
	}
}

func TestFlashModel_Show(t *testing.T) {
	m := NewFlashModel(SwissgitPalette)
	m, cmd := m.Show("saved", 3*time.Second)

	if !m.Active() {
		t.Error("flash should be active after Show")
	}
	if cmd == nil {
		t.Error("Show should return a tick cmd")
	}

	v := m.View()
	if v == "" {
		t.Error("active flash should render non-empty")
	}
}

func TestFlashModel_ShowError(t *testing.T) {
	m := NewFlashModel(SwissgitPalette)
	m, cmd := m.ShowError("failed", 3*time.Second)

	if !m.Active() {
		t.Error("flash should be active after ShowError")
	}
	if cmd == nil {
		t.Error("ShowError should return a tick cmd")
	}
	if !m.isError {
		t.Error("isError should be true after ShowError")
	}
}

func TestFlashModel_Clear(t *testing.T) {
	m := NewFlashModel(SwissgitPalette)
	m, _ = m.Show("test", 3*time.Second)

	m = m.Clear()
	if m.Active() {
		t.Error("flash should not be active after Clear")
	}
	if m.View() != "" {
		t.Error("cleared flash should render empty")
	}
}
