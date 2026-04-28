package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func initialModel() model {
	todos, _ := loadTodos()
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
	return model{todos: todos, categories: categories, activeTab: 0, input: ti, mode: viewMode}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode == addMode || m.mode == editMode || m.mode == categoryAddMode {
			switch msg.String() {
			case "enter":
				val := m.input.Value()
				if val != "" {
					switch m.mode {
					case addMode:
						m.todos = append(m.todos, Todo{Title: val, DueDate: time.Now().AddDate(0, 0, 1), Category: m.categories[m.activeTab]})
						saveTodos(m.todos)
					case editMode:
						m.todos[m.getFilteredIndex(m.cursor)].Title = val
						saveTodos(m.todos)
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
		switch msg.String() {
		case "ctrl+c", "q": return m, tea.Quit
		case "h": if m.activeTab > 0 { m.activeTab--; m.cursor = 0 }
		case "l": if m.activeTab < len(m.categories)-1 { m.activeTab++; m.cursor = 0 }
		case "up", "k": if m.cursor > 0 { m.cursor-- }
		case "down", "j": if m.cursor < m.filteredCount()-1 { m.cursor++ }
		case "enter":
			if m.filteredCount() > 0 {
				m.todos[m.getFilteredIndex(m.cursor)].Completed = !m.todos[m.getFilteredIndex(m.cursor)].Completed
				saveTodos(m.todos)
			}
		case "a": m.mode = addMode; m.input.Placeholder = "タスク名..."; m.input.SetValue(""); m.input.Focus(); return m, textinput.Blink
		case "n": m.mode = categoryAddMode; m.input.Placeholder = "カテゴリ名..."; m.input.SetValue(""); m.input.Focus(); return m, textinput.Blink
		case "e":
			if m.filteredCount() > 0 { m.mode = editMode; m.input.SetValue(m.todos[m.getFilteredIndex(m.cursor)].Title); m.input.Focus() }
		case "p":
			if m.filteredCount() > 0 { m.todos[m.getFilteredIndex(m.cursor)].Priority = (m.todos[m.getFilteredIndex(m.cursor)].Priority + 1) % 3; saveTodos(m.todos) }
		case "f": m.filterDone = !m.filterDone
		case "d":
			if m.filteredCount() > 0 {
				idx := m.getFilteredIndex(m.cursor)
				m.todos = append(m.todos[:idx], m.todos[idx+1:]...)
				if m.cursor >= m.filteredCount() && m.cursor > 0 { m.cursor-- }
				saveTodos(m.todos)
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
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1).Bold(true)
	tabStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("240"))
	activeTabStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("255")).Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false)
)

func (m model) View() string {
	if m.mode == addMode || m.mode == editMode || m.mode == categoryAddMode {
		return "\n  " + m.input.Placeholder + "\n" + m.input.View() + "\n\n  enter: 確定 • esc: キャンセル"
	}
	tabs := ""
	for i, cat := range m.categories {
		if i == m.activeTab { tabs += activeTabStyle.Render(cat) } else { tabs += tabStyle.Render(cat) }
	}
	s := titleStyle.Render("TODO リスト") + "\n" + tabs + "\n\n"
	
	pStr := []string{"Low", "Mid", "High"}
	count := 0
	for _, todo := range m.todos {
		if todo.Category == m.categories[m.activeTab] && (!m.filterDone || !todo.Completed) {
			cursor := " "
			if m.cursor == count { cursor = ">" }
			checked := " "
			if todo.Completed { checked = "x" }
			s += fmt.Sprintf("%s [%s] [%-4s] %-10s %s\n", cursor, checked, pStr[todo.Priority], todo.DueDate.Format("01/02"), todo.Title)
			count++
		}
	}
	s += "\n h/l:タブ移動 • n:カテゴリ追加 • a:タスク追加 • j/k:移動 • e:編集 • p:優先度 • f:フィルター • d:削除 • q:終了"
	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v", err); os.Exit(1)
	}
}
