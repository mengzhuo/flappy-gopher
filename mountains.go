package main

import (
	"math"

	"github.com/egoist/mygo/ui"
)

const maxVisibleRidgePoints = 4096

type mountainLayer struct {
	seed         uint64
	segmentWidth float64
	minY         float64
	maxY         float64
	parallax     float64
	color        ui.Color
}

var mountainLayers = [...]mountainLayer{
	{seed: 0x476f4d6f756e7461, segmentWidth: 140, minY: 190, maxY: 290, parallax: 0.12, color: ui.Hex("#DDF4F8")},
	{seed: 0x476c696465727365, segmentWidth: 110, minY: 250, maxY: 345, parallax: 0.25, color: ui.Hex("#B6E8F0")},
	{seed: 0x4865617052696467, segmentWidth: 90, minY: 315, maxY: 410, parallax: 0.45, color: ui.Hex("#7ED6E6")},
}

type ridgePoint struct {
	index int64
	x     float64
	y     float64
}

func ridgeHash(value uint64) uint64 {
	value += 0x9e3779b97f4a7c15
	value = (value ^ (value >> 30)) * 0xbf58476d1ce4e5b9
	value = (value ^ (value >> 27)) * 0x94d049bb133111eb
	return value ^ (value >> 31)
}

func ridgeY(layer mountainLayer, index int64) float64 {
	hash := ridgeHash(layer.seed + uint64(index)*0x9e3779b97f4a7c15)
	unit := float64(hash>>11) * (1.0 / (1 << 53))
	middle := (layer.minY + layer.maxY) / 2
	low, high := layer.minY, middle
	if index%2 != 0 {
		low, high = middle, layer.maxY
	}
	return low + (high-low)*unit
}

func visibleRidge(layer mountainLayer, travel, width float64) []ridgePoint {
	if travel < 0 || width <= 0 || math.IsNaN(travel) || math.IsNaN(width) || math.IsInf(travel, 0) || math.IsInf(width, 0) {
		return nil
	}
	offset := travel * layer.parallax
	if offset < 0 || math.IsNaN(offset) || math.IsInf(offset, 0) || offset > 1<<62 || offset+width > 1<<62 {
		return nil
	}
	firstFloat := math.Floor(offset/layer.segmentWidth) - 1
	lastFloat := math.Floor((offset+width)/layer.segmentWidth) + 1
	if lastFloat-firstFloat+1 > maxVisibleRidgePoints {
		return nil
	}
	first, last := int64(firstFloat), int64(lastFloat)
	points := make([]ridgePoint, 0, int(last-first+1))
	for index := first; index <= last; index++ {
		points = append(points, ridgePoint{
			index: index,
			x:     float64(index)*layer.segmentWidth - offset,
			y:     ridgeY(layer, index),
		})
	}
	return points
}

func drawMountains(r *logicalRenderer, travel, width float64) {
	for _, layer := range mountainLayers {
		points := visibleRidge(layer, travel, width)
		if len(points) < 2 {
			continue
		}
		path := &ui.Path{}
		path.MoveTo(r.x(points[0].x), r.y(groundY))
		path.LineTo(r.x(points[0].x), r.y(points[0].y))
		for _, point := range points[1:] {
			path.LineTo(r.x(point.x), r.y(point.y))
		}
		path.LineTo(r.x(points[len(points)-1].x), r.y(groundY))
		path.Close()
		r.p.FillPath(path, layer.color)
	}
}
