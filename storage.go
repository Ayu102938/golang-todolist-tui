package main

import (
	"encoding/json"
	"os"
)

const filename = "todos.json"

func loadTodos() ([]Todo, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Todo{}, nil
		}
		return nil, err
	}
	var todos []Todo
	err = json.Unmarshal(file, &todos)
	return todos, err
}

func saveTodos(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
