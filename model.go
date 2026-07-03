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
	Title       string    `json:"title"`
	Completed   bool      `json:"completed"`
	Priority    Priority  `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
}

type mode int

const (
	viewMode mode = iota
	addMode
	addDateMode
	editMode
	categoryAddMode
	categoryDeleteMode
	searchMode
	descMode
)

type model struct {
	todos        []Todo
	categories   []string
	activeTab    int
	cursor       int
	input        textinput.Model
	mode         mode
	filterDone   bool
	width        int
	height       int
	sortField    sortField
	storage      Storage
	pendingTitle string
	searchQuery  string
	lastTodos    []Todo
	config       Config
	ui           styles
}

const maxInputLen = 100

func (m *model) sortTodos() {
	if m.sortField == sortNone {
		return
	}
	activeCat := m.categories[m.activeTab]

	var activeItems []Todo
	for _, t := range m.todos {
		if t.Category == activeCat {
			activeItems = append(activeItems, t)
		}
	}

	sort.SliceStable(activeItems, func(i, j int) bool {
		switch m.sortField {
		case sortPriority:
			return activeItems[i].Priority > activeItems[j].Priority
		case sortDueDate:
			return activeItems[i].DueDate.Before(activeItems[j].DueDate)
		case sortName:
			return activeItems[i].Title < activeItems[j].Title
		default:
			return false
		}
	})

	idx := 0
	for i := range m.todos {
		if m.todos[i].Category == activeCat {
			m.todos[i] = activeItems[idx]
			idx++
		}
	}
}
