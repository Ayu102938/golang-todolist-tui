package main

import (
	"time"
	"github.com/charmbracelet/bubbles/textinput"
)

type Priority int
const ( Low Priority = iota; Medium; High )

type Todo struct {
	Title      string    `json:"title"`
	Completed  bool      `json:"completed"`
	Priority   Priority  `json:"priority"`
	DueDate    time.Time `json:"due_date"`
	Category   string    `json:"category"`
}

type mode int
const ( 
    viewMode mode = iota
    addMode
    editMode
    categoryAddMode // ここが不足していました
)

type model struct {
	todos       []Todo
	categories  []string
	activeTab   int
	cursor      int
	input       textinput.Model
	mode        mode
	width       int
	height      int
	filterDone  bool
}
