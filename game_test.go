package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/plugins/game"
	"github.com/egoist/mygo/ui"
)

func testGameApp(t *testing.T, w *world) *gameApp {
	t.Helper()
	assets, err := loadGameAssets()
	if err != nil {
		t.Fatal(err)
	}
	return newGameApp(w, scoreStore{path: filepath.Join(t.TempDir(), "score.json")}, assets)
}

func TestMuteKeyDoesNotChangeGameplay(t *testing.T) {
	w := testWorld()
	g := testGameApp(t, w)
	tt := ui.NewTester(g.view, 1152, 648)
	step := func() { tt.Advance(3 * time.Second / game.Rate) }
	tt.Key(0, ui.KeyM)
	step()
	if !g.muted || w.state != ready || w.score != 0 || w.velocity != 0 {
		t.Fatalf("muted=%v world=%+v", g.muted, w)
	}
	tt.Key(0, ui.KeyM)
	step()
	if g.muted || w.state != ready || w.score != 0 || w.velocity != 0 {
		t.Fatalf("muted=%v world=%+v", g.muted, w)
	}
}

func TestNativeViewInputAndFocusPause(t *testing.T) {
	w := testWorld()
	g := testGameApp(t, w)
	tt := ui.NewTester(g.view, 1152, 648)
	step := func() { tt.Advance(3 * time.Second / game.Rate) }
	tt.Key(0, ui.KeySpace)
	step()
	if w.state != playing || w.velocity == 0 {
		t.Fatalf("Space did not start: %+v", w)
	}
	tt.Key(0, ui.KeyP)
	step()
	if w.state != paused {
		t.Fatalf("P did not pause: %v", w.state)
	}
	tt.Key(0, ui.KeyP)
	step()
	if w.state != playing {
		t.Fatalf("P did not resume: %v", w.state)
	}
	tt.Key(0, ui.KeyR)
	step()
	if w.state != ready || w.score != 0 {
		t.Fatalf("R did not reset: %+v", w)
	}
}

func TestPointerAndFocusPause(t *testing.T) {
	w := testWorld()
	g := testGameApp(t, w)
	tt := ui.NewTester(g.view, 1152, 648)
	step := func() { tt.Advance(3 * time.Second / game.Rate) }
	tt.ClickAt(576, 324)
	step()
	if w.state != playing {
		t.Fatalf("click did not start: %v", w.state)
	}
	tt.SetFocused(false)
	step()
	if w.state != paused {
		t.Fatalf("focus loss did not pause: %v", w.state)
	}
	tt.SetFocused(true)
	step()
	if w.state != paused {
		t.Fatal("focus return auto-resumed")
	}
	tt.Key(0, ui.KeyEscape)
	step()
	if w.state != playing {
		t.Fatalf("Escape did not resume: %v", w.state)
	}
}

func TestNaturalGameOverAndReset(t *testing.T) {
	w := testWorld()
	g := testGameApp(t, w)
	tt := ui.NewTester(g.view, 1152, 648)
	step := func() { tt.Advance(3 * time.Second / game.Rate) }
	tt.Key(0, ui.KeySpace)
	step()
	for range 180 {
		step()
	}
	if w.state != gameOver {
		t.Fatalf("bird never collided: %+v", w)
	}
	tt.Key(0, ui.KeyR)
	step()
	if w.state != ready || w.score != 0 {
		t.Fatalf("reset after collision failed: %+v", w)
	}
}

func TestConsumeEventsAndMute(t *testing.T) {
	w := testWorld()
	played := make(map[soundEvent]int)
	g := &gameApp{world: w, effects: make(map[soundEvent]func())}
	for _, event := range []soundEvent{eventFlap, eventScore, eventHit} {
		e := event
		g.effects[e] = func() { played[e]++ }
	}
	w.events = []soundEvent{eventFlap, eventScore}
	g.consumeEvents()
	if played[eventFlap] != 1 || played[eventScore] != 1 || len(w.events) != 0 {
		t.Fatalf("played=%v events=%v", played, w.events)
	}
	g.muted = true
	w.events = []soundEvent{eventHit}
	g.consumeEvents()
	if played[eventHit] != 0 || len(w.events) != 0 {
		t.Fatalf("muted played=%v", played)
	}
	g.muted = false
	w.events = []soundEvent{eventHit}
	g.consumeEvents()
	if played[eventHit] != 1 {
		t.Fatalf("unmuted hit count=%d", played[eventHit])
	}
}
