package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func initialModel() model {
	todos, err := loadTodos()
	if err != nil {
		log.Printf("warning: failed to load todos: %v", err)
		todos = []Todo{}
	}
	categories := []string{"Home"}
	catMap := map[string]bool{"Home": true}
	for i := range todos {
		if todos[i].Category == "" {
			todos[i].Category = "Home"
		}
		if !catMap[todos[i].Category] {
			catMap[todos[i].Category] = true
			categories = append(categories, todos[i].Category)
		}
	}
	ti := textinput.New()
	ti.Placeholder = "タスク名..."
	ti.CharLimit = maxInputLen
	return model{todos: todos, categories: categories, activeTab: 0, input: ti, mode: viewMode}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.mode == addMode || m.mode == editMode || m.mode == categoryAddMode {
			switch msg.String() {
			case "enter":
				val := m.input.Value()
				if val != "" {
				switch m.mode {
					case addMode:
						m.todos = append(m.todos, Todo{Title: val, DueDate: time.Now().AddDate(0, 0, 1), Category: m.categories[m.activeTab]})
						if err := saveTodos(m.todos); err != nil {
							log.Printf("error: failed to save todos: %v", err)
						}
					case editMode:
						idx := m.getFilteredIndex(m.cursor)
						if idx >= 0 {
							m.todos[idx].Title = val
							if err := saveTodos(m.todos); err != nil {
								log.Printf("error: failed to save todos: %v", err)
							}
						}
					case categoryAddMode:
						m.categories = append(m.categories, val)
						m.activeTab = len(m.categories) - 1
					}
					m.input.SetValue(""); m.mode = viewMode
				}
			case "esc": m.input.SetValue(""); m.mode = viewMode
			default: m.input, cmd = m.input.Update(msg); return m, cmd
			}
			return m, nil
		}
		if m.mode == categoryDeleteMode {
			switch msg.String() {
			case "enter":
				if m.activeTab > 0 {
					removed := m.categories[m.activeTab]
					m.categories = append(m.categories[:m.activeTab], m.categories[m.activeTab+1:]...)
					var remaining []Todo
					for _, t := range m.todos {
						if t.Category != removed {
							remaining = append(remaining, t)
						}
					}
					m.todos = remaining
					m.activeTab = 0
					m.cursor = 0
					if err := saveTodos(m.todos); err != nil {
						log.Printf("error: failed to save todos: %v", err)
					}
				}
				m.mode = viewMode
			case "esc":
				m.mode = viewMode
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q": return m, tea.Quit
		case "h": if m.activeTab > 0 { m.activeTab--; m.cursor = 0 }
		case "l": if m.activeTab < len(m.categories)-1 { m.activeTab++; m.cursor = 0 }
		case "up", "k": if m.cursor > 0 { m.cursor-- }
		case "down", "j": if m.cursor < m.filteredCount()-1 { m.cursor++ }
		case "enter":
			if m.filteredCount() > 0 {
				idx := m.getFilteredIndex(m.cursor)
				if idx >= 0 {
					m.todos[idx].Completed = !m.todos[idx].Completed
					if err := saveTodos(m.todos); err != nil {
						log.Printf("error: failed to save todos: %v", err)
					}
				}
			}
		case "a": m.mode = addMode; m.input.Placeholder = "タスク名..."; m.input.SetValue(""); m.input.Focus(); return m, textinput.Blink
		case "n": m.mode = categoryAddMode; m.input.Placeholder = "カテゴリ名..."; m.input.SetValue(""); m.input.Focus(); return m, textinput.Blink
		case "e":
			if m.filteredCount() > 0 {
				idx := m.getFilteredIndex(m.cursor)
				if idx >= 0 {
					m.mode = editMode
					m.input.SetValue(m.todos[idx].Title)
					m.input.Focus()
				}
			}
		case "p":
			if m.filteredCount() > 0 {
				idx := m.getFilteredIndex(m.cursor)
				if idx >= 0 {
					m.todos[idx].Priority = (m.todos[idx].Priority + 1) % 3
					if err := saveTodos(m.todos); err != nil {
						log.Printf("error: failed to save todos: %v", err)
					}
				}
			}
		case "f": m.filterDone = !m.filterDone
		case "s":
			m.sortField = (m.sortField + 1) % 4
			m.sortTodos()
			m.cursor = 0
		case "x": m.mode = categoryDeleteMode
		case "d":
			if m.filteredCount() > 0 {
				idx := m.getFilteredIndex(m.cursor)
				if idx >= 0 {
					m.todos = append(m.todos[:idx], m.todos[idx+1:]...)
					if m.cursor >= m.filteredCount() && m.cursor > 0 {
						m.cursor--
					}
					if err := saveTodos(m.todos); err != nil {
						log.Printf("error: failed to save todos: %v", err)
					}
				}
			}
		}
	}
	return m, nil
}

func (m model) filteredCount() int {
	c := 0
	for _, t := range m.todos {
		if t.Category == m.categories[m.activeTab] && (!m.filterDone || !t.Completed) { c++ }
	}
	return c
}

func (m model) getFilteredIndex(target int) int {
	c := 0
	for i, t := range m.todos {
		if t.Category == m.categories[m.activeTab] && (!m.filterDone || !t.Completed) {
			if c == target { return i }
			c++
		}
	}
	return -1
}

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
			return "\n  Home カテゴリは削除できません\n\n  enter: 戻る"
		}
		return "\n  \"" + cat + "\" カテゴリとそのタスクを削除しますか？\n\n  enter: 確定 • esc: キャンセル"
	}
	if m.mode == addMode || m.mode == editMode || m.mode == categoryAddMode {
		return "\n  " + m.input.Placeholder + "\n" + m.input.View() + "\n\n  enter: 確定 • esc: キャンセル"
	}
	var tabViews []string
	for i, cat := range m.categories {
		if i == m.activeTab { tabViews = append(tabViews, activeTabStyle.Render(cat)) } else { tabViews = append(tabViews, tabStyle.Render(cat)) }
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
			if m.cursor == count { cursor = ">" }
			checked := " "
			if todo.Completed { checked = "x" }
			dueStr := todo.DueDate.Format("01/02")
			title := todo.Title
			if len(title) > contentWidth-25 {
				title = title[:contentWidth-25] + "..."
			}
			line := fmt.Sprintf("%s [%s] [%-4s] %-10s %s", cursor, checked, pStr[todo.Priority], dueStr, title)
			if todo.Completed {
				line = completedStyle.Render(line)
			} else if todo.DueDate.Before(now) {
				line = overdueStyle.Render(line)
			}
			s += line + "\n"
			count++
		}
	}
	s += "\n" + lipgloss.NewStyle().Width(contentWidth).Render(" h/l:タブ • n:カテゴリ追加 • x:カテゴリ削除 • a:追加 • j/k:移動 • e:編集 • p:優先度 • s:ソート • f:フィルター • d:削除 • q:終了")
	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v", err); os.Exit(1)
	}
}
