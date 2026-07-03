package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func initialModel(storage Storage) model {
	todos, err := storage.Load()
	if err != nil {
		log.Printf("warning: failed to load todos: %v", err)
		todos = []Todo{}
	}
	categories := []string{defaultCategory}
	catMap := map[string]bool{defaultCategory: true}
	for i := range todos {
		if todos[i].Category == "" {
			todos[i].Category = defaultCategory
		}
		if !catMap[todos[i].Category] {
			catMap[todos[i].Category] = true
			categories = append(categories, todos[i].Category)
		}
	}
	ti := textinput.New()
	ti.Placeholder = "タスク名..."
	ti.CharLimit = maxInputLen
	return model{todos: todos, categories: categories, activeTab: 0, input: ti, mode: viewMode, storage: storage}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.mode == addMode || m.mode == addDateMode || m.mode == editMode || m.mode == categoryAddMode || m.mode == searchMode || m.mode == descMode {
			return m.handleInputKey(msg)
		}
		if m.mode == categoryDeleteMode {
			return m.handleCategoryDeleteKey(msg)
		}
		return m.handleViewKey(msg)
	}
	return m, nil
}

func (m model) handleInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		val := m.input.Value()
		switch m.mode {
		case addMode:
			if val != "" {
				m.pendingTitle = val
				m.mode = addDateMode
				m.input.Placeholder = "期限 YYYY-MM-DD (enter=明日)..."
				m.input.SetValue("")
			}
		case addDateMode:
			m.undoSnapshot()
			if val != "" {
				if dueDate, err := time.Parse("2006-01-02", val); err == nil {
					m.todos = append(m.todos, Todo{Title: m.pendingTitle, DueDate: dueDate, Category: m.categories[m.activeTab]})
				} else {
					m.todos = append(m.todos, Todo{Title: m.pendingTitle, DueDate: time.Now().AddDate(0, 0, 1), Category: m.categories[m.activeTab]})
				}
			} else {
				m.todos = append(m.todos, Todo{Title: m.pendingTitle, DueDate: time.Now().AddDate(0, 0, 1), Category: m.categories[m.activeTab]})
			}
			m.input.SetValue("")
			m.mode = viewMode
		case editMode:
			idx := m.getFilteredIndex(m.cursor)
			if idx >= 0 && val != "" {
				m.undoSnapshot()
				m.todos[idx].Title = val
			}
			m.input.SetValue("")
			m.mode = viewMode
		case categoryAddMode:
			if val != "" {
				m.categories = append(m.categories, val)
				m.activeTab = len(m.categories) - 1
			}
			m.input.SetValue("")
			m.mode = viewMode
		case searchMode:
			m.searchQuery = val
			m.cursor = 0
			m.mode = viewMode
		case descMode:
			idx := m.getFilteredIndex(m.cursor)
			if idx >= 0 {
				m.undoSnapshot()
				m.todos[idx].Description = val
			}
			m.input.SetValue("")
			m.mode = viewMode
		}
	case "esc":
		m.input.SetValue("")
		if m.mode == searchMode {
			m.searchQuery = ""
			m.cursor = 0
		}
		m.mode = viewMode
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) handleCategoryDeleteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.activeTab > 0 {
			m.undoSnapshot()
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

		}
		m.mode = viewMode
	case "esc":
		m.mode = viewMode
	}
	return m, nil
}

func (m model) handleViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if err := m.storage.Save(m.todos); err != nil {
			log.Printf("error: failed to save todos: %v", err)
		}
		return m, tea.Quit
	case "h":
		if m.activeTab > 0 {
			m.activeTab--
			m.cursor = 0
		}
	case "l":
		if m.activeTab < len(m.categories)-1 {
			m.activeTab++
			m.cursor = 0
		}
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < m.filteredCount()-1 {
			m.cursor++
		}
	case "enter":
		if m.filteredCount() > 0 {
			idx := m.getFilteredIndex(m.cursor)
			if idx >= 0 {
				m.undoSnapshot()
				m.todos[idx].Completed = !m.todos[idx].Completed
			}
		}
	case "a":
		m.mode = addMode
		m.input.Placeholder = "タスク名..."
		m.input.SetValue("")
		m.input.Focus()
		return m, textinput.Blink
	case "n":
		m.mode = categoryAddMode
		m.input.Placeholder = "カテゴリ名..."
		m.input.SetValue("")
		m.input.Focus()
		return m, textinput.Blink
	case "e":
		if m.filteredCount() > 0 {
			idx := m.getFilteredIndex(m.cursor)
			if idx >= 0 {
				m.mode = editMode
				m.input.SetValue(m.todos[idx].Title)
				m.input.Focus()
			}
		}
	case "i":
		if m.filteredCount() > 0 {
			idx := m.getFilteredIndex(m.cursor)
			if idx >= 0 {
				m.mode = descMode
				m.input.Placeholder = "詳細..."
				m.input.SetValue(m.todos[idx].Description)
				m.input.Focus()
			}
		}
	case "p":
		if m.filteredCount() > 0 {
			idx := m.getFilteredIndex(m.cursor)
			if idx >= 0 {
				m.undoSnapshot()
				m.todos[idx].Priority = (m.todos[idx].Priority + 1) % 3
			}
		}
	case "f":
		m.filterDone = !m.filterDone
	case "s":
		m.sortField = (m.sortField + 1) % 4
		m.sortTodos()
		m.cursor = 0
	case "x":
		m.mode = categoryDeleteMode
	case "/":
		m.mode = searchMode
		m.input.Placeholder = "検索..."
		m.input.SetValue(m.searchQuery)
		m.input.Focus()
		return m, textinput.Blink
	case "d":
		if m.filteredCount() > 0 {
			idx := m.getFilteredIndex(m.cursor)
			if idx >= 0 {
				m.undoSnapshot()
				m.todos = append(m.todos[:idx], m.todos[idx+1:]...)
				if m.cursor >= m.filteredCount() && m.cursor > 0 {
					m.cursor--
				}
			}
		}
	case "u":
		if len(m.lastTodos) > 0 {
			m.todos = m.lastTodos
			m.lastTodos = nil
			if m.cursor >= m.filteredCount() && m.cursor > 0 {
				m.cursor = m.filteredCount() - 1
			}
		}
	}
	return m, nil
}

func (m model) filteredCount() int {
	c := 0
	for _, t := range m.todos {
		if t.Category == m.categories[m.activeTab] && (!m.filterDone || !t.Completed) && (m.searchQuery == "" || containsIgnoreCase(t.Title, m.searchQuery)) {
			c++
		}
	}
	return c
}

func (m model) getFilteredIndex(target int) int {
	c := 0
	for i, t := range m.todos {
		if t.Category == m.categories[m.activeTab] && (!m.filterDone || !t.Completed) && (m.searchQuery == "" || containsIgnoreCase(t.Title, m.searchQuery)) {
			if c == target { return i }
			c++
		}
	}
	return -1
}

func main() {
	storage := NewFileStorage(DefaultFilePath())
	if _, err := tea.NewProgram(initialModel(storage), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v", err); os.Exit(1)
	}
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func (m *model) undoSnapshot() {
	m.lastTodos = make([]Todo, len(m.todos))
	copy(m.lastTodos, m.todos)
}
