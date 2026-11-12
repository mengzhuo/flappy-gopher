package main

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestInitialWindowSizeUsesSixtyPercentOfWorkArea(t *testing.T) {
	width, height := initialWindowSize(1920, 1080)
	if width != 1152 || height != 648 {
		t.Fatalf("size=%dx%d, want 1152x648", width, height)
	}
	width, height = initialWindowSize(2560, 1440)
	if width != 1536 || height != 864 {
		t.Fatalf("size=%dx%d, want 1536x864", width, height)
	}
}

func TestInitialWindowSizeFallbackAndBounds(t *testing.T) {
	if width, height := initialWindowSize(0, 0); width != 1152 || height != 648 {
		t.Fatalf("fallback=%dx%d, want 1152x648", width, height)
	}
	if width, height := initialWindowSize(800, 500); width != 640 || height != 360 {
		t.Fatalf("minimum=%dx%d, want 640x360", width, height)
	}
	if width, height := initialWindowSize(500, 300); width != 500 || height != 300 {
		t.Fatalf("small work area=%dx%d, want 500x300", width, height)
	}
}

func TestLogicalWidthTracksContentAspectRatio(t *testing.T) {
	width, ok := logicalWidthForContent(1152, 648)
	want := 512.0 * 1152.0 / 648.0
	if !ok || math.Abs(width-want) > 1e-9 {
		t.Fatalf("width=%v ok=%v, want %v", width, ok, want)
	}
	for _, size := range [][2]float32{{0, 648}, {1152, 0}, {-1, 648}} {
		if _, ok := logicalWidthForContent(size[0], size[1]); ok {
			t.Fatalf("invalid content size accepted: %v", size)
		}
	}
	if _, ok := logicalWidthForContent(float32(math.NaN()), 648); ok {
		t.Fatal("NaN content width accepted")
	}
}

func TestSyncViewportUpdatesWorldWithoutIntermediateCanvas(t *testing.T) {
	w := testWorld()
	g := &gameApp{world: w}
	if !g.syncViewport(ui.Rect{W: 1152, H: 648}) {
		t.Fatal("valid viewport rejected")
	}
	wantWidth := 512.0 * 1152.0 / 648.0
	if math.Abs(w.width-wantWidth) > 1e-9 {
		t.Fatalf("world width=%v, want %v", w.width, wantWidth)
	}
	validWidth := w.width
	if g.syncViewport(ui.Rect{}) {
		t.Fatal("zero viewport accepted")
	}
	if w.width != validWidth {
		t.Fatalf("invalid viewport changed width to %v", w.width)
	}
}
