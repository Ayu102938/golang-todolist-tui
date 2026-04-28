package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func initialModel() model {
	todos, _ := loadTodos()
	ti := textinput.New()
	ti.Placeholder = "タスク名..."
	return model{todos: todos, input: ti, mode: viewMode}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode == addMode || m.mode == editMode {
			switch msg.String() {
			case "enter":
				val := m.input.Value()
				if val != "" {
					if m.mode == addMode {
						m.todos = append(m.todos, Todo{Title: val, DueDate: time.Now().AddDate(0, 0, 1)})
					} else {
						m.todos[m.cursor].Title = val
					}
					saveTodos(m.todos)
					m.input.SetValue(""); m.mode = viewMode
				}
			case "esc": m.input.SetValue(""); m.mode = viewMode
			default: m.input, cmd = m.input.Update(msg); return m, cmd
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q": return m, tea.Quit
		case "up", "k": if m.cursor > 0 { m.cursor-- }
		case "down", "j": if m.cursor < len(m.todos)-1 { m.cursor++ }
		case "enter":
			if len(m.todos) > 0 {
				m.todos[m.cursor].Completed = !m.todos[m.cursor].Completed
				saveTodos(m.todos)
			}
		case "a": m.mode = addMode; m.input.Focus(); return m, textinput.Blink
		case "e":
			if len(m.todos) > 0 { m.mode = editMode; m.input.SetValue(m.todos[m.cursor].Title); m.input.Focus() }
		case "p":
			if len(m.todos) > 0 { m.todos[m.cursor].Priority = (m.todos[m.cursor].Priority + 1) % 3; saveTodos(m.todos) }
		case "f": m.filterDone = !m.filterDone
		case "d":
			if len(m.todos) > 0 {
				m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)
				if m.cursor >= len(m.todos) && m.cursor > 0 { m.cursor-- }
				saveTodos(m.todos)
			}
		}
	}
	return m, nil
}

var (
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1).Bold(true)
	progStyle = lipgloss.NewStyle().Margin(1, 0)
)

func (m model) View() string {
	if m.mode == addMode || m.mode == editMode {
		return "\n  タスクを入力:\n" + m.input.View() + "\n\n  enter: 確定 • esc: キャンセル"
	}

	doneCount := 0
	for _, t := range m.todos { if t.Completed { doneCount++ } }
	percent := 0.0
	if len(m.todos) > 0 { percent = float64(doneCount) / float64(len(m.todos)) }
	pBar := progress.New(progress.WithDefaultGradient())
	pBar.Width = 40
	
	s := titleStyle.Render("TODO リスト") + "\n"
	s += progStyle.Render(pBar.ViewAs(percent)) + fmt.Sprintf(" %.0f%%\n\n", percent*100)
	
	pStr := []string{"Low", "Mid", "High"}
	for i, todo := range m.todos {
		if m.filterDone && todo.Completed { continue }
		cursor := " "
		if m.cursor == i { cursor = ">" }
		checked := " "
		if todo.Completed { checked = "x" }
		s += fmt.Sprintf("%s [%s] [%-4s] %-10s %s\n", cursor, checked, pStr[todo.Priority], todo.DueDate.Format("01/02"), todo.Title)
	}
	s += "\n j/k:移動 • e:編集 • p:優先度 • f:フィルター • a:追加 • d:削除 • q:終了"
	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v", err); os.Exit(1)
	}
}
