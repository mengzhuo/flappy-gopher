package main

import "testing"

func TestReadyCardIsCenteredBehindTitle(t *testing.T) {
	x, y, width, height := readyCardRect(910)
	if x != 295 || y != 132 || width != 320 || height != 108 {
		t.Fatalf("ready card=(%v, %v, %v, %v), want (295, 132, 320, 108)", x, y, width, height)
	}
}
