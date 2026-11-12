package main

import (
	"math/rand/v2"
	"testing"
)

func testWorld() *world {
	return newWorld(7, rand.New(rand.NewPCG(1, 2)))
}

func TestResetRestoresReadyState(t *testing.T) {
	w := testWorld()
	if w.state != ready || w.score != 0 || w.best != 7 {
		t.Fatalf("initial world: %+v", w)
	}
	if w.birdY != logicalHeight/2 || w.velocity != 0 || len(w.pipes) != 3 {
		t.Fatalf("initial bird/pipes: y=%v v=%v pipes=%d", w.birdY, w.velocity, len(w.pipes))
	}
	w.birdY, w.score, w.state = 100, 9, gameOver
	w.reset()
	if w.state != ready || w.score != 0 || w.best != 7 {
		t.Fatalf("reset did not restore state: %+v", w)
	}
	if w.birdY != logicalHeight/2 || w.previousBirdY != logicalHeight/2 || w.velocity != 0 || len(w.pipes) != 3 {
		t.Fatalf("reset did not restore bird/pipes: y=%v py=%v v=%v pipes=%d", w.birdY, w.previousBirdY, w.velocity, len(w.pipes))
	}
}

func TestResetRestoresAbsolutePipeSequence(t *testing.T) {
	w := testWorld()
	for range 8 {
		w.makePipe(0)
	}
	w.reset()
	fresh := testWorld()
	if w.nextTransition != fresh.nextTransition || w.lastGapCenter != fresh.lastGapCenter {
		t.Fatalf("sequence state=%d/%v, want %d/%v", w.nextTransition, w.lastGapCenter, fresh.nextTransition, fresh.lastGapCenter)
	}
	for i := range fresh.pipes {
		if w.pipes[i] != fresh.pipes[i] {
			t.Fatalf("pipe %d=%+v, want %+v", i, w.pipes[i], fresh.pipes[i])
		}
	}
	if w.pipes[0].actionBudget != 0 || w.pipes[1].actionBudget != 1 {
		t.Fatalf("initial budgets=%d/%d", w.pipes[0].actionBudget, w.pipes[1].actionBudget)
	}
}
