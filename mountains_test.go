package main

import (
	"math"
	"reflect"
	"testing"
)

func TestRidgeYIsDeterministicBoundedAndAlternating(t *testing.T) {
	for _, layer := range mountainLayers {
		middle := (layer.minY + layer.maxY) / 2
		for index := int64(-8); index < 24; index++ {
			y := ridgeY(layer, index)
			if y != ridgeY(layer, index) || y < layer.minY || y > layer.maxY {
				t.Fatalf("layer seed=%x index=%d y=%v outside [%v,%v]", layer.seed, index, y, layer.minY, layer.maxY)
			}
			if index%2 == 0 && y > middle {
				t.Fatalf("even peak index=%d y=%v above middle %v", index, y, middle)
			}
			if index%2 != 0 && y < middle {
				t.Fatalf("odd valley index=%d y=%v below middle %v", index, y, middle)
			}
		}
	}
}

func TestVisibleRidgeParallaxAndResizePreserveIdentity(t *testing.T) {
	for _, layer := range mountainLayers {
		base := visibleRidge(layer, 80, 420)
		wide := visibleRidge(layer, 80, 620)
		if len(base) < 4 || len(wide) <= len(base) {
			t.Fatalf("unexpected ridge counts base=%d wide=%d", len(base), len(wide))
		}
		if !reflect.DeepEqual(base, wide[:len(base)]) {
			t.Fatal("resize changed existing ridge points")
		}
		shifted := visibleRidge(layer, 130, 420)
		byIndex := make(map[int64]ridgePoint, len(shifted))
		for _, point := range shifted {
			byIndex[point.index] = point
		}
		compared := 0
		for _, point := range base {
			moved, ok := byIndex[point.index]
			if !ok {
				continue
			}
			if moved.y != point.y || math.Abs((moved.x-point.x)+50*layer.parallax) > 1e-9 {
				t.Fatalf("index=%d moved from %+v to %+v", point.index, point, moved)
			}
			compared++
		}
		if compared == 0 {
			t.Fatal("no stable ridge identities found")
		}
	}
}

func TestVisibleRidgeRejectsInvalidInputAndBoundsWork(t *testing.T) {
	layer := mountainLayers[0]
	for _, input := range [][2]float64{
		{-1, 400}, {400, 0}, {400, -1},
		{math.NaN(), 400}, {400, math.NaN()},
		{math.Inf(1), 400}, {400, math.Inf(1)},
		{0, layer.segmentWidth * 5000},
	} {
		if points := visibleRidge(layer, input[0], input[1]); points != nil {
			t.Fatalf("invalid input accepted: %v, points=%d", input, len(points))
		}
	}
}

func TestMountainsDoNotConsumeGateRNG(t *testing.T) {
	withMountains := testWorld()
	withoutMountains := testWorld()
	for _, layer := range mountainLayers {
		_ = visibleRidge(layer, 12345, 1400)
		_ = ridgeY(layer, 999)
	}
	got := withMountains.makePipe(500)
	want := withoutMountains.makePipe(500)
	if got != want {
		t.Fatalf("mountain generation changed gate RNG: got %+v want %+v", got, want)
	}
}

func TestMountainColorsAreVisibleAndDistinct(t *testing.T) {
	seen := make(map[string]bool)
	for _, layer := range mountainLayers {
		if layer.color.A == 0 {
			t.Fatalf("transparent mountain layer seed=%x", layer.seed)
		}
		key := string([]byte{layer.color.R, layer.color.G, layer.color.B, layer.color.A})
		if seen[key] {
			t.Fatalf("duplicate mountain color at seed=%x", layer.seed)
		}
		seen[key] = true
	}
}
