package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestFlappyGopherMetadata(t *testing.T) {
	data, err := os.ReadFile("mygo.json")
	if err != nil {
		t.Fatalf("read mygo.json: %v", err)
	}
	var config struct {
		Name       string `json:"name"`
		Slug       string `json:"slug"`
		Identifier string `json:"identifier"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("decode mygo.json: %v", err)
	}
	if config.Name != "Flappy Gopher" {
		t.Errorf("name=%q, want Flappy Gopher", config.Name)
	}
	if config.Slug != "flappy-gopher" {
		t.Errorf("slug=%q, want flappy-gopher", config.Slug)
	}
	if config.Identifier != "io.mzh.flappygopher" {
		t.Errorf("identifier=%q, want io.mzh.flappygopher", config.Identifier)
	}
}

func TestFlappyGopherWindowStateKey(t *testing.T) {
	if windowStateKey != "main-flappy-gopher-v1" {
		t.Fatalf("windowStateKey=%q, want main-flappy-gopher-v1", windowStateKey)
	}
}
