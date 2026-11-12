package main

import (
	"math"
	"testing"
)

func TestTargetBirdAngleLimits(t *testing.T) {
	for _, test := range []struct {
		velocity float64
		want     float64
	}{
		{-10000, -12}, {-240, -12}, {0, 0},
		{120, 6}, {360, 18}, {10000, 18},
	} {
		if got := targetBirdAngle(test.velocity); got != test.want {
			t.Fatalf("targetBirdAngle(%v)=%v, want %v", test.velocity, got, test.want)
		}
	}
}

func TestBirdAngleSmoothsAndFreezes(t *testing.T) {
	w := testWorld()
	w.velocity = flapVelocity
	w.updateBirdAngle()
	if w.birdAngle >= 0 || w.birdAngle < -12 {
		t.Fatalf("first flap angle=%v, want [-12,0)", w.birdAngle)
	}
	first := w.birdAngle
	for range 240 {
		w.updateBirdAngle()
	}
	if math.Abs(w.birdAngle-targetBirdAngle(flapVelocity)) > 1e-6 {
		t.Fatalf("settled angle=%v, want %v", w.birdAngle, targetBirdAngle(flapVelocity))
	}
	w.previousBirdAngle = first
	interpolated := w.previousBirdAngle + (w.birdAngle-w.previousBirdAngle)*0.5
	if interpolated >= first || interpolated <= w.birdAngle {
		t.Fatalf("interpolated angle=%v, first=%v current=%v", interpolated, first, w.birdAngle)
	}
	w.state = paused
	before := w.birdAngle
	w.step(controls{focused: true})
	if w.birdAngle != before || w.previousBirdAngle != before {
		t.Fatalf("paused angles changed: previous=%v current=%v want %v", w.previousBirdAngle, w.birdAngle, before)
	}
	w.reset()
	if w.birdAngle != 0 || w.previousBirdAngle != 0 {
		t.Fatalf("reset angles=%v/%v", w.previousBirdAngle, w.birdAngle)
	}
}
