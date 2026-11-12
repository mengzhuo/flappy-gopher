package main

import "testing"

func scoreOneGate(t *testing.T, w *world) {
	t.Helper()
	w.state = playing
	w.birdY, w.previousBirdY = 100, 100
	w.velocity = -100
	w.pipes = []pipe{{
		x: birdX - pipeWidth, previousX: birdX - pipeWidth, gapY: 80,
	}}
	w.step(activeControls())
	if w.state != playing {
		t.Fatalf("scoring fixture collided: %+v", w)
	}
}

func TestScoreFeedbackExpiresAndResetClearsIt(t *testing.T) {
	w := testWorld()
	scoreOneGate(t, w)
	if w.score != 1 || w.feedback != "GC OK +1" || w.feedbackTime != scoreFeedbackDuration {
		t.Fatalf("score=%d feedback=%q time=%v", w.score, w.feedback, w.feedbackTime)
	}
	w.pipes = nil
	for range 61 {
		w.step(activeControls())
	}
	if w.feedback != "" || w.feedbackTime != 0 {
		t.Fatalf("feedback did not expire: %q %v", w.feedback, w.feedbackTime)
	}
	scoreOneGate(t, w)
	w.reset()
	if w.feedback != "" || w.feedbackTime != 0 {
		t.Fatalf("reset retained feedback: %q %v", w.feedback, w.feedbackTime)
	}
}

func TestEveryFifthGateUsesMemoryJokeWithoutBonusScore(t *testing.T) {
	w := testWorld()
	w.score = 4
	scoreOneGate(t, w)
	if w.score != 5 || w.feedback != "128 MB freed" {
		t.Fatalf("score=%d feedback=%q", w.score, w.feedback)
	}
	w.finish()
	if w.feedback != "" || w.feedbackTime != 0 {
		t.Fatalf("game over retained feedback: %q %v", w.feedback, w.feedbackTime)
	}
}
