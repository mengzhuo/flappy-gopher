package main

import "testing"

func TestPausedAndGameOverDoNotAdvancePipeSequence(t *testing.T) {
	for _, current := range []state{paused, gameOver} {
		w := testWorld()
		w.state = current
		before := append([]pipe(nil), w.pipes...)
		transition := w.nextTransition
		center := w.lastGapCenter
		for range 120 {
			w.step(activeControls())
		}
		if w.nextTransition != transition || w.lastGapCenter != center {
			t.Fatalf("state=%v advanced sequence", current)
		}
		for i := range before {
			if w.pipes[i] != before[i] {
				t.Fatalf("state=%v changed pipe %d", current, i)
			}
		}
		got, want := w.makePipe(0), testWorld().makePipe(0)
		if got != want {
			t.Fatalf("state=%v consumed gate RNG: got %+v want %+v", current, got, want)
		}
	}
}

func TestFocusLossDoesNotAdvancePipeSequence(t *testing.T) {
	w := testWorld()
	w.state = playing
	before := append([]pipe(nil), w.pipes...)
	transition := w.nextTransition
	w.step(controls{focused: false})
	if w.state != paused || w.nextTransition != transition {
		t.Fatalf("focus loss advanced sequence: state=%v transition=%d", w.state, w.nextTransition)
	}
	for i := range before {
		if w.pipes[i] != before[i] {
			t.Fatalf("focus loss changed pipe %d", i)
		}
	}
}

func TestPlayingResizePreservesAbsoluteTransitionOrder(t *testing.T) {
	wide, baseline := testWorld(), testWorld()
	wide.state, baseline.state = playing, playing
	wide.resize(910)
	for baseline.nextTransition < wide.nextTransition {
		baseline.makePipe(0)
	}
	got, want := wide.makePipe(0), baseline.makePipe(0)
	if got != want {
		t.Fatalf("next absolute transition differs: got %+v want %+v", got, want)
	}
}
