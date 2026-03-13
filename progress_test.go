package curd

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func testTasks() []RepoTask {
	return []RepoTask{
		{Name: "repo-a", Path: "/a", Status: TaskPending},
		{Name: "repo-b", Path: "/b", Status: TaskRunning},
		{Name: "repo-c", Path: "/c", Status: TaskDone},
		{Name: "repo-d", Path: "/d", Status: TaskFailed},
	}
}

func TestCountFinished(t *testing.T) {
	tests := []struct {
		name         string
		tasks        []RepoTask
		wantFinished int
		wantAllDone  bool
	}{
		{
			name:         "mixed statuses",
			tasks:        testTasks(),
			wantFinished: 2, // done + failed
			wantAllDone:  false,
		},
		{
			name: "all done",
			tasks: []RepoTask{
				{Status: TaskDone},
				{Status: TaskFailed},
				{Status: TaskDone},
			},
			wantFinished: 3,
			wantAllDone:  true,
		},
		{
			name: "none finished",
			tasks: []RepoTask{
				{Status: TaskPending},
				{Status: TaskPending},
				{Status: TaskRunning},
			},
			wantFinished: 0,
			wantAllDone:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ProgressModel{Tasks: tt.tasks}
			finished, allDone := m.countFinished()
			if finished != tt.wantFinished {
				t.Errorf("finished = %d, want %d", finished, tt.wantFinished)
			}
			if allDone != tt.wantAllDone {
				t.Errorf("allDone = %v, want %v", allDone, tt.wantAllDone)
			}
		})
	}
}

func TestHandleTaskUpdate(t *testing.T) {
	tasks := []RepoTask{
		{Name: "a", Status: TaskPending},
		{Name: "b", Status: TaskPending},
		{Name: "c", Status: TaskPending},
	}
	m := NewProgressModel(tasks, SwissgitPalette)

	msg := RepoTaskUpdateMsg{Index: 0, Status: TaskDone, Result: "ok"}
	m, _ = m.handleTaskUpdate(msg)

	if m.Tasks[0].Status != TaskDone {
		t.Errorf("task 0 status = %d, want %d", m.Tasks[0].Status, TaskDone)
	}
	if m.Tasks[0].Result != "ok" {
		t.Errorf("task 0 result = %q, want %q", m.Tasks[0].Result, "ok")
	}
	if m.finished != 1 {
		t.Errorf("finished = %d, want 1", m.finished)
	}
	if m.done {
		t.Error("should not be done yet")
	}
}

func TestHandleTaskUpdate_AllDone(t *testing.T) {
	tasks := []RepoTask{
		{Name: "a", Status: TaskDone},
		{Name: "b", Status: TaskPending},
	}
	m := NewProgressModel(tasks, SwissgitPalette)

	msg := RepoTaskUpdateMsg{Index: 1, Status: TaskDone, Result: "ok"}
	m, cmd := m.handleTaskUpdate(msg)

	if !m.done {
		t.Error("expected done = true after all tasks complete")
	}
	if m.finished != 2 {
		t.Errorf("finished = %d, want 2", m.finished)
	}
	if cmd == nil {
		t.Fatal("expected a batched cmd, got nil")
	}

	// Execute the batched cmd to cover the AllTasksDoneMsg anonymous function
	batchMsg := cmd()
	batch, ok := batchMsg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("expected tea.BatchMsg, got %T", batchMsg)
	}
	// Execute each cmd in the batch to cover the anonymous function
	var foundAllDone bool
	for _, c := range batch {
		if c != nil {
			innerMsg := c()
			if _, ok := innerMsg.(AllTasksDoneMsg); ok {
				foundAllDone = true
			}
		}
	}
	if !foundAllDone {
		t.Error("expected AllTasksDoneMsg in batch")
	}
}

func TestIsDone(t *testing.T) {
	m := NewProgressModel([]RepoTask{{Status: TaskPending}}, SwissgitPalette)

	if m.IsDone() {
		t.Error("new model should not be done")
	}

	m.done = true
	if !m.IsDone() {
		t.Error("expected IsDone() = true after setting done")
	}
}

func TestNewProgressModel(t *testing.T) {
	tasks := testTasks()
	m := NewProgressModel(tasks, SwissgitPalette)

	if len(m.Tasks) != len(tasks) {
		t.Errorf("task count = %d, want %d", len(m.Tasks), len(tasks))
	}
	if m.done {
		t.Error("new model should not be done")
	}
	if m.finished != 0 {
		t.Errorf("finished = %d, want 0", m.finished)
	}
}

func TestProgressView(t *testing.T) {
	m := NewProgressModel(testTasks(), SwissgitPalette)
	v := m.View()
	if v == "" {
		t.Error("View() returned empty string")
	}
}

func TestProgressView_Done(t *testing.T) {
	tasks := []RepoTask{
		{Name: "a", Status: TaskDone},
		{Name: "b", Status: TaskDone},
	}
	m := NewProgressModel(tasks, SwissgitPalette)
	m.done = true
	m.finished = 2
	v := m.View()
	if v == "" {
		t.Error("View() returned empty string when done")
	}
}

func TestProgressView_ManyRunning(t *testing.T) {
	tasks := []RepoTask{
		{Name: "a", Status: TaskRunning},
		{Name: "b", Status: TaskRunning},
		{Name: "c", Status: TaskRunning},
		{Name: "d", Status: TaskRunning},
		{Name: "e", Status: TaskRunning},
	}
	m := NewProgressModel(tasks, SwissgitPalette)
	v := m.View()
	if v == "" {
		t.Error("View() returned empty string with many running tasks")
	}
}

func TestProgressInit(t *testing.T) {
	m := NewProgressModel(testTasks(), SwissgitPalette)
	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() should return spinner tick cmd")
	}
}

func TestProgressUpdate_WindowSize(t *testing.T) {
	m := NewProgressModel(testTasks(), SwissgitPalette)

	tests := []struct {
		name  string
		width int
	}{
		{"wide terminal", 200},
		{"narrow terminal", 30},
		{"very narrow terminal", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, cmd := m.Update(tea.WindowSizeMsg{Width: tt.width, Height: 40})
			// Should not produce a command from resize
			if cmd != nil {
				t.Error("WindowSizeMsg should not produce a cmd")
			}
			_ = updated
		})
	}
}

func TestProgressHandleResize(t *testing.T) {
	m := NewProgressModel(testTasks(), SwissgitPalette)

	// Wide: capped at 60
	m.handleResize(tea.WindowSizeMsg{Width: 200, Height: 40})

	// Narrow: clamped to 20
	m.handleResize(tea.WindowSizeMsg{Width: 25, Height: 40})

	// Normal range
	m.handleResize(tea.WindowSizeMsg{Width: 60, Height: 40})
}

func TestProgressUpdate_TaskUpdateMsg(t *testing.T) {
	tasks := []RepoTask{
		{Name: "a", Status: TaskPending},
		{Name: "b", Status: TaskPending},
	}
	m := NewProgressModel(tasks, SwissgitPalette)

	msg := RepoTaskUpdateMsg{Index: 0, Status: TaskDone, Result: "ok"}
	m, cmd := m.Update(msg)
	if cmd == nil {
		t.Error("expected cmd from task update (progress bar percent)")
	}
	if m.Tasks[0].Status != TaskDone {
		t.Errorf("task 0 status = %d, want %d", m.Tasks[0].Status, TaskDone)
	}
	if m.finished != 1 {
		t.Errorf("finished = %d, want 1", m.finished)
	}
}

func TestProgressUpdate_TaskUpdateMsg_Error(t *testing.T) {
	tasks := []RepoTask{
		{Name: "a", Status: TaskPending},
	}
	m := NewProgressModel(tasks, SwissgitPalette)

	msg := RepoTaskUpdateMsg{Index: 0, Status: TaskFailed, Error: "something broke"}
	m, _ = m.Update(msg)
	if m.Tasks[0].Status != TaskFailed {
		t.Errorf("task 0 status = %d, want %d", m.Tasks[0].Status, TaskFailed)
	}
	if m.Tasks[0].Error != "something broke" {
		t.Errorf("task 0 error = %q, want %q", m.Tasks[0].Error, "something broke")
	}
}

func TestProgressUpdate_TaskUpdateMsg_OutOfBounds(t *testing.T) {
	tasks := []RepoTask{
		{Name: "a", Status: TaskPending},
	}
	m := NewProgressModel(tasks, SwissgitPalette)

	// Out of bounds index should not panic
	msg := RepoTaskUpdateMsg{Index: 99, Status: TaskDone}
	m, _ = m.Update(msg)
	// Task should remain unchanged
	if m.Tasks[0].Status != TaskPending {
		t.Errorf("task 0 should remain pending, got %d", m.Tasks[0].Status)
	}
}

func TestProgressUpdate_UnhandledMsg(t *testing.T) {
	m := NewProgressModel(testTasks(), SwissgitPalette)
	// Send an unrelated key press
	m2, cmd := m.Update(tea.KeyPressMsg{Code: 'x'})
	if cmd != nil {
		t.Error("unhandled msg should not produce a cmd")
	}
	_ = m2
}

func TestProgressUpdate_SpinnerTick_NotDone(t *testing.T) {
	m := NewProgressModel(testTasks(), SwissgitPalette)
	// Get a real spinner tick message
	tickMsg := m.spinner.Tick()
	m2, cmd := m.Update(tickMsg)
	// Should forward to spinner and return a cmd (next tick)
	if cmd == nil {
		t.Error("expected cmd from spinner tick when not done")
	}
	_ = m2
}

func TestProgressUpdate_SpinnerTick_Done(t *testing.T) {
	m := NewProgressModel(testTasks(), SwissgitPalette)
	m.done = true
	tickMsg := m.spinner.Tick()
	_, cmd := m.Update(tickMsg)
	// When done, spinner tick should return nil cmd
	if cmd != nil {
		t.Error("expected nil cmd from spinner tick when done")
	}
}

func TestProgressUpdate_FrameMsg(t *testing.T) {
	m := NewProgressModel(testTasks(), SwissgitPalette)
	// Get a real FrameMsg by calling SetPercent on the progress bar
	frameCmd := m.bar.SetPercent(0.5)
	if frameCmd == nil {
		t.Fatal("SetPercent should return a cmd")
	}
	frameMsg := frameCmd()
	m2, cmd := m.Update(frameMsg)
	// FrameMsg should be forwarded to the progress bar
	_ = m2
	_ = cmd // may or may not produce a follow-up cmd
}

func TestHandleTaskUpdate_OutOfBoundsNotAllDone(t *testing.T) {
	// Test OOB index with remaining pending tasks (not all done)
	tasks := []RepoTask{
		{Name: "a", Status: TaskPending},
		{Name: "b", Status: TaskPending},
	}
	m := NewProgressModel(tasks, SwissgitPalette)

	msg := RepoTaskUpdateMsg{Index: 99, Status: TaskDone}
	m, cmd := m.handleTaskUpdate(msg)
	// Tasks should remain unchanged
	if m.Tasks[0].Status != TaskPending {
		t.Errorf("task 0 should remain pending, got %d", m.Tasks[0].Status)
	}
	if m.done {
		t.Error("should not be done with pending tasks")
	}
	// Should still return a progress cmd (SetPercent(0))
	_ = cmd
}

func TestViewFailCount(t *testing.T) {
	tests := []struct {
		name      string
		tasks     []RepoTask
		wantEmpty bool
	}{
		{
			name: "no failures",
			tasks: []RepoTask{
				{Status: TaskDone},
				{Status: TaskPending},
			},
			wantEmpty: true,
		},
		{
			name: "has failures",
			tasks: []RepoTask{
				{Status: TaskFailed},
				{Status: TaskDone},
				{Status: TaskFailed},
			},
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ProgressModel{
				Tasks:  tt.tasks,
				styles: SwissgitPalette.Styles(),
			}
			result := m.viewFailCount(m.styles)
			if tt.wantEmpty && result != "" {
				t.Errorf("expected empty string, got %q", result)
			}
			if !tt.wantEmpty && result == "" {
				t.Error("expected non-empty fail count string")
			}
		})
	}
}
