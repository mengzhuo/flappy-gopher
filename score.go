package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type scoreStore struct{ path string }

type scoreDocument struct {
	Best int `json:"best"`
}

func (s scoreStore) load() (int, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	var document scoreDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return 0, err
	}
	if document.Best < 0 {
		return 0, fmt.Errorf("negative high score %d", document.Best)
	}
	return document.Best, nil
}

func (s scoreStore) save(best int) error {
	if best < 0 {
		return fmt.Errorf("negative high score %d", best)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(s.path), ".score-*.tmp")
	if err != nil {
		return err
	}
	temporary := file.Name()
	defer os.Remove(temporary)
	if err := json.NewEncoder(file).Encode(scoreDocument{Best: best}); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temporary, s.path)
}
