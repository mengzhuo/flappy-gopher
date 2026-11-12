package main

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestLogicalRendererMapsGeometryAndTextAtNativeScale(t *testing.T) {
	r, ok := newLogicalRenderer(nil, ui.Rect{X: 10, Y: 20, W: 1152, H: 648})
	if !ok {
		t.Fatal("valid native content rejected")
	}
	got := r.rect(72, 51, 72, 51)
	want := ui.Rect{X: 101.125, Y: 84.546875, W: 91.125, H: 64.546875}
	for name, pair := range map[string][2]float32{
		"x": {got.X, want.X}, "y": {got.Y, want.Y},
		"width": {got.W, want.W}, "height": {got.H, want.H},
	} {
		if math.Abs(float64(pair[0]-pair[1])) > 1e-4 {
			t.Fatalf("%s=%v, want %v", name, pair[0], pair[1])
		}
	}
	if size := r.span("GC", 10, colorText).Size; math.Abs(float64(size-12.65625)) > 1e-4 {
		t.Fatalf("font size=%v, want 12.65625", size)
	}
	if math.Abs(float64(r.x(birdX)-101.125)) > 1e-4 || math.Abs(float64(r.y(birdX)-111.125)) > 1e-4 {
		t.Fatalf("rotation origin=%v,%v", r.x(birdX), r.y(birdX))
	}
}

func TestLogicalRendererRejectsInvalidContent(t *testing.T) {
	for _, content := range []ui.Rect{
		{}, {W: -1, H: 648}, {W: 1152, H: -1},
		{W: float32(math.NaN()), H: 648}, {W: 1152, H: float32(math.Inf(1))},
	} {
		if _, ok := newLogicalRenderer(nil, content); ok {
			t.Fatalf("invalid content accepted: %+v", content)
		}
	}
}
