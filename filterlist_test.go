package curd

import (
	"fmt"
	"io"
	"testing"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// testItem is a minimal list.Item for testing.
type testItem struct{ name string }

func (i testItem) FilterValue() string { return i.name }

// testDelegate is a minimal delegate for testing.
type testDelegate struct{}

func (d testDelegate) Height() int                             { return 1 }
func (d testDelegate) Spacing() int                            { return 0 }
func (d testDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d testDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	ti := item.(testItem)
	fmt.Fprint(w, ti.name)
}

func testFilterListConfig() FilterListConfig {
	items := []list.Item{
		testItem{"alpha"},
		testItem{"beta"},
		testItem{"gamma"},
	}
	return FilterListConfig{
		Items:    items,
		Delegate: testDelegate{},
		Title:    "Pick one",
		Palette:  FonduePalette,
		Width:    80,
		Height:   24,
	}
}

func TestNewFilterListModel(t *testing.T) {
	m := NewFilterListModel(testFilterListConfig())

	if m.List.Title != "Pick one" {
		t.Errorf("Title = %q, want %q", m.List.Title, "Pick one")
	}
	if len(m.List.Items()) != 3 {
		t.Errorf("items = %d, want 3", len(m.List.Items()))
	}
}

func TestFilterListModel_SelectedItem(t *testing.T) {
	m := NewFilterListModel(testFilterListConfig())
	item := m.SelectedItem()
	if item == nil {
		t.Fatal("expected selected item, got nil")
	}
	ti, ok := item.(testItem)
	if !ok {
		t.Fatalf("expected testItem, got %T", item)
	}
	if ti.name != "alpha" {
		t.Errorf("selected = %q, want %q", ti.name, "alpha")
	}
}

func TestFilterListModel_SetSize(t *testing.T) {
	m := NewFilterListModel(testFilterListConfig())
	m.SetSize(120, 40)
	if m.List.Width() != 120 {
		t.Errorf("width = %d, want 120", m.List.Width())
	}
}

func TestFilterListModel_View(t *testing.T) {
	m := NewFilterListModel(testFilterListConfig())
	v := m.View()
	if v == "" {
		t.Error("View() returned empty string")
	}
}

func TestFilterListModel_Init(t *testing.T) {
	m := NewFilterListModel(testFilterListConfig())
	cmd := m.Init()
	if cmd != nil {
		t.Error("Init() should return nil cmd")
	}
}

func TestFilterListModel_Filtering(t *testing.T) {
	m := NewFilterListModel(testFilterListConfig())
	if m.Filtering() {
		t.Error("should not be filtering initially")
	}
}
