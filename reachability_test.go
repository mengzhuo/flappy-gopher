package main

import (
	"math"
	"math/rand/v2"
	"testing"
)

func referenceControls(w *world) controls {
	in := activeControls()
	if w.state == ready {
		in.flap = true
		return in
	}
	target := logicalHeight / 2
	current := -1
	for i := range w.pipes {
		if w.pipes[i].x+pipeWidth >= birdX-birdRadius {
			current = i
			break
		}
	}
	if current >= 0 {
		center := w.pipes[current].gapY + pipeGap/2
		target = center
		if current+1 < len(w.pipes) {
			nextCenter := w.pipes[current+1].gapY + pipeGap/2
			distance := math.Max(0, w.pipes[current].x-birdX)
			weight := math.Max(0, math.Min(1, (pipeSpacing-distance)/pipeSpacing)) * 0.15
			target = center*(1-weight) + nextCenter*weight
		}
	}
	velocity := math.Min(w.velocity+verticalAcceleration(w.velocity)*fixedStep, maxFallVelocity)
	predictedY := w.birdY + velocity*fixedStep
	in.flap = predictedY > target+3
	return in
}

func TestReferenceControllerClearsProgressiveSequences(t *testing.T) {
	for _, seed := range []uint64{1, 7, 42, 99} {
		w := newWorld(0, rand.New(rand.NewPCG(seed, seed+1)))
		for steps := 0; steps < 60*190 && w.score < 100; steps++ {
			w.step(referenceControls(w))
			if w.state == gameOver {
				t.Fatalf("seed=%d score=%d birdY=%v velocity=%v pipes=%+v", seed, w.score, w.birdY, w.velocity, w.pipes)
			}
		}
		if w.score < 100 {
			t.Fatalf("seed=%d only scored %d", seed, w.score)
		}
	}
}
