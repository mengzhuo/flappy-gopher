package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func decodeSprite(t *testing.T, path string) image.Image {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoded, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	return decoded
}

func spriteRGBA(image image.Image, x, y int) color.NRGBA {
	return color.NRGBAModel.Convert(image.At(x, y)).(color.NRGBA)
}

func TestPropellerSplitDimensionsAndContent(t *testing.T) {
	for _, path := range []string{
		"resources/sprites/gopher-body.png",
		"resources/sprites/gopher-propeller.png",
	} {
		sprite := decodeSprite(t, path)
		if sprite.Bounds().Dx() != 880 || sprite.Bounds().Dy() != 624 {
			t.Fatalf("%s size=%dx%d", path, sprite.Bounds().Dx(), sprite.Bounds().Dy())
		}
		if spriteRGBA(sprite, 0, 0).A != 0 {
			t.Fatalf("%s top-left is not transparent", path)
		}
		nonTransparent := 0
		for y := 0; y < 624; y++ {
			for x := 0; x < 880; x++ {
				if spriteRGBA(sprite, x, y).A != 0 {
					nonTransparent++
				}
			}
		}
		if nonTransparent == 0 {
			t.Fatalf("%s is empty", path)
		}
	}
}

func TestPropellerSplitPreservesEverySourcePixel(t *testing.T) {
	source := decodeSprite(t, "resources/sprites/gopher.png")
	body := decodeSprite(t, "resources/sprites/gopher-body.png")
	propeller := decodeSprite(t, "resources/sprites/gopher-propeller.png")
	for y := 0; y < 624; y++ {
		for x := 0; x < 880; x++ {
			sourcePixel := spriteRGBA(source, x, y)
			bodyPixel := spriteRGBA(body, x, y)
			propellerPixel := spriteRGBA(propeller, x, y)
			if sourcePixel.A == 0 {
				if bodyPixel.A != 0 || propellerPixel.A != 0 {
					t.Fatalf("transparent source gained content at %d,%d", x, y)
				}
				continue
			}
			bodyOnly := bodyPixel == sourcePixel && propellerPixel.A == 0
			propellerOnly := propellerPixel == sourcePixel && bodyPixel.A == 0
			if !bodyOnly && !propellerOnly {
				t.Fatalf("source pixel not partitioned at %d,%d", x, y)
			}
		}
	}
}

func TestPropellerHubStaysWithBody(t *testing.T) {
	source := decodeSprite(t, "resources/sprites/gopher.png")
	body := decodeSprite(t, "resources/sprites/gopher-body.png")
	propeller := decodeSprite(t, "resources/sprites/gopher-propeller.png")
	hub := spriteRGBA(source, 766, 390)
	if hub.A == 0 || spriteRGBA(body, 766, 390) != hub || spriteRGBA(propeller, 766, 390).A != 0 {
		t.Fatal("hub must remain opaque in body and absent from propeller")
	}
}
