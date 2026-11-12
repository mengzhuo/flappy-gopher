package main

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func requireColor(t *testing.T, name string, got, want ui.Color) {
	t.Helper()
	if got.R != want.R || got.G != want.G || got.B != want.B || got.A != want.A {
		t.Fatalf("%s=%v, want %v", name, got, want)
	}
}

func TestGoLightPalette(t *testing.T) {
	for _, test := range []struct {
		name string
		got  ui.Color
		want ui.Color
	}{
		{"background", colorBackground, ui.Hex("#F3FBFD")},
		{"panel", colorPanel, ui.Hex("#E8F8FB")},
		{"go blue", colorGo, ui.Hex("#00ADD8")},
		{"deep aqua text", colorText, ui.Hex("#00758D")},
		{"gray", colorMuted, ui.Hex("#555759")},
		{"fuchsia", colorDanger, ui.Hex("#CE3262")},
		{"aqua", colorSuccess, ui.Hex("#00A29C")},
		{"light blue", colorLightBlue, ui.Hex("#5DC9E2")},
		{"ground", colorGround, ui.Hex("#DBD9D6")},
		{"on accent", colorOnAccent, ui.Hex("#FFFFFF")},
	} {
		requireColor(t, test.name, test.got, test.want)
	}
	wantMountains := []ui.Color{ui.Hex("#DDF4F8"), ui.Hex("#B6E8F0"), ui.Hex("#7ED6E6")}
	for i, want := range wantMountains {
		requireColor(t, "mountain", mountainLayers[i].color, want)
	}
}

func luminanceChannel(value uint8) float64 {
	channel := float64(value) / 255
	if channel <= 0.04045 {
		return channel / 12.92
	}
	return math.Pow((channel+0.055)/1.055, 2.4)
}

func relativeLuminance(color ui.Color) float64 {
	return 0.2126*luminanceChannel(color.R) + 0.7152*luminanceChannel(color.G) + 0.0722*luminanceChannel(color.B)
}

func contrastRatio(a, b ui.Color) float64 {
	la, lb := relativeLuminance(a), relativeLuminance(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

func TestPrimaryTextContrast(t *testing.T) {
	for _, test := range []struct {
		name       string
		foreground ui.Color
		background ui.Color
	}{
		{"main text", colorText, colorBackground},
		{"ground addresses", colorMuted, colorGround},
		{"gate cap labels", colorOnAccent, colorText},
		{"danger", colorDanger, colorBackground},
	} {
		if ratio := contrastRatio(test.foreground, test.background); ratio < 4.5 {
			t.Fatalf("%s contrast=%v, want at least 4.5", test.name, ratio)
		}
	}
}
