package curd

import (
	"strings"
	"testing"
)

func TestTruncateError(t *testing.T) {
	tests := []struct {
		name   string
		err    string
		maxLen int
		check  func(t *testing.T, got string)
	}{
		{
			name:   "empty string",
			err:    "",
			maxLen: 80,
			check: func(t *testing.T, got string) {
				if got != "" {
					t.Errorf("expected empty, got %q", got)
				}
			},
		},
		{
			name:   "short string no truncation",
			err:    "something went wrong",
			maxLen: 80,
			check: func(t *testing.T, got string) {
				if got != "something went wrong" {
					t.Errorf("expected unchanged string, got %q", got)
				}
			},
		},
		{
			name:   "exact length no truncation",
			err:    strings.Repeat("a", 50),
			maxLen: 50,
			check: func(t *testing.T, got string) {
				if len(got) != 50 {
					t.Errorf("expected length 50, got %d", len(got))
				}
			},
		},
		{
			name:   "long string gets truncated with ellipsis prefix",
			err:    strings.Repeat("x", 100),
			maxLen: 50,
			check: func(t *testing.T, got string) {
				if !strings.HasPrefix(got, "\u2026") {
					t.Error("truncated string should start with ellipsis")
				}
				// The ellipsis is 3 bytes UTF-8, plus (maxLen-1) ASCII bytes.
				wantLen := len("\u2026") + 50 - 1
				if len(got) != wantLen {
					t.Errorf("truncated byte length should be %d, got %d", wantLen, len(got))
				}
			},
		},
		{
			name:   "maxLen floors to 40",
			err:    strings.Repeat("y", 100),
			maxLen: 10,
			check: func(t *testing.T, got string) {
				// maxLen floors to 40; result = ellipsis (3 bytes) + 39 ASCII bytes = 42
				wantLen := len("\u2026") + 40 - 1
				if len(got) != wantLen {
					t.Errorf("expected byte length %d (floored to 40), got %d", wantLen, len(got))
				}
			},
		},
		{
			name:   "extracts tail after last error prefix when string exceeds maxLen",
			err:    strings.Repeat("a", 50) + " error: short tail",
			maxLen: 50,
			check: func(t *testing.T, got string) {
				// Total length > 50, and tail after "error: " is "short tail" (10 chars <= 50).
				if got != "short tail" {
					t.Errorf("expected tail after 'error: ', got %q", got)
				}
			},
		},
		{
			name:   "uses last error prefix when multiple exist",
			err:    strings.Repeat("b", 50) + " error: first error: second",
			maxLen: 50,
			check: func(t *testing.T, got string) {
				// Total > 50, last "error: " tail is "second" (6 chars <= 50).
				if got != "second" {
					t.Errorf("expected tail after last 'error: ', got %q", got)
				}
			},
		},
		{
			name:   "short string with error prefix returns unchanged",
			err:    "cmd error: oops",
			maxLen: 80,
			check: func(t *testing.T, got string) {
				// String is <= maxLen, returned as-is without extraction.
				if got != "cmd error: oops" {
					t.Errorf("short string should be returned unchanged, got %q", got)
				}
			},
		},
		{
			name:   "error tail still too long falls through to ellipsis",
			err:    "error: " + strings.Repeat("z", 100),
			maxLen: 50,
			check: func(t *testing.T, got string) {
				// error tail is 100 chars > maxLen 50, so falls through to ellipsis truncation.
				if !strings.HasPrefix(got, "\u2026") {
					t.Error("should fall through to ellipsis truncation when error tail is too long")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := TruncateError(tc.err, tc.maxLen)
			tc.check(t, got)
		})
	}
}

func TestNewResultModel(t *testing.T) {
	tasks := []RepoTask{
		{Name: "repo-a", Path: "/a", Status: TaskDone, Result: "ok"},
		{Name: "repo-b", Path: "/b", Status: TaskFailed, Error: "fail"},
	}
	m := NewResultModel("Test Op", tasks, SwissgitPalette)

	if m.Title != "Test Op" {
		t.Errorf("expected title %q, got %q", "Test Op", m.Title)
	}
	if len(m.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(m.Tasks))
	}
	if m.Width != 80 {
		t.Errorf("expected default width 80, got %d", m.Width)
	}
}

func TestResultModelViewAllSuccess(t *testing.T) {
	tasks := []RepoTask{
		{Name: "alpha", Path: "/a", Status: TaskDone, Result: "merged"},
		{Name: "beta", Path: "/b", Status: TaskDone, Result: "merged"},
	}
	m := NewResultModel("PR Merge", tasks, SwissgitPalette)
	view := m.View()

	if view == "" {
		t.Fatal("View() should not return empty string")
	}
	if !strings.Contains(view, "alpha") {
		t.Error("view should contain task name 'alpha'")
	}
	if !strings.Contains(view, "beta") {
		t.Error("view should contain task name 'beta'")
	}
	if !strings.Contains(view, "PR Merge") {
		t.Error("view should contain the title")
	}
}

func TestResultModelViewWithFailures(t *testing.T) {
	tasks := []RepoTask{
		{Name: "ok-repo", Path: "/ok", Status: TaskDone},
		{Name: "bad-repo", Path: "/bad", Status: TaskFailed, Error: "push rejected"},
	}
	m := NewResultModel("Push", tasks, RaclettePalette)
	view := m.View()

	if !strings.Contains(view, "bad-repo") {
		t.Error("view should contain failed task name")
	}
	if !strings.Contains(view, "push rejected") {
		t.Error("view should contain error message")
	}
	if !strings.Contains(view, "ok-repo") {
		t.Error("view should contain successful task name")
	}
}

func TestResultModelViewOnlyFailures(t *testing.T) {
	tasks := []RepoTask{
		{Name: "fail-1", Path: "/f1", Status: TaskFailed, Error: "timeout"},
		{Name: "fail-2", Path: "/f2", Status: TaskFailed, Error: "auth error"},
	}
	m := NewResultModel("Deploy", tasks, FonduePalette)
	view := m.View()

	if !strings.Contains(view, "fail-1") {
		t.Error("view should contain fail-1")
	}
	if !strings.Contains(view, "fail-2") {
		t.Error("view should contain fail-2")
	}
}

func TestResultModelViewNarrowWidth(t *testing.T) {
	tasks := []RepoTask{
		{Name: "fail-repo", Path: "/f", Status: TaskFailed, Error: "some error message"},
	}
	m := NewResultModel("Narrow", tasks, SwissgitPalette)
	m.Width = 30 // Width - 10 = 20, which is < 40, triggers floor to 40
	view := m.View()
	if view == "" {
		t.Error("View() with narrow width should not return empty")
	}
	if !strings.Contains(view, "fail-repo") {
		t.Error("view should contain task name")
	}
}

func TestResultModelViewSuccessNoResult(t *testing.T) {
	// Test success task with empty Result string (no extra dim text)
	tasks := []RepoTask{
		{Name: "clean-repo", Path: "/c", Status: TaskDone, Result: ""},
	}
	m := NewResultModel("Test", tasks, SwissgitPalette)
	view := m.View()
	if !strings.Contains(view, "clean-repo") {
		t.Error("view should contain task name")
	}
}

func TestResultModelViewEmpty(t *testing.T) {
	m := NewResultModel("Nothing", nil, GruyerePalette)
	view := m.View()

	// Should still render the summary banner even with no tasks.
	if view == "" {
		t.Error("View() with no tasks should still render summary banner")
	}
}
