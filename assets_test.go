package main

import (
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestEmbeddedAssetsDecode(t *testing.T) {
	assets, err := loadGameAssets()
	if err != nil {
		t.Fatal(err)
	}
	for name, bitmap := range map[string]*ui.Bitmap{
		"body": assets.gopherBody, "propeller": assets.gopherPropeller,
	} {
		if bitmap == nil {
			t.Fatalf("missing %s bitmap", name)
		}
		if w, h := bitmap.Size(); w != 880 || h != 624 {
			t.Fatalf("%s size=%dx%d, want 880x624", name, w, h)
		}
	}
	for _, event := range []soundEvent{eventFlap, eventScore, eventHit} {
		if assets.effects[event] == nil {
			t.Fatalf("missing sound %v", event)
		}
	}
}

func TestBiplaneHasTransparentBackground(t *testing.T) {
	file, err := os.Open("resources/sprites/gopher.png")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	image, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	corner := color.NRGBAModel.Convert(image.At(0, 0)).(color.NRGBA)
	if corner.A != 0 {
		t.Fatalf("top-left alpha=%d, want 0", corner.A)
	}
}

func TestIconDimensions(t *testing.T) {
	file, err := os.Open("resources/icon.png")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	config, err := png.DecodeConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	if config.Width != 256 || config.Height != 256 {
		t.Fatalf("icon size=%dx%d, want 256x256", config.Width, config.Height)
	}
}
