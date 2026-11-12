package main

import (
	_ "embed"
	"fmt"

	"github.com/egoist/mygo/audio"
	"github.com/egoist/mygo/ui"
)

//go:embed resources/sprites/gopher-body.png
var gopherBodyPNG []byte

//go:embed resources/sprites/gopher-propeller.png
var gopherPropellerPNG []byte

//go:embed resources/audio/flap.wav
var flapWAV []byte

//go:embed resources/audio/score.wav
var scoreWAV []byte

//go:embed resources/audio/hit.wav
var hitWAV []byte

type gameAssets struct {
	gopherBody      *ui.Bitmap
	gopherPropeller *ui.Bitmap
	effects         map[soundEvent]func()
}

func loadGameAssets() (*gameAssets, error) {
	assets := &gameAssets{}
	body, err := ui.DecodeBitmap(gopherBodyPNG)
	if err != nil {
		return nil, fmt.Errorf("decode gopher body: %w", err)
	}
	propeller, err := ui.DecodeBitmap(gopherPropellerPNG)
	if err != nil {
		return nil, fmt.Errorf("decode gopher propeller: %w", err)
	}
	assets.gopherBody = body
	assets.gopherPropeller = propeller

	sounds := make(map[soundEvent]*audio.Sound)
	for event, data := range map[soundEvent][]byte{
		eventFlap:  flapWAV,
		eventScore: scoreWAV,
		eventHit:   hitWAV,
	} {
		sound, decodeErr := audio.DecodeSound(data)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode sound %d: %w", event, decodeErr)
		}
		sounds[event] = sound
	}
	assets.effects = make(map[soundEvent]func(), len(sounds))
	for event, sound := range sounds {
		assets.effects[event] = func() { sound.Play() }
	}
	return assets, nil
}
