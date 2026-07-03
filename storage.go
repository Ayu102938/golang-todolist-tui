package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Storage interface {
	Load() ([]Todo, error)
	Save(todos []Todo) error
}

type FileStorage struct {
	path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

func DefaultFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	return filepath.Join(configDir, "my-tui-app", "todos.json")
}

func (fs *FileStorage) Load() ([]Todo, error) {
	file, err := os.ReadFile(fs.path)
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

func (fs *FileStorage) Save(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(fs.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(fs.path, data, 0644)
}
