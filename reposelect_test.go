package curd

import (
	"fmt"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func testRepoSelectModel() RepoSelectModel {
	return RepoSelectModel{
		repos: []RepoInfo{
			{Path: "/a", Name: "repo-a"},
			{Path: "/b", Name: "repo-b"},
			{Path: "/c", Name: "repo-c"},
		},
		selected:   map[int]bool{0: true},
		caller:     "test",
		styles:     SwissgitPalette.Styles(),
		termHeight: 50,
	}
}

func TestSelectedPaths(t *testing.T) {
	m := testRepoSelectModel()
	m.selected = map[int]bool{0: true, 2: true}

	paths := m.SelectedPaths()

	if len(paths) != 2 {
		t.Fatalf("len(paths) = %d, want 2", len(paths))
	}
	want := map[string]bool{"/a": true, "/c": true}
	for _, p := range paths {
		if !want[p] {
			t.Errorf("unexpected path %q", p)
		}
	}
}

func TestSelectedPaths_Empty(t *testing.T) {
	m := testRepoSelectModel()
	m.selected = map[int]bool{}

	paths := m.SelectedPaths()
	if len(paths) != 0 {
		t.Errorf("len(paths) = %d, want 0", len(paths))
	}
}

func TestVisibleRepoCount(t *testing.T) {
	m := testRepoSelectModel()
	m.termHeight = 50

	count := m.visibleRepoCount()

	if count < 5 {
		t.Errorf("visibleRepoCount() = %d, want >= 5", count)
	}
	if count > m.termHeight {
		t.Errorf("visibleRepoCount() = %d, should not exceed termHeight %d", count, m.termHeight)
	}
}

func TestEnsureCursorVisible(t *testing.T) {
	m := testRepoSelectModel()
	// Create enough repos to require scrolling
	m.repos = make([]RepoInfo, 100)
	for i := range m.repos {
		m.repos[i] = RepoInfo{Path: "/r", Name: "r"}
	}
	m.termHeight = 30

	// Move cursor past the visible window
	m.cursor = 50
	m.winOffset = 0
	m.ensureCursorVisible()

	wh := m.visibleRepoCount()
	if m.cursor < m.winOffset || m.cursor >= m.winOffset+wh {
		t.Errorf("cursor %d not visible in window [%d, %d)", m.cursor, m.winOffset, m.winOffset+wh)
	}
}

func TestHandleScanResult(t *testing.T) {
	m := testRepoSelectModel()
	m.loading = true
	m.selected = make(map[int]bool)

	repos := []RepoInfo{
		{Path: "/x", Name: "x", IsDirty: true},
		{Path: "/y", Name: "y", IsDirty: false},
		{Path: "/z", Name: "z", IsDirty: true},
	}
	result := m.handleScanResult(reposScanResultMsg{repos: repos})

	if result.loading {
		t.Error("loading should be false after scan result")
	}
	if len(result.repos) != 3 {
		t.Errorf("repos count = %d, want 3", len(result.repos))
	}
	// Dirty repos should be pre-selected
	if !result.selected[0] {
		t.Error("dirty repo 0 should be selected")
	}
	if result.selected[1] {
		t.Error("clean repo 1 should not be selected")
	}
	if !result.selected[2] {
		t.Error("dirty repo 2 should be selected")
	}
}

func TestHandleToggle(t *testing.T) {
	m := testRepoSelectModel()
	m.cursor = 1
	m.selected = map[int]bool{}

	// Toggle on
	m.handleToggle()
	if !m.selected[1] {
		t.Error("expected repo 1 to be selected after toggle on")
	}

	// Toggle off
	m.handleToggle()
	if m.selected[1] {
		t.Error("expected repo 1 to be deselected after toggle off")
	}
}

func TestHandleToggle_SingleSelect(t *testing.T) {
	m := testRepoSelectModel()
	m.singleSelect = true
	m.cursor = 1
	m.selected = map[int]bool{}

	m.handleToggle()
	if m.selected[1] {
		t.Error("toggle should be no-op in single select mode")
	}
}

func TestHandleSelectAll(t *testing.T) {
	m := testRepoSelectModel()
	m.selected = map[int]bool{}

	// Select all
	m.handleSelectAll()
	for i := range m.repos {
		if !m.selected[i] {
			t.Errorf("repo %d should be selected after select all", i)
		}
	}

	// Deselect all (all are selected, so toggles to none)
	m.handleSelectAll()
	for i := range m.repos {
		if m.selected[i] {
			t.Errorf("repo %d should be deselected after deselect all", i)
		}
	}
}

func TestHandleSelectAll_SingleSelect(t *testing.T) {
	m := testRepoSelectModel()
	m.singleSelect = true
	m.selected = map[int]bool{}

	m.handleSelectAll()
	for i := range m.repos {
		if m.selected[i] {
			t.Errorf("repo %d should not be selected in single select mode", i)
		}
	}
}

func TestHandleEnter(t *testing.T) {
	m := testRepoSelectModel()
	m.selected = map[int]bool{0: true, 2: true}
	m.caller = "commit"

	_, cmd := m.handleEnter()
	if cmd == nil {
		t.Fatal("expected cmd, got nil")
	}

	msg := cmd()
	done, ok := msg.(RepoSelectDoneMsg)
	if !ok {
		t.Fatalf("expected RepoSelectDoneMsg, got %T", msg)
	}
	if done.Caller != "commit" {
		t.Errorf("caller = %q, want %q", done.Caller, "commit")
	}
	if len(done.Paths) != 2 {
		t.Errorf("paths count = %d, want 2", len(done.Paths))
	}
}

func TestHandleEnter_SingleSelect(t *testing.T) {
	m := testRepoSelectModel()
	m.singleSelect = true
	m.cursor = 1

	_, cmd := m.handleEnter()
	if cmd == nil {
		t.Fatal("expected cmd, got nil")
	}

	msg := cmd()
	done, ok := msg.(RepoSelectDoneMsg)
	if !ok {
		t.Fatalf("expected RepoSelectDoneMsg, got %T", msg)
	}
	if len(done.Paths) != 1 || done.Paths[0] != "/b" {
		t.Errorf("paths = %v, want [/b]", done.Paths)
	}
}

func TestHandleKeyPress_Back(t *testing.T) {
	m := testRepoSelectModel()

	_, cmd := m.handleKeyPress(tea.KeyPressMsg{Code: tea.KeyEscape})
	if cmd == nil {
		t.Fatal("expected cmd for back, got nil")
	}
	msg := cmd()
	if _, ok := msg.(BackToMenuMsg); !ok {
		t.Fatalf("expected BackToMenuMsg, got %T", msg)
	}
}

func TestMeasureRepoChrome(t *testing.T) {
	st := SwissgitPalette.Styles()

	tests := []struct {
		name       string
		scrollable bool
	}{
		{"not scrollable", false},
		{"scrollable", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chrome := measureRepoChrome(st, tt.scrollable)
			if chrome <= 0 {
				t.Errorf("chrome = %d, want > 0", chrome)
			}
		})
	}

	// Scrollable should add more chrome than non-scrollable
	scrollChrome := measureRepoChrome(st, true)
	noScrollChrome := measureRepoChrome(st, false)
	if scrollChrome <= noScrollChrome {
		t.Errorf("scrollable chrome (%d) should be > non-scrollable chrome (%d)", scrollChrome, noScrollChrome)
	}
}

func TestRepoSelectView_Loading(t *testing.T) {
	m := testRepoSelectModel()
	m.loading = true

	v := m.View()
	if v == "" {
		t.Error("View() returned empty string during loading")
	}
}

func TestRepoSelectView_Empty(t *testing.T) {
	m := testRepoSelectModel()
	m.repos = nil
	m.loading = false

	v := m.View()
	if v == "" {
		t.Error("View() returned empty string for empty repos")
	}
}

func TestRepoSelectView_WithRepos(t *testing.T) {
	m := testRepoSelectModel()
	m.loading = false

	v := m.View()
	if v == "" {
		t.Error("View() returned empty string with repos")
	}
}

func TestRepoSelectUpdate_KeyUpDown(t *testing.T) {
	m := testRepoSelectModel()

	// Down
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1 after down", m.cursor)
	}

	// Up
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 after up", m.cursor)
	}

	// Up again wraps to last
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.cursor != 2 {
		t.Errorf("cursor = %d, want 2 after wrap up", m.cursor)
	}
}

func TestRepoSelectUpdate_Toggle(t *testing.T) {
	m := testRepoSelectModel()
	m.selected = map[int]bool{}

	// Toggle via Update (space key)
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if !m.selected[0] {
		t.Error("expected repo 0 to be selected after space toggle")
	}

	// Toggle off
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if m.selected[0] {
		t.Error("expected repo 0 to be deselected after second toggle")
	}
}

func TestRepoSelectUpdate_SelectAll(t *testing.T) {
	m := testRepoSelectModel()
	m.selected = map[int]bool{}

	// Ctrl+A to select all
	m, _ = m.Update(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})
	for i := range m.repos {
		if !m.selected[i] {
			t.Errorf("repo %d should be selected after ctrl+a", i)
		}
	}

	// Ctrl+A again to deselect all
	m, _ = m.Update(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})
	for i := range m.repos {
		if m.selected[i] {
			t.Errorf("repo %d should be deselected after second ctrl+a", i)
		}
	}
}

func TestRepoSelectUpdate_Enter(t *testing.T) {
	m := testRepoSelectModel()
	m.selected = map[int]bool{1: true}

	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected cmd from enter")
	}
	msg := cmd()
	done, ok := msg.(RepoSelectDoneMsg)
	if !ok {
		t.Fatalf("expected RepoSelectDoneMsg, got %T", msg)
	}
	if len(done.Paths) != 1 || done.Paths[0] != "/b" {
		t.Errorf("paths = %v, want [/b]", done.Paths)
	}
}

func TestRepoSelectUpdate_WindowSize(t *testing.T) {
	m := testRepoSelectModel()

	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 60})
	if m.termHeight != 60 {
		t.Errorf("termHeight = %d, want 60", m.termHeight)
	}
}

func TestRepoSelectUpdate_KeyPressWhileLoading(t *testing.T) {
	m := testRepoSelectModel()
	m.loading = true

	// Key press during loading should be a no-op
	m2, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if cmd != nil {
		t.Error("key press while loading should not produce a cmd")
	}
	if m2.cursor != m.cursor {
		t.Error("cursor should not change while loading")
	}
}

func TestMoveCursorUp(t *testing.T) {
	m := testRepoSelectModel()
	m.cursor = 0

	m.moveCursorUp()
	if m.cursor != 2 {
		t.Errorf("cursor = %d, want 2 after wrapping up from 0", m.cursor)
	}

	m.moveCursorUp()
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1", m.cursor)
	}
}

func TestMoveCursorDown(t *testing.T) {
	m := testRepoSelectModel()
	m.cursor = 2

	m.moveCursorDown()
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 after wrapping down from last", m.cursor)
	}
}

func TestRepoSelectMouseWheel(t *testing.T) {
	m := testRepoSelectModel()

	// Wheel down
	m, _ = m.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelDown}))
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1 after wheel down", m.cursor)
	}

	// Wheel up
	m, _ = m.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelUp}))
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 after wheel up", m.cursor)
	}
}

func TestRepoSelectMouseWheel_WhileLoading(t *testing.T) {
	m := testRepoSelectModel()
	m.loading = true

	m, _ = m.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelDown}))
	if m.cursor != 0 {
		t.Error("mouse wheel while loading should not change cursor")
	}
}

func TestRepoSelectMouseWheel_EmptyRepos(t *testing.T) {
	m := testRepoSelectModel()
	m.repos = nil

	m, _ = m.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelDown}))
	if m.cursor != 0 {
		t.Error("mouse wheel with empty repos should not change cursor")
	}
}

func TestRepoPrefix(t *testing.T) {
	st := SwissgitPalette.Styles()
	m := testRepoSelectModel()

	tests := []struct {
		name         string
		isCursor     bool
		isSelected   bool
		singleSelect bool
		wantEmpty    bool
	}{
		{"multi-select cursor+selected", true, true, false, false},
		{"multi-select not cursor, selected", false, true, false, false},
		{"multi-select not cursor, not selected", false, false, false, false},
		{"single-select cursor", true, false, true, false},
		{"single-select not cursor", false, false, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.singleSelect = tt.singleSelect
			result := m.repoPrefix(st, tt.isCursor, tt.isSelected)
			if tt.wantEmpty && result != "" {
				t.Errorf("expected empty prefix, got %q", result)
			}
			if !tt.wantEmpty && result == "" {
				t.Error("expected non-empty prefix")
			}
		})
	}
}

func TestRepoName(t *testing.T) {
	st := SwissgitPalette.Styles()
	m := testRepoSelectModel()

	// cursor item
	name := m.repoName(st, "test-repo", true, false)
	if name == "" {
		t.Error("repoName for cursor should not be empty")
	}

	// selected item
	name = m.repoName(st, "test-repo", false, true)
	if name == "" {
		t.Error("repoName for selected should not be empty")
	}

	// unselected item
	name = m.repoName(st, "test-repo", false, false)
	if name == "" {
		t.Error("repoName for unselected should not be empty")
	}
}

func TestRepoInfo(t *testing.T) {
	st := SwissgitPalette.Styles()
	m := testRepoSelectModel()

	// Dirty repo
	dirty := RepoInfo{Modified: 3}
	info := m.repoInfo(st, dirty, false, false)
	if info == "" {
		t.Error("dirty repo should have non-empty info")
	}

	// Clean repo, selected
	clean := RepoInfo{}
	info = m.repoInfo(st, clean, false, true)
	if info == "" {
		t.Error("clean selected repo should show checkmark")
	}

	// Clean repo, cursor
	info = m.repoInfo(st, clean, true, false)
	if info == "" {
		t.Error("clean cursor repo should show checkmark")
	}

	// Clean repo, not selected, not cursor
	info = m.repoInfo(st, clean, false, false)
	if info != "" {
		t.Errorf("clean unselected non-cursor repo should have empty info, got %q", info)
	}
}

func TestRepoBranch(t *testing.T) {
	st := SwissgitPalette.Styles()
	m := testRepoSelectModel()

	// Non-default branch
	r := RepoInfo{Branch: "feature/x", DefaultBranch: "main"}
	branch := m.repoBranch(st, r)
	if branch == "" {
		t.Error("non-default branch should render branch info")
	}

	// Default branch
	r2 := RepoInfo{Branch: "main", DefaultBranch: "main"}
	branch2 := m.repoBranch(st, r2)
	if branch2 != "" {
		t.Errorf("default branch should return empty, got %q", branch2)
	}

	// Empty branch
	r3 := RepoInfo{Branch: "", DefaultBranch: "main"}
	branch3 := m.repoBranch(st, r3)
	if branch3 != "" {
		t.Errorf("empty branch should return empty, got %q", branch3)
	}
}

func TestRepoSelectView_SingleSelect(t *testing.T) {
	m := testRepoSelectModel()
	m.singleSelect = true
	m.loading = false

	v := m.View()
	if v == "" {
		t.Error("View() returned empty for single select mode")
	}
}

func TestRepoSelectView_Scrollable(t *testing.T) {
	m := testRepoSelectModel()
	m.loading = false
	// Create many repos to trigger scrolling
	m.repos = make([]RepoInfo, 50)
	for i := range m.repos {
		m.repos[i] = RepoInfo{Path: "/r", Name: "r"}
	}
	m.termHeight = 20

	v := m.View()
	if v == "" {
		t.Error("View() returned empty for scrollable list")
	}

	// Move cursor to middle to trigger scroll-up indicator
	m.cursor = 25
	m.ensureCursorVisible()
	v2 := m.View()
	if v2 == "" {
		t.Error("View() returned empty after scrolling to middle")
	}
}

func TestNewRepoSelectModel(t *testing.T) {
	scanner := func(rootPath string) ([]RepoInfo, error) {
		return []RepoInfo{{Path: "/a", Name: "a"}}, nil
	}
	cfg := RepoSelectConfig{
		Palette:      SwissgitPalette,
		RootPath:     "/tmp/test",
		Caller:       "test-caller",
		ParentOffset: 5,
		TermHeight:   40,
		Scanner:      scanner,
		SingleSelect: true,
	}
	m := NewRepoSelectModel(cfg)

	if !m.loading {
		t.Error("new model should be loading")
	}
	if m.caller != "test-caller" {
		t.Errorf("caller = %q, want %q", m.caller, "test-caller")
	}
	if m.rootPath != "/tmp/test" {
		t.Errorf("rootPath = %q, want %q", m.rootPath, "/tmp/test")
	}
	if m.termHeight != 40 {
		t.Errorf("termHeight = %d, want 40", m.termHeight)
	}
	if m.parentOffset != 5 {
		t.Errorf("parentOffset = %d, want 5", m.parentOffset)
	}
	if !m.singleSelect {
		t.Error("expected singleSelect = true")
	}
	if m.selected == nil {
		t.Error("selected map should be initialized")
	}
}

func TestRepoSelectInit(t *testing.T) {
	scanner := func(rootPath string) ([]RepoInfo, error) {
		return nil, nil
	}
	cfg := RepoSelectConfig{
		Palette: SwissgitPalette,
		Scanner: scanner,
	}
	m := NewRepoSelectModel(cfg)
	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() should return a batched cmd")
	}
}

func TestScanRepos(t *testing.T) {
	called := false
	scanner := func(rootPath string) ([]RepoInfo, error) {
		called = true
		if rootPath != "/test/root" {
			t.Errorf("rootPath = %q, want %q", rootPath, "/test/root")
		}
		return []RepoInfo{{Path: "/a", Name: "a"}}, nil
	}
	cfg := RepoSelectConfig{
		Palette:  SwissgitPalette,
		RootPath: "/test/root",
		Scanner:  scanner,
	}
	m := NewRepoSelectModel(cfg)
	cmd := m.scanRepos()
	if cmd == nil {
		t.Fatal("scanRepos should return a cmd")
	}
	msg := cmd()
	result, ok := msg.(reposScanResultMsg)
	if !ok {
		t.Fatalf("expected reposScanResultMsg, got %T", msg)
	}
	if !called {
		t.Error("scanner function should have been called")
	}
	if len(result.repos) != 1 {
		t.Errorf("expected 1 repo, got %d", len(result.repos))
	}
}

func TestScanRepos_EmptyRoot(t *testing.T) {
	scanner := func(rootPath string) ([]RepoInfo, error) {
		if rootPath != "." {
			t.Errorf("expected default root '.', got %q", rootPath)
		}
		return nil, nil
	}
	cfg := RepoSelectConfig{
		Palette:  SwissgitPalette,
		RootPath: "", // empty root should default to "."
		Scanner:  scanner,
	}
	m := NewRepoSelectModel(cfg)
	cmd := m.scanRepos()
	msg := cmd()
	_, ok := msg.(reposScanResultMsg)
	if !ok {
		t.Fatalf("expected reposScanResultMsg, got %T", msg)
	}
}

func TestScanRepos_Error(t *testing.T) {
	scanner := func(rootPath string) ([]RepoInfo, error) {
		return nil, fmt.Errorf("scan failed")
	}
	cfg := RepoSelectConfig{
		Palette:  SwissgitPalette,
		RootPath: "/test",
		Scanner:  scanner,
	}
	m := NewRepoSelectModel(cfg)
	cmd := m.scanRepos()
	msg := cmd()
	result, ok := msg.(reposScanResultMsg)
	if !ok {
		t.Fatalf("expected reposScanResultMsg, got %T", msg)
	}
	if result.repos != nil {
		t.Error("expected nil repos on error")
	}
}

func TestVisibleRepoCount_SmallTerminal(t *testing.T) {
	m := testRepoSelectModel()
	m.termHeight = 5 // very small terminal
	m.parentOffset = 10
	count := m.visibleRepoCount()
	if count < 5 {
		t.Errorf("visibleRepoCount should be at least 5, got %d", count)
	}
}

func TestEnsureCursorVisible_CursorAboveWindow(t *testing.T) {
	m := testRepoSelectModel()
	m.repos = make([]RepoInfo, 100)
	for i := range m.repos {
		m.repos[i] = RepoInfo{Path: "/r", Name: "r"}
	}
	m.termHeight = 30
	// Set winOffset ahead of cursor so cursor < winOffset
	m.cursor = 5
	m.winOffset = 20
	m.ensureCursorVisible()
	if m.winOffset != 5 {
		t.Errorf("winOffset = %d, want 5 (should snap to cursor)", m.winOffset)
	}
}

func TestRepoSelectUpdate_SpinnerTickWhileLoading(t *testing.T) {
	scanner := func(rootPath string) ([]RepoInfo, error) {
		return nil, nil
	}
	cfg := RepoSelectConfig{
		Palette: SwissgitPalette,
		Scanner: scanner,
	}
	m := NewRepoSelectModel(cfg)
	// m.loading is true
	tickMsg := m.spinner.Tick()
	m2, cmd := m.Update(tickMsg)
	if cmd == nil {
		t.Error("spinner tick while loading should return a cmd")
	}
	_ = m2
}

func TestRepoSelectUpdate_SpinnerTickNotLoading(t *testing.T) {
	m := testRepoSelectModel()
	m.loading = false
	// Create a spinner tick - even though not loading, should be a no-op
	tickMsg := m.spinner.Tick()
	_, cmd := m.Update(tickMsg)
	// When not loading, spinner tick should not produce a cmd
	if cmd != nil {
		t.Error("spinner tick when not loading should return nil cmd")
	}
}

func TestRepoSelectUpdate_ScanResultMsg(t *testing.T) {
	m := testRepoSelectModel()
	m.loading = true
	repos := []RepoInfo{
		{Path: "/x", Name: "x", IsDirty: true},
		{Path: "/y", Name: "y"},
	}
	m2, cmd := m.Update(reposScanResultMsg{repos: repos})
	if m2.loading {
		t.Error("should not be loading after scan result")
	}
	if len(m2.repos) != 2 {
		t.Errorf("repos count = %d, want 2", len(m2.repos))
	}
	if cmd != nil {
		t.Error("scan result should return nil cmd")
	}
}

func TestHandleEnter_SingleSelect_OutOfBounds(t *testing.T) {
	m := testRepoSelectModel()
	m.singleSelect = true
	m.repos = nil // empty repos
	m.cursor = 5  // out of bounds

	_, cmd := m.handleEnter()
	if cmd != nil {
		t.Error("single select with cursor out of bounds should return nil cmd")
	}
}

func TestRepoSelectView_VisibleStartOverflow(t *testing.T) {
	m := testRepoSelectModel()
	m.loading = false
	m.repos = []RepoInfo{{Path: "/a", Name: "a"}}
	// Force winOffset past the end of repos
	m.winOffset = 100
	m.termHeight = 50

	v := m.View()
	if v == "" {
		t.Error("View() should handle visibleStart > len(repos)")
	}
}

func TestViewScrollDown_ScrollableNoRemaining(t *testing.T) {
	m := testRepoSelectModel()
	st := m.styles
	// scrollable=true, but visibleEnd == len(repos), so remaining == 0
	m.repos = make([]RepoInfo, 5)
	result := m.viewScrollDown(st, 5, true) // visibleEnd=5 == len(repos)
	if result != "\n" {
		t.Errorf("expected newline for scrollable with no remaining, got %q", result)
	}
}

func TestViewScrollDown_NotScrollableNoRemaining(t *testing.T) {
	m := testRepoSelectModel()
	st := m.styles
	m.repos = make([]RepoInfo, 3)
	result := m.viewScrollDown(st, 3, false)
	if result != "" {
		t.Errorf("expected empty for not scrollable with no remaining, got %q", result)
	}
}

func TestRepoSelectUpdate_UnhandledMsg(t *testing.T) {
	m := testRepoSelectModel()
	// Send a message type that's not handled by any case
	type customMsg struct{}
	m2, cmd := m.Update(customMsg{})
	if cmd != nil {
		t.Error("unhandled msg should return nil cmd")
	}
	_ = m2
}

func TestViewRepoLine(t *testing.T) {
	m := testRepoSelectModel()
	m.repos = []RepoInfo{
		{Path: "/a", Name: "repo-a", Branch: "feature/x", DefaultBranch: "main", Modified: 2},
		{Path: "/b", Name: "repo-b"},
	}
	m.selected = map[int]bool{0: true}
	m.cursor = 0
	st := m.styles

	// Cursor line
	line := m.viewRepoLine(st, 0)
	if line == "" {
		t.Error("viewRepoLine for cursor should not be empty")
	}

	// Non-cursor line
	line2 := m.viewRepoLine(st, 1)
	if line2 == "" {
		t.Error("viewRepoLine for non-cursor should not be empty")
	}
}
