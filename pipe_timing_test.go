package main

import (
	"math"
	"testing"
)

func TestFirstPipeArrivesAtFourSecondsAcrossWidths(t *testing.T) {
	for _, width := range []float64{288, 420, 910, 1400} {
		w := testWorld()
		w.resize(width)
		w.state = playing
		if w.pipes[0].x != firstPipeX || w.pipes[0].gapY != logicalHeight/2-pipeGap/2 {
			t.Fatalf("width=%v initial pipe=%+v", width, w.pipes[0])
		}
		for frame := 1; frame <= 240; frame++ {
			in := activeControls()
			if w.birdY > logicalHeight/2 {
				in.flap = true
			}
			w.step(in)
			if w.state != playing {
				t.Fatalf("width=%v frame=%d state=%v", width, frame, w.state)
			}
			if frame == 239 && w.pipes[0].x <= birdX {
				t.Fatalf("width=%v first pipe arrived early: x=%v", width, w.pipes[0].x)
			}
		}
		if math.Abs(w.pipes[0].x-birdX) > 1e-9 {
			t.Fatalf("width=%v first pipe x=%v after 4 seconds, want %v", width, w.pipes[0].x, birdX)
		}
	}
}
