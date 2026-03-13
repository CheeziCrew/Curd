package curd

import (
	"strings"
	"testing"
)

func TestRenderHintBarEmpty(t *testing.T) {
	ss := SwissgitPalette.Styles()
	out := RenderHintBar(ss, nil)
	// Even with no hints, the bar frame should still render.
	if out == "" {
		t.Error("RenderHintBar with nil hints should return non-empty (bar frame)")
	}
}

func TestRenderHintBarSingle(t *testing.T) {
	ss := SwissgitPalette.Styles()
	hints := []Hint{{Key: "esc", Desc: "back"}}
	out := RenderHintBar(ss, hints)

	if !strings.Contains(out, "esc") {
		t.Error("output should contain the key text 'esc'")
	}
	if !strings.Contains(out, "back") {
		t.Error("output should contain the desc text 'back'")
	}
}

func TestRenderHintBarMultiple(t *testing.T) {
	ss := RaclettePalette.Styles()
	hints := []Hint{
		{Key: "j/k", Desc: "move"},
		{Key: "enter", Desc: "select"},
		{Key: "q", Desc: "quit"},
	}
	out := RenderHintBar(ss, hints)

	for _, h := range hints {
		if !strings.Contains(out, h.Key) {
			t.Errorf("output should contain key %q", h.Key)
		}
		if !strings.Contains(out, h.Desc) {
			t.Errorf("output should contain desc %q", h.Desc)
		}
	}
}

func TestRenderHintBarSeparator(t *testing.T) {
	ss := FonduePalette.Styles()
	hints := []Hint{
		{Key: "a", Desc: "first"},
		{Key: "b", Desc: "second"},
	}
	out := RenderHintBar(ss, hints)

	// The separator dot should appear between items.
	if !strings.Contains(out, "\u00b7") { // middle dot
		t.Error("output should contain the separator dot between hints")
	}
}

func TestRenderHintBarDifferentPalettes(t *testing.T) {
	hints := []Hint{{Key: "x", Desc: "do"}}
	sg := RenderHintBar(SwissgitPalette.Styles(), hints)
	ra := RenderHintBar(RaclettePalette.Styles(), hints)

	// Different palettes should produce different styled output.
	if sg == ra {
		t.Error("hint bars with different palettes should differ in output")
	}
}
