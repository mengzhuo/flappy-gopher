package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScoreStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "score.json")
	s := scoreStore{path: path}
	if best, err := s.load(); err != nil || best != 0 {
		t.Fatalf("missing file: best=%d err=%v", best, err)
	}
	if err := s.save(12); err != nil {
		t.Fatal(err)
	}
	if best, err := s.load(); err != nil || best != 12 {
		t.Fatalf("round trip: best=%d err=%v", best, err)
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil || len(entries) != 1 {
		t.Fatalf("temporary files left behind: %v %v", entries, err)
	}
}

func TestScoreStoreRejectsCorruptionAndNegatives(t *testing.T) {
	path := filepath.Join(t.TempDir(), "score.json")
	if err := os.WriteFile(path, []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	s := scoreStore{path: path}
	if best, err := s.load(); err == nil || best != 0 {
		t.Fatalf("corrupt file: best=%d err=%v", best, err)
	}
	if err := s.save(-1); err == nil {
		t.Fatal("negative high score was accepted")
	}
}
