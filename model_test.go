package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestFilteredCount_NoFilter(t *testing.T) {
	m := model{
		categories: []string{"Home", "Work"},
		activeTab:  0,
		todos: []Todo{
			{Title: "a", Category: "Home", Completed: false},
			{Title: "b", Category: "Home", Completed: true},
			{Title: "c", Category: "Work", Completed: false},
		},
	}
	if got := m.filteredCount(); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestFilteredCount_FilterDone(t *testing.T) {
	m := model{
		categories: []string{"Home"},
		activeTab:  0,
		filterDone: true,
		todos: []Todo{
			{Title: "a", Category: "Home", Completed: false},
			{Title: "b", Category: "Home", Completed: true},
			{Title: "c", Category: "Home", Completed: true},
		},
	}
	if got := m.filteredCount(); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestFilteredCount_Empty(t *testing.T) {
	m := model{
		categories: []string{"Home"},
		activeTab:  0,
		todos:      []Todo{},
	}
	if got := m.filteredCount(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestFilteredCount_NoMatchingCategory(t *testing.T) {
	m := model{
		categories: []string{"Home", "Work"},
		activeTab:  1,
		todos: []Todo{
			{Title: "a", Category: "Home"},
		},
	}
	if got := m.filteredCount(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestGetFilteredIndex(t *testing.T) {
	m := model{
		categories: []string{"Home"},
		activeTab:  0,
		todos: []Todo{
			{Title: "a", Category: "Home", Completed: false},
			{Title: "b", Category: "Home", Completed: false},
			{Title: "c", Category: "Home", Completed: false},
		},
	}
	tests := []struct {
		target int
		want   int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
	}
	for _, tt := range tests {
		got := m.getFilteredIndex(tt.target)
		if got != tt.want {
			t.Errorf("getFilteredIndex(%d) = %d, want %d", tt.target, got, tt.want)
		}
	}
}

func TestGetFilteredIndex_WithFilter(t *testing.T) {
	m := model{
		categories: []string{"Home"},
		activeTab:  0,
		filterDone: true,
		todos: []Todo{
			{Title: "a", Category: "Home", Completed: false},
			{Title: "b", Category: "Home", Completed: true},
			{Title: "c", Category: "Home", Completed: false},
		},
	}
	tests := []struct {
		target int
		want   int
	}{
		{0, 0},
		{1, 2},
	}
	for _, tt := range tests {
		got := m.getFilteredIndex(tt.target)
		if got != tt.want {
			t.Errorf("getFilteredIndex(%d) = %d, want %d", tt.target, got, tt.want)
		}
	}
}

func TestGetFilteredIndex_NotFound(t *testing.T) {
	m := model{
		categories: []string{"Home"},
		activeTab:  0,
		todos:      []Todo{},
	}
	if got := m.getFilteredIndex(0); got != -1 {
		t.Errorf("expected -1, got %d", got)
	}
}

func TestGetFilteredIndex_WithCategories(t *testing.T) {
	m := model{
		categories: []string{"Home", "Work"},
		activeTab:  1,
		todos: []Todo{
			{Title: "a", Category: "Home"},
			{Title: "b", Category: "Work"},
			{Title: "c", Category: "Work"},
		},
	}
	tests := []struct {
		target int
		want   int
	}{
		{0, 1},
		{1, 2},
	}
	for _, tt := range tests {
		got := m.getFilteredIndex(tt.target)
		if got != tt.want {
			t.Errorf("getFilteredIndex(%d) = %d, want %d", tt.target, got, tt.want)
		}
	}
}

func TestPriorityConstants(t *testing.T) {
	if Low != 0 {
		t.Errorf("Low should be 0, got %d", Low)
	}
	if Medium != 1 {
		t.Errorf("Medium should be 1, got %d", Medium)
	}
	if High != 2 {
		t.Errorf("High should be 2, got %d", High)
	}
}

func TestTodoStruct_JSONTags(t *testing.T) {
	todo := Todo{
		Title:     "test",
		Completed: true,
		Priority:  High,
		DueDate:   time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		Category:  "Work",
	}
	expected := `{"title":"test","completed":true,"priority":2,"due_date":"2025-01-15T00:00:00Z","category":"Work"}`

	data, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != expected {
		t.Errorf("JSON mismatch:\ngot:      %s\nexpected: %s", string(data), expected)
	}
}

func TestMaxContentWidth_Default(t *testing.T) {
	m := model{width: 0}
	if got := maxContentWidth(m); got != 80 {
		t.Errorf("expected 80 for zero width, got %d", got)
	}
}

func TestMaxContentWidth_Small(t *testing.T) {
	m := model{width: 20}
	if got := maxContentWidth(m); got != 16 {
		t.Errorf("expected 16 for width 20, got %d", got)
	}
}

func TestMaxContentWidth_Large(t *testing.T) {
	m := model{width: 120}
	if got := maxContentWidth(m); got != 116 {
		t.Errorf("expected 116 for width 120, got %d", got)
	}
}

func TestView_OverdueTodo(t *testing.T) {
	m := model{
		width:      80,
		height:     24,
		categories: []string{"Home"},
		activeTab:  0,
		todos: []Todo{
			{Title: "overdue task", Category: "Home", Completed: false, DueDate: time.Now().AddDate(0, 0, -1)},
		},
	}
	view := m.View()
	if !strings.Contains(view, "overdue task") {
		t.Error("view should contain the overdue task title")
	}
}

func TestView_CompletedTodo(t *testing.T) {
	m := model{
		width:      80,
		height:     24,
		categories: []string{"Home"},
		activeTab:  0,
		todos: []Todo{
			{Title: "done task", Category: "Home", Completed: true, DueDate: time.Now().AddDate(0, 0, 1)},
		},
	}
	view := m.View()
	if !strings.Contains(view, "done task") {
		t.Error("view should contain the completed task title")
	}
}

func TestView_LongTitleTruncated(t *testing.T) {
	m := model{
		width:      50,
		height:     24,
		categories: []string{"Home"},
		activeTab:  0,
		todos: []Todo{
			{Title: "this is a very long title that should be truncated in the view output", Category: "Home", Completed: false, DueDate: time.Now().AddDate(0, 0, 1)},
		},
	}
	view := m.View()
	if !strings.Contains(view, "...") {
		t.Error("long title should be truncated with '...'")
	}
}

func TestView_TooManyTodosEllipsis(t *testing.T) {
	todos := make([]Todo, 10)
	now := time.Now().AddDate(0, 0, 1)
	for i := range todos {
		todos[i] = Todo{Title: "task", Category: "Home", Completed: false, DueDate: now}
	}
	m := model{
		width:      80,
		height:     12,
		categories: []string{"Home"},
		activeTab:  0,
		todos:      todos,
	}
	view := m.View()
	if !strings.Contains(view, "...") {
		t.Error("view should show '...' when todos exceed available height")
	}
}

func TestMaxInputLen(t *testing.T) {
	if maxInputLen <= 0 {
		t.Errorf("maxInputLen should be positive, got %d", maxInputLen)
	}
}

func TestSortTodos_ByPriority(t *testing.T) {
	m := &model{
		categories: []string{"Home"},
		activeTab:  0,
		sortField:  sortPriority,
		todos: []Todo{
			{Title: "low task", Category: "Home", Priority: Low},
			{Title: "high task", Category: "Home", Priority: High},
			{Title: "mid task", Category: "Home", Priority: Medium},
		},
	}
	m.sortTodos()
	if m.todos[0].Priority != High {
		t.Errorf("first should be High, got %d", m.todos[0].Priority)
	}
	if m.todos[1].Priority != Medium {
		t.Errorf("second should be Medium, got %d", m.todos[1].Priority)
	}
	if m.todos[2].Priority != Low {
		t.Errorf("third should be Low, got %d", m.todos[2].Priority)
	}
}

func TestSortTodos_ByDueDate(t *testing.T) {
	m := &model{
		categories: []string{"Home"},
		activeTab:  0,
		sortField:  sortDueDate,
		todos: []Todo{
			{Title: "late", Category: "Home", DueDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)},
			{Title: "early", Category: "Home", DueDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
			{Title: "mid", Category: "Home", DueDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)},
		},
	}
	m.sortTodos()
	if m.todos[0].Title != "early" {
		t.Errorf("first should be early, got %s", m.todos[0].Title)
	}
	if m.todos[2].Title != "late" {
		t.Errorf("last should be late, got %s", m.todos[2].Title)
	}
}

func TestSortTodos_ByName(t *testing.T) {
	m := &model{
		categories: []string{"Home"},
		activeTab:  0,
		sortField:  sortName,
		todos: []Todo{
			{Title: "charlie", Category: "Home"},
			{Title: "alpha", Category: "Home"},
			{Title: "bravo", Category: "Home"},
		},
	}
	m.sortTodos()
	if m.todos[0].Title != "alpha" {
		t.Errorf("first should be alpha, got %s", m.todos[0].Title)
	}
	if m.todos[1].Title != "bravo" {
		t.Errorf("second should be bravo, got %s", m.todos[1].Title)
	}
	if m.todos[2].Title != "charlie" {
		t.Errorf("third should be charlie, got %s", m.todos[2].Title)
	}
}

func TestSortTodos_None(t *testing.T) {
	original := []Todo{
		{Title: "b", Category: "Home", Priority: Low},
		{Title: "a", Category: "Home", Priority: High},
	}
	m := &model{
		categories: []string{"Home"},
		activeTab:  0,
		sortField:  sortNone,
		todos:      append([]Todo{}, original...),
	}
	m.sortTodos()
	if m.todos[0].Title != "b" || m.todos[1].Title != "a" {
		t.Error("sortNone should not reorder")
	}
}

func TestSortTodos_OtherCategoryUnaffected(t *testing.T) {
	m := &model{
		categories: []string{"Home", "Work"},
		activeTab:  1,
		sortField:  sortPriority,
		todos: []Todo{
			{Title: "home low", Category: "Home", Priority: Low},
			{Title: "work low", Category: "Work", Priority: Low},
			{Title: "work high", Category: "Work", Priority: High},
		},
	}
	m.sortTodos()
	if m.todos[0].Category != "Home" {
		t.Error("non-active category items should stay in place")
	}
}

func TestSortLabel(t *testing.T) {
	tests := []struct {
		sf  sortField
		exp string
	}{
		{sortNone, "なし"},
		{sortPriority, "優先度"},
		{sortDueDate, "期限"},
		{sortName, "名前"},
	}
	for _, tt := range tests {
		if got := sortLabel(tt.sf); got != tt.exp {
			t.Errorf("sortLabel(%v) = %q, want %q", tt.sf, got, tt.exp)
		}
	}
}

func TestRenderProgressBar(t *testing.T) {
	m := model{ui: buildStyles(defaultTheme)}
	bar := renderProgressBar(m, 10, 5, 10)
	if bar == "" {
		t.Error("progress bar should not be empty")
	}

	barEmpty := renderProgressBar(m, 0, 0, 10)
	if barEmpty == "" {
		t.Error("progress bar for zero total should not be empty")
	}

	barFull := renderProgressBar(m, 10, 10, 10)
	if barFull == "" {
		t.Error("progress bar for full should not be empty")
	}
}

func TestView_ProgressBar(t *testing.T) {
	m := model{
		width:      80,
		height:     24,
		categories: []string{"Home"},
		activeTab:  0,
		todos: []Todo{
			{Title: "a", Category: "Home", Completed: true, DueDate: time.Now()},
			{Title: "b", Category: "Home", Completed: false, DueDate: time.Now()},
			{Title: "c", Category: "Home", Completed: true, DueDate: time.Now()},
		},
	}
	view := m.View()
	if !strings.Contains(view, "2/3") {
		t.Error("view should show progress as 2/3")
	}
}

func TestView_CategoryDeleteConfirm(t *testing.T) {
	m := model{
		width:      80,
		height:     24,
		categories: []string{"Home", "Work"},
		activeTab:  1,
		mode:       categoryDeleteMode,
	}
	view := m.View()
	if !strings.Contains(view, "Work") {
		t.Error("delete confirm should show category name")
	}
	if !strings.Contains(view, "確定") {
		t.Error("delete confirm should show confirmation prompt")
	}
}

func TestView_CategoryDeleteHome(t *testing.T) {
	m := model{
		width:      80,
		height:     24,
		categories: []string{"Home"},
		activeTab:  0,
		mode:       categoryDeleteMode,
	}
	view := m.View()
	if !strings.Contains(view, "削除できません") {
		t.Error("Home should not be deletable")
	}
}
