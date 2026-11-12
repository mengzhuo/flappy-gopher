package main

import (
	"fmt"
	"math"
)

const (
	groundTileWidth   = 32.0
	heapBaseAddress   = uint64(0x00A0)
	heapAddressStride = uint64(0x20)
)

type heapBlock struct {
	x       float64
	address uint64
}

func heapAddress(blockIndex uint64) uint64 {
	return heapBaseAddress + blockIndex*heapAddressStride
}

func visibleHeapBlocks(travel, width float64) []heapBlock {
	if travel < 0 || width <= 0 || math.IsNaN(travel) || math.IsNaN(width) || math.IsInf(travel, 0) || math.IsInf(width, 0) {
		return nil
	}
	offset := math.Mod(travel, groundTileWidth)
	first := int64(math.Floor(travel / groundTileWidth))
	blocks := make([]heapBlock, 0, int(width/groundTileWidth)+2)
	for i := int64(0); ; i++ {
		x := float64(i)*groundTileWidth - offset
		if x >= width {
			break
		}
		blocks = append(blocks, heapBlock{x: x, address: heapAddress(uint64(first + i))})
	}
	return blocks
}

func drawGround(r *logicalRenderer, travel, width float64) {
	r.Fill(0, groundY, width, logicalHeight-groundY, colorGround, 0)
	r.Fill(0, groundY, width, 2, colorGo, 0)
	for _, block := range visibleHeapBlocks(travel, width) {
		r.Fill(block.x+1, groundY+6, 30, 25, colorPanel, 2)
		r.Text(fmt.Sprintf("0x%04X", block.address), block.x+3, groundY+14, 5, colorMuted)
	}
	r.CenteredText("HEAP", width/2, 480, 8, colorText)
}
