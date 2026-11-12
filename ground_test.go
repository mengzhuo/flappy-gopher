package main

import (
	"fmt"
	"math"
	"testing"
)

func TestVisibleHeapAddressesContinueAcrossTileBoundary(t *testing.T) {
	initial := visibleHeapBlocks(0, 96)
	if len(initial) != 3 {
		t.Fatalf("initial blocks=%d, want 3", len(initial))
	}
	for i, want := range []uint64{0x00A0, 0x00C0, 0x00E0} {
		if initial[i].address != want || initial[i].x != float64(i*32) {
			t.Fatalf("block %d=%+v, want address %#04X x=%d", i, initial[i], want, i*32)
		}
	}
	next := visibleHeapBlocks(32, 96)
	for i, want := range []uint64{0x00C0, 0x00E0, 0x0100} {
		if next[i].address != want {
			t.Fatalf("next block %d=%#04X, want %#04X", i, next[i].address, want)
		}
	}
}

func TestVisibleHeapBlocksInterpolateAndResizeWithoutReset(t *testing.T) {
	blocks := visibleHeapBlocks(31.9, 96)
	if math.Abs(blocks[0].x+31.9) > 1e-9 || blocks[0].address != 0x00A0 {
		t.Fatalf("partial first block=%+v", blocks[0])
	}
	if math.Abs(blocks[1].x-0.1) > 1e-9 || blocks[1].address != 0x00C0 {
		t.Fatalf("incoming block=%+v", blocks[1])
	}
	wide := visibleHeapBlocks(31.9, 160)
	for i := range blocks {
		if wide[i] != blocks[i] {
			t.Fatalf("resize changed block %d: %+v -> %+v", i, blocks[i], wide[i])
		}
	}
}

func TestHeapAddressFormattingExtendsBeyondFourDigits(t *testing.T) {
	if got := fmt.Sprintf("0x%04X", heapAddress(0x1000)); got != "0x200A0" {
		t.Fatalf("formatted address=%q, want 0x200A0", got)
	}
}

func TestGroundTravelPausesAndResets(t *testing.T) {
	w := testWorld()
	w.state = playing
	w.birdY, w.previousBirdY = 100, 100
	w.pipes = nil
	w.step(activeControls())
	if w.groundTravel != pipeSpeed*fixedStep {
		t.Fatalf("travel=%v", w.groundTravel)
	}
	w.state = paused
	before := w.groundTravel
	w.step(activeControls())
	if w.groundTravel != before {
		t.Fatalf("paused travel=%v, want %v", w.groundTravel, before)
	}
	w.reset()
	if w.groundTravel != 0 || w.previousGroundTravel != 0 {
		t.Fatalf("reset travel=%v previous=%v", w.groundTravel, w.previousGroundTravel)
	}
}
