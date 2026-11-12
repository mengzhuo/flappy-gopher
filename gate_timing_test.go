package main

import "testing"

func TestGateTimingsAtReferenceHeights(t *testing.T) {
	for _, test := range []struct {
		name             string
		center           float64
		wantMax, wantMin int
	}{
		{name: "highest", center: minPipeCenter, wantMax: 6, wantMin: 1},
		{name: "centered", center: logicalHeight / 2, wantMax: 12, wantMin: 3},
		{name: "lowest", center: maxPipeCenter, wantMax: 14, wantMin: 4},
	} {
		t.Run(test.name, func(t *testing.T) {
			gotMax, gotMin := gateTimings(test.center - pipeGap/2)
			if gotMax != test.wantMax || gotMin != test.wantMin {
				t.Fatalf("gateTimings(%v)=%d/%d, want %d/%d", test.center, gotMax, gotMin, test.wantMax, test.wantMin)
			}
		})
	}
}

func TestGateTimingsAreMonotonicAndOrdered(t *testing.T) {
	previousMax, previousMin := 0, 0
	for center := minPipeCenter; center <= maxPipeCenter+1e-9; center += 0.25 {
		gotMax, gotMin := gateTimings(center - pipeGap/2)
		if gotMax < previousMax || gotMin < previousMin {
			t.Fatalf("timing decreased at center=%v: %d/%d after %d/%d", center, gotMax, gotMin, previousMax, previousMin)
		}
		if gotMax < 6 || gotMax > 14 || gotMin < 1 || gotMin > 4 || gotMax < gotMin {
			t.Fatalf("invalid timing at center=%v: %d/%d", center, gotMax, gotMin)
		}
		previousMax, previousMin = gotMax, gotMin
	}
}
