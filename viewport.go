package main

import "math"

const (
	minWindowWidth       = 640
	minWindowHeight      = 360
	fallbackWindowWidth  = 1152
	fallbackWindowHeight = 648
)

func initialWindowSize(workWidth, workHeight int) (int, int) {
	if workWidth <= 0 || workHeight <= 0 {
		return fallbackWindowWidth, fallbackWindowHeight
	}
	width := int(math.Round(float64(workWidth) * 0.60))
	height := int(math.Round(float64(workHeight) * 0.60))
	width = max(minWindowWidth, width)
	height = max(minWindowHeight, height)
	return min(width, workWidth), min(height, workHeight)
}

func logicalWidthForContent(contentWidth, contentHeight float32) (float64, bool) {
	if contentWidth <= 0 || contentHeight <= 0 {
		return 0, false
	}
	width := logicalHeight * float64(contentWidth) / float64(contentHeight)
	if width <= 0 || math.IsNaN(width) || math.IsInf(width, 0) {
		return 0, false
	}
	return width, true
}
