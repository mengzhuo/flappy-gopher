package main

import "math"

func gateTimings(gapY float64) (maxMS, minMS int) {
	center := gapY + pipeGap/2
	t := (center - minPipeCenter) / (maxPipeCenter - minPipeCenter)
	return 6 + int(math.Round(8*t)), 1 + int(math.Round(3*t))
}
