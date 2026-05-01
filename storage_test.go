package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "todotest")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func withTempFile(t *testing.T, fn func(string)) {
	t.Helper()
	dir := tempDir(t)
	origFilename := filename
	t.Cleanup(func() { filename = origFilename })

	tmpPath := filepath.Join(dir, "todos.json")
	filename = tmpPath
	fn(tmpPath)
}

func TestLoadTodos_FileNotFound(t *testing.T) {
	withTempFile(t, func(path string) {
		todos, err := loadTodos()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(todos) != 0 {
			t.Fatalf("expected empty slice, got %d todos", len(todos))
		}
	})
}

func TestLoadTodos_InvalidJSON(t *testing.T) {
	withTempFile(t, func(path string) {
		os.WriteFile(path, []byte("not json"), 0644)
		_, err := loadTodos()
		if err == nil {
			t.Fatal("expected error for invalid JSON, got nil")
		}
	})
}

func TestLoadTodos_ValidJSON(t *testing.T) {
	withTempFile(t, func(path string) {
		now := time.Now().Truncate(time.Second).UTC()
		data, _ := json.MarshalIndent([]Todo{
			{Title: "test", Completed: true, Priority: High, DueDate: now, Category: "Work"},
		}, "", "  ")
		os.WriteFile(path, data, 0644)

		todos, err := loadTodos()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("expected 1 todo, got %d", len(todos))
		}
		if todos[0].Title != "test" {
			t.Errorf("expected title 'test', got %q", todos[0].Title)
		}
		if !todos[0].Completed {
			t.Error("expected Completed=true")
		}
		if todos[0].Priority != High {
			t.Errorf("expected Priority=High, got %d", todos[0].Priority)
		}
		if todos[0].Category != "Work" {
			t.Errorf("expected Category='Work', got %q", todos[0].Category)
		}
	})
}

func TestSaveTodos(t *testing.T) {
	withTempFile(t, func(path string) {
		now := time.Now().Truncate(time.Second).UTC()
		todos := []Todo{
			{Title: "task1", Completed: false, Priority: Low, DueDate: now, Category: "Home"},
			{Title: "task2", Completed: true, Priority: Medium, DueDate: now, Category: "Work"},
		}
		if err := saveTodos(todos); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read saved file: %v", err)
		}

		var loaded []Todo
		if err := json.Unmarshal(raw, &loaded); err != nil {
			t.Fatalf("failed to unmarshal saved file: %v", err)
		}
		if len(loaded) != 2 {
			t.Fatalf("expected 2 todos, got %d", len(loaded))
		}
		if loaded[0].Title != "task1" || loaded[1].Title != "task2" {
			t.Errorf("titles mismatch: %q, %q", loaded[0].Title, loaded[1].Title)
		}
	})
}

func TestSaveAndLoadRoundtrip(t *testing.T) {
	withTempFile(t, func(path string) {
		now := time.Now().Truncate(time.Second).UTC()
		original := []Todo{
			{Title: "buy milk", Completed: false, Priority: High, DueDate: now, Category: "Shopping"},
		}
		if err := saveTodos(original); err != nil {
			t.Fatalf("save error: %v", err)
		}

		loaded, err := loadTodos()
		if err != nil {
			t.Fatalf("load error: %v", err)
		}
		if len(loaded) != 1 || loaded[0].Title != "buy milk" {
			t.Errorf("roundtrip mismatch: %+v", loaded)
		}
	})
}

func TestSaveTodos_CreatesParentDir(t *testing.T) {
	dir := tempDir(t)
	origFilename := filename
	t.Cleanup(func() { filename = origFilename })

	nested := filepath.Join(dir, "subdir1", "subdir2", "todos.json")
	filename = nested

	todos := []Todo{{Title: "nested", Category: "Home"}}
	if err := saveTodos(todos); err != nil {
		t.Fatalf("save error: %v", err)
	}
	raw, err := os.ReadFile(nested)
	if err != nil {
		t.Fatalf("failed to read nested file: %v", err)
	}
	if !strings.Contains(string(raw), "nested") {
		t.Error("saved content missing expected data")
	}
}
