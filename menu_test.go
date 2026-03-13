package curd

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func testMenuConfig() MenuConfig {
	return MenuConfig{
		Banner:  []string{"Test App"},
		Tagline: "testing",
		Items: []MenuItem{
			{Icon: "1", Name: "First", Command: "first", Desc: "do first"},
			{Icon: "2", Name: "Second", Command: "second", Desc: "do second"},
			{Icon: "3", Name: "Third", Command: "third", Desc: "do third"},
		},
		Palette: SwissgitPalette,
	}
}

func TestNewMenuModel(t *testing.T) {
	m := NewMenuModel(testMenuConfig())

	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}
	// styles should be set — verify a style renders non-empty output.
	if m.styles.Title.Render("x") == "" {
		t.Error("styles should be initialized, Title renders empty")
	}
	if len(m.config.Items) != 3 {
		t.Errorf("items = %d, want 3", len(m.config.Items))
	}
}

func TestMenuUpdate_KeyDown(t *testing.T) {
	tests := []struct {
		name       string
		presses    int
		wantCursor int
	}{
		{"one down", 1, 1},
		{"two down", 2, 2},
		{"wrap to zero", 3, 0},
		{"wrap plus one", 4, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMenuModel(testMenuConfig())
			for i := 0; i < tt.presses; i++ {
				m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
			}
			if m.cursor != tt.wantCursor {
				t.Errorf("cursor = %d, want %d", m.cursor, tt.wantCursor)
			}
		})
	}
}

func TestMenuUpdate_KeyUp(t *testing.T) {
	tests := []struct {
		name       string
		presses    int
		wantCursor int
	}{
		{"one up wraps to last", 1, 2},
		{"two up", 2, 1},
		{"three up back to start", 3, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMenuModel(testMenuConfig())
			for i := 0; i < tt.presses; i++ {
				m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
			}
			if m.cursor != tt.wantCursor {
				t.Errorf("cursor = %d, want %d", m.cursor, tt.wantCursor)
			}
		})
	}
}

func TestMenuUpdate_Enter(t *testing.T) {
	tests := []struct {
		name        string
		cursor      int
		wantCommand string
	}{
		{"first item", 0, "first"},
		{"second item", 1, "second"},
		{"third item", 2, "third"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMenuModel(testMenuConfig())
			// Move cursor to desired position
			for i := 0; i < tt.cursor; i++ {
				m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
			}
			_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			if cmd == nil {
				t.Fatal("expected a cmd, got nil")
			}
			msg := cmd()
			sel, ok := msg.(MenuSelectionMsg)
			if !ok {
				t.Fatalf("expected MenuSelectionMsg, got %T", msg)
			}
			if sel.Command != tt.wantCommand {
				t.Errorf("command = %q, want %q", sel.Command, tt.wantCommand)
			}
		})
	}
}

func TestMenuView(t *testing.T) {
	m := NewMenuModel(testMenuConfig())
	m.width = 80
	m.height = 24

	v := m.View()
	if v == "" {
		t.Error("View() returned empty string")
	}
}

func TestMenuUpdate_WindowSize(t *testing.T) {
	m := NewMenuModel(testMenuConfig())

	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	if m.width != 120 {
		t.Errorf("width = %d, want 120", m.width)
	}
	if m.height != 40 {
		t.Errorf("height = %d, want 40", m.height)
	}
}

func TestMenuInit(t *testing.T) {
	m := NewMenuModel(testMenuConfig())
	cmd := m.Init()
	if cmd != nil {
		t.Error("Init() should return nil cmd")
	}
}

func TestMenuUpdate_MouseWheel(t *testing.T) {
	tests := []struct {
		name       string
		button     tea.MouseButton
		wantCursor int
	}{
		{"wheel down moves cursor down", tea.MouseWheelDown, 1},
		{"wheel up wraps to last", tea.MouseWheelUp, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMenuModel(testMenuConfig())
			msg := tea.MouseWheelMsg(tea.Mouse{Button: tt.button})
			m, _ = m.Update(msg)
			if m.cursor != tt.wantCursor {
				t.Errorf("cursor = %d, want %d", m.cursor, tt.wantCursor)
			}
		})
	}
}

func TestMenuUpdate_MouseWheelWrap(t *testing.T) {
	m := NewMenuModel(testMenuConfig())
	// Move to last item (index 2)
	for i := 0; i < 2; i++ {
		m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	}
	// Wheel down should wrap to 0
	m, _ = m.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelDown}))
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 after wrap", m.cursor)
	}
}

func TestMenuView_WithTags(t *testing.T) {
	cfg := MenuConfig{
		Banner:  []string{"Logo Line 1", "Logo Line 2", "Logo Line 3"},
		Tagline: "test tagline",
		Items: []MenuItem{
			{Icon: "1", Name: "First", Command: "first", Desc: "desc1", Tag: "mvn"},
			{Icon: "2", Name: "Second", Command: "second", Desc: "desc2"},
			{Icon: "3", Name: "Third", Command: "third", Desc: "desc3", Tag: "git"},
			{Icon: "4", Name: "Fourth", Command: "fourth", Desc: "desc4"},
			{Icon: "5", Name: "Fifth", Command: "fifth", Desc: "desc5"},
		},
		Palette: SwissgitPalette,
	}
	m := NewMenuModel(cfg)
	m.width = 100
	m.height = 40

	v := m.View()
	if v == "" {
		t.Error("View() returned empty string")
	}

	// Move cursor to a non-zero position and render again
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	v2 := m.View()
	if v2 == "" {
		t.Error("View() returned empty string after moving cursor")
	}
	if v == v2 {
		t.Error("View() should differ when cursor moves")
	}
}

func TestMenuUpdate_UnhandledKey(t *testing.T) {
	m := NewMenuModel(testMenuConfig())
	// Send a key that doesn't match any handler (e.g., 'x')
	m2, cmd := m.Update(tea.KeyPressMsg{Code: 'x'})
	if cmd != nil {
		t.Error("unhandled key should not produce a cmd")
	}
	if m2.cursor != m.cursor {
		t.Error("unhandled key should not change cursor")
	}
}
