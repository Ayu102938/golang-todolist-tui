package main

import (
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

type Priority int

const (
	Low    Priority = iota
	Medium
	High
)

type sortField int

const (
	sortNone sortField = iota
	sortPriority
	sortDueDate
	sortName
)

type Todo struct {
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Priority  Priority  `json:"priority"`
	DueDate   time.Time `json:"due_date"`
	Category  string    `json:"category"`
}

type mode int

const (
	viewMode mode = iota
	addMode
	editMode
	categoryAddMode
	categoryDeleteMode
)

type model struct {
	todos       []Todo
	categories  []string
	activeTab   int
	cursor      int
	input       textinput.Model
	mode        mode
	filterDone  bool
	width       int
	height      int
	sortField   sortField
}

const maxInputLen = 100

func (m *model) sortTodos() {
	if m.sortField == sortNone {
		return
	}
	activeCat := m.categories[m.activeTab]
	sort.SliceStable(m.todos, func(i, j int) bool {
		if m.todos[i].Category != activeCat || m.todos[j].Category != activeCat {
			return false
		}
		switch m.sortField {
		case sortPriority:
			return m.todos[i].Priority > m.todos[j].Priority
		case sortDueDate:
			return m.todos[i].DueDate.Before(m.todos[j].DueDate)
		case sortName:
			return m.todos[i].Title < m.todos[j].Title
		default:
			return false
		}
	})
}
