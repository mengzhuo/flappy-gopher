package main

import "math"

const (
	planeSpriteWidth       = 72.0
	planeSpriteHeight      = 51.0
	planeSourceWidth       = 880.0
	planeSourceHeight      = 624.0
	propellerSourcePivotX  = 766.0
	propellerSourcePivotY  = 390.0
	propellerDegreesPerSec = 720.0
)

func propellerAngle(animationTime float64) float64 {
	return math.Mod(animationTime*propellerDegreesPerSec, 360)
}

func propellerPivot(birdY float64) (float64, float64) {
	x := birdX - planeSpriteWidth/2 + propellerSourcePivotX*planeSpriteWidth/planeSourceWidth
	y := birdY - planeSpriteHeight/2 + propellerSourcePivotY*planeSpriteHeight/planeSourceHeight
	return x, y
}
