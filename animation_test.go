package main

import (
	"math"
	"testing"
)

func TestAnimationClockAdvancesOnlyDuringPlay(t *testing.T) {
	w := testWorld()
	w.step(controls{focused: true})
	if w.animationTime != 0 {
		t.Fatalf("ready animation time=%v", w.animationTime)
	}
	w.step(controls{flap: true, focused: true})
	if w.animationTime != fixedStep {
		t.Fatalf("playing animation time=%v", w.animationTime)
	}
	w.step(controls{pause: true, focused: true})
	if w.animationTime != fixedStep {
		t.Fatalf("paused animation advanced=%v", w.animationTime)
	}
	w.step(controls{pause: true, focused: true})
	w.step(controls{focused: true})
	if math.Abs(w.animationTime-2*fixedStep) > 1e-9 {
		t.Fatalf("resumed animation time=%v", w.animationTime)
	}
	w.reset()
	if w.animationTime != 0 {
		t.Fatalf("reset animation time=%v", w.animationTime)
	}
}

func TestGroundTravelAccumulatesAcrossTileWidth(t *testing.T) {
	w := testWorld()
	w.state = playing
	w.birdY, w.previousBirdY = 100, 100
	w.pipes = nil
	w.groundTravel = 31.9
	w.previousGroundTravel = 31.9
	w.step(controls{focused: true})
	want := 31.9 + pipeSpeed*fixedStep
	if math.Abs(w.groundTravel-want) > 1e-9 {
		t.Fatalf("ground travel=%v, want %v", w.groundTravel, want)
	}
}
