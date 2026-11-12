package main

import (
	"math"
	"testing"
)

func TestPropellerAngleUsesAnimationClock(t *testing.T) {
	for _, test := range []struct {
		time float64
		want float64
	}{
		{0, 0}, {0.125, 90}, {0.25, 180}, {0.5, 0}, {1.25, 180},
	} {
		if got := propellerAngle(test.time); math.Abs(got-test.want) > 1e-9 {
			t.Fatalf("propellerAngle(%v)=%v, want %v", test.time, got, test.want)
		}
	}
}

func TestPropellerPivotMatchesSourceHub(t *testing.T) {
	x, y := propellerPivot(256)
	if math.Abs(x-98.6727272727) > 1e-9 || math.Abs(y-262.375) > 1e-9 {
		t.Fatalf("pivot=%v,%v, want 98.6727272727,262.375", x, y)
	}
}
