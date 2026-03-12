package curd

import (
	"fmt"
	"strings"
)

// ResultModel displays a summary of completed operations.
type ResultModel struct {
	Tasks  []RepoTask
	Title  string
	Width  int
	styles StyleSet
}

func NewResultModel(title string, tasks []RepoTask, palette Palette) ResultModel {
	return ResultModel{
		Tasks:  tasks,
		Title:  title,
		Width:  80,
		styles: palette.Styles(),
	}
}

// TruncateError trims long error messages, keeping the juicy tail.
func TruncateError(err string, maxLen int) string {
	if maxLen < 40 {
		maxLen = 40
	}
	if len(err) <= maxLen {
		return err
	}
	if idx := strings.LastIndex(err, "error: "); idx >= 0 {
		tail := err[idx+len("error: "):]
		if len(tail) <= maxLen {
			return tail
		}
	}
	return "…" + err[len(err)-maxLen+1:]
}

func (m ResultModel) View() string {
	st := m.styles
	var succeeded, failed int
	var okTasks, failTasks []RepoTask
	for _, t := range m.Tasks {
		switch t.Status {
		case TaskDone:
			succeeded++
			okTasks = append(okTasks, t)
		case TaskFailed:
			failed++
			failTasks = append(failTasks, t)
		}
	}

	var s string

	// Summary banner
	banner := st.ResultAccent.Render(m.Title) + st.ResultDim.Render("  ")
	banner += st.ResultOk.Render(fmt.Sprintf("✔ %d", succeeded))
	if failed > 0 {
		banner += st.ResultDim.Render("  ") + st.ResultFail.Render(fmt.Sprintf("✗ %d", failed))
	}
	s += st.SummaryBox.Render(banner) + "\n\n"

	maxErr := m.Width - 10
	if maxErr < 40 {
		maxErr = 40
	}

	// Failures first
	if len(failTasks) > 0 {
		var failContent string
		for _, t := range failTasks {
			errMsg := TruncateError(t.Error, maxErr)
			failContent += fmt.Sprintf("  %s %s\n", st.ResultFail.Render("✗"), st.NameStyle.Render(t.Name))
			failContent += fmt.Sprintf("    %s\n", st.ResultDim.Render(errMsg))
		}
		failContent = strings.TrimRight(failContent, "\n")
		s += st.FailBox.Render(failContent) + "\n\n"
	}

	// Successes
	if len(okTasks) > 0 {
		var okContent string
		for _, t := range okTasks {
			line := fmt.Sprintf("  %s %s", st.ResultOk.Render("✔"), st.NameStyle.Render(t.Name))
			if t.Result != "" {
				line += "  " + st.ResultDim.Render(t.Result)
			}
			okContent += line + "\n"
		}
		okContent = strings.TrimRight(okContent, "\n")
		s += st.SuccessBox.Render(okContent) + "\n"
	}

	return s
}
