package main

import "testing"

func TestThemeCopy(t *testing.T) {
	if appTitle != "Flappy Gopher" || readyHint != "SPACE: go func()" {
		t.Fatalf("title=%q ready=%q", appTitle, readyHint)
	}
	for _, state := range []state{ready, playing} {
		if title, hint := stateOverlay(state); title != "" || hint != "" {
			t.Fatalf("state %v unexpectedly has overlay %q %q", state, title, hint)
		}
	}
	if title, hint := stateOverlay(paused); title != "runtime.Gosched()" || hint != "P / ESC TO RESUME" {
		t.Fatalf("paused overlay=%q %q", title, hint)
	}
	if title, hint := stateOverlay(gameOver); title != "STOP THE WORLD" || hint != "SPACE: restart goroutine" {
		t.Fatalf("game-over overlay=%q %q", title, hint)
	}
}
