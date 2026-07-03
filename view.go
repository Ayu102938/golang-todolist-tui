package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle      = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1).Bold(true)
	tabStyle        = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("240"))
	activeTabStyle  = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("255")).Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false)
	overdueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	completedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	progressBgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	progressFgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
)

func maxContentWidth(m model) int {
	if m.width <= 0 {
		return 80
	}
	return m.width - 4
}

func renderProgressBar(total, done int, width int) string {
	if total == 0 {
		return progressBgStyle.Render(strings.Repeat("░", width))
	}
	filled := int(float64(done) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}
	return progressFgStyle.Render(strings.Repeat("█", filled)) + progressBgStyle.Render(strings.Repeat("░", width-filled))
}

func sortLabel(sf sortField) string {
	switch sf {
	case sortPriority:
		return "優先度"
	case sortDueDate:
		return "期限"
	case sortName:
		return "名前"
	default:
		return "なし"
	}
}

func (m model) View() string {
	if m.mode == categoryDeleteMode {
		cat := m.categories[m.activeTab]
		if m.activeTab == 0 {
			return "\n  " + defaultCategory + " カテゴリは削除できません\n\n  enter: 戻る"
		}
		return "\n  \"" + cat + "\" カテゴリとそのタスクを削除しますか？\n\n  enter: 確定 • esc: キャンセル"
	}
	if m.mode == addMode || m.mode == addDateMode || m.mode == editMode || m.mode == categoryAddMode || m.mode == searchMode || m.mode == descMode {
		return "\n  " + m.input.Placeholder + "\n" + m.input.View() + "\n\n  enter: 確定 • esc: キャンセル"
	}
	var tabViews []string
	for i, cat := range m.categories {
		if i == m.activeTab {
			tabViews = append(tabViews, activeTabStyle.Render(cat))
		} else {
			tabViews = append(tabViews, tabStyle.Render(cat))
		}
	}
	tabs := lipgloss.JoinHorizontal(lipgloss.Top, tabViews...)
	s := titleStyle.Render("TODO リスト") + "\n" + tabs + "\n"

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
	bar := renderProgressBar(catTotal, catDone, barWidth)
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
				line = completedStyle.Render(line)
			} else if todo.DueDate.Before(now) {
				line = overdueStyle.Render(line)
			}
			s += line + "\n"
			count++
		}
	}
	if count == 0 {
		if catTotal == 0 {
			s += "  タスクがありません\n"
		} else {
			s += "  表示できるタスクがありません\n"
		}
	}
	s += "\n" + lipgloss.NewStyle().Width(contentWidth).Render(" h/l:タブ • n:カテゴリ追加 • x:カテゴリ削除 • a:追加 • j/k:移動 • e:編集 • i:詳細 • p:優先度 • s:ソート • f:フィルター • /:検索 • d:削除 • u:元に戻す • q:終了")
	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
