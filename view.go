package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func maxContentWidth(m model) int {
	if m.width <= 0 {
		return 80
	}
	return m.width - 4
}

func renderProgressBar(m model, total, done, width int) string {
	if total == 0 {
		return m.ui.progressBg.Render(strings.Repeat("░", width))
	}
	filled := int(float64(done) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}
	return m.ui.progressFg.Render(strings.Repeat("█", filled)) + m.ui.progressBg.Render(strings.Repeat("░", width-filled))
}

func sortLabel(sf sortField) string {
	switch sf {
	case sortPriority:
		return msg.SortPriority
	case sortDueDate:
		return msg.SortDueDate
	case sortName:
		return msg.SortName
	default:
		return msg.SortNone
	}
}

func (m model) View() string {
	if m.mode == categoryDeleteMode {
		cat := m.categories[m.activeTab]
		if m.activeTab == 0 {
			return "\n  " + m.config.DefaultCategory + " " + msg.CannotDeleteHome + "\n\n  enter: " + msg.Return
		}
		return "\n  \"" + cat + "\" " + msg.DeleteConfirm + "\n\n  enter: " + msg.Confirm + " • esc: " + msg.Cancel
	}
	if m.mode == addMode || m.mode == addDateMode || m.mode == editMode || m.mode == categoryAddMode || m.mode == searchMode || m.mode == descMode {
		return "\n  " + m.input.Placeholder + "\n" + m.input.View() + "\n\n  enter: " + msg.Confirm + " • esc: " + msg.Cancel
	}
	var tabViews []string
	for i, cat := range m.categories {
		if i == m.activeTab {
			tabViews = append(tabViews, m.ui.activeTab.Render(cat))
		} else {
			tabViews = append(tabViews, m.ui.tab.Render(cat))
		}
	}
	tabs := lipgloss.JoinHorizontal(lipgloss.Top, tabViews...)
	s := m.ui.title.Render(msg.Title) + "\n" + tabs + "\n"

	activeCat := m.categories[m.activeTab]
	var catTotal, catDone int
	for _, t := range m.todos {
		if t.Category == activeCat {
			catTotal++
			if t.Completed {
				catDone++
			}
		}
	}
	barWidth := maxContentWidth(m) - 12
	if barWidth < 5 {
		barWidth = 5
	}
	bar := renderProgressBar(m, catTotal, catDone, barWidth)
	s += fmt.Sprintf(" %d/%d %s\n\n", catDone, catTotal, bar)

	contentWidth := maxContentWidth(m)
	pStr := []string{"Low", "Mid", "High"}
	now := time.Now()
	availableHeight := m.height - 10
	if availableHeight < 1 {
		availableHeight = 1
	}
	count := 0
	for _, todo := range m.todos {
		if count >= availableHeight {
			s += "  ...\n"
			break
		}
		if todo.Category == m.categories[m.activeTab] && (!m.filterDone || !todo.Completed) {
			cursor := " "
			if m.cursor == count {
				cursor = ">"
			}
			checked := " "
			if todo.Completed {
				checked = "x"
			}
			dueStr := todo.DueDate.Format("01/02")
			title := todo.Title
			if len(title) > contentWidth-28 {
				title = title[:contentWidth-28] + "..."
			}
			descMarker := ""
			if todo.Description != "" {
				descMarker = "…"
			}
			line := fmt.Sprintf("%s [%s] [%-4s] %-10s %s%s", cursor, checked, pStr[todo.Priority], dueStr, title, descMarker)
			if todo.Completed {
				line = m.ui.completed.Render(line)
			} else if todo.DueDate.Before(now) {
				line = m.ui.overdue.Render(line)
			}
			s += line + "\n"
			count++
		}
	}
	if count == 0 {
		if catTotal == 0 {
			s += "  " + msg.NoTasks + "\n"
		} else {
			s += "  " + msg.NoVisibleTasks + "\n"
		}
	}
	s += "\n" + lipgloss.NewStyle().Width(contentWidth).Render(msg.Help)
	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
