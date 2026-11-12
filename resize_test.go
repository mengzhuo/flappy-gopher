package main

import (
	"math"
	"testing"
)

func TestReadyResizeKeepsFirstArrivalAndRebuildsCoverage(t *testing.T) {
	w := testWorld()
	w.resize(910)
	if w.width != 910 || len(w.pipes) != 6 {
		t.Fatalf("width=%v pipes=%d, want 910 and 6", w.width, len(w.pipes))
	}
	if w.pipes[0].x != firstPipeX {
		t.Fatalf("first gate x=%v, want %v", w.pipes[0].x, firstPipeX)
	}
}

func TestPlayingResizePreservesExistingGatesAndExtendsCoverage(t *testing.T) {
	w := testWorld()
	w.state = playing
	original := append([]pipe(nil), w.pipes...)
	w.resize(910)
	if len(w.pipes) <= len(original) {
		t.Fatalf("gate count did not grow: %d -> %d", len(original), len(w.pipes))
	}
	for i := range original {
		if w.pipes[i] != original[i] {
			t.Fatalf("existing gate %d changed: got %+v want %+v", i, w.pipes[i], original[i])
		}
	}
	for i := 1; i < len(w.pipes); i++ {
		spacing := w.pipes[i].x - w.pipes[i-1].x
		if math.Abs(spacing-pipeSpacing) > 1e-9 {
			t.Fatalf("spacing %d=%v, want %v", i, spacing, pipeSpacing)
		}
	}
	if math.Abs(pipeSpacing/pipeSpeed-1.6) > 1e-9 {
		t.Fatalf("gate interval=%v, want 1.6 seconds", pipeSpacing/pipeSpeed)
	}
}

func TestResizeRejectsInvalidWidthWithoutMovingWorld(t *testing.T) {
	w := testWorld()
	before := w.width
	for _, width := range []float64{0, -10, math.NaN(), math.Inf(1)} {
		w.resize(width)
		if w.width != before {
			t.Fatalf("resize(%v) changed width to %v", width, w.width)
		}
	}
}
