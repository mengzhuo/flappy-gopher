package main

import (
	"reflect"
	"testing"
)

func TestFlapEventIsEmittedAndDrainedOnce(t *testing.T) {
	w := testWorld()
	w.step(controls{flap: true, focused: true})
	if got, want := w.drainEvents(), []soundEvent{eventFlap}; !reflect.DeepEqual(got, want) {
		t.Fatalf("flap events = %v, want %v", got, want)
	}
	if got := w.drainEvents(); len(got) != 0 {
		t.Fatalf("events not drained: %v", got)
	}
}

func TestScoreEventIsEmittedOnce(t *testing.T) {
	w := testWorld()
	w.state = playing
	w.birdY, w.previousBirdY = 100, 100
	w.pipes = []pipe{{x: birdX - pipeWidth, previousX: birdX - pipeWidth, gapY: 80}}
	w.step(controls{focused: true})
	if w.score != 1 || !reflect.DeepEqual(w.events, []soundEvent{eventScore}) {
		t.Fatalf("score=%d events=%v", w.score, w.events)
	}
}

func TestHitEventIsNotRepeatedAndResetClearsEvents(t *testing.T) {
	w := testWorld()
	w.state = playing
	w.birdY = groundY - birdRadius
	w.step(controls{focused: true})
	w.finish()
	if !reflect.DeepEqual(w.events, []soundEvent{eventHit}) {
		t.Fatalf("hit events = %v", w.events)
	}
	w.reset()
	if len(w.events) != 0 {
		t.Fatalf("reset retained events: %v", w.events)
	}
}
