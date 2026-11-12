package main

import (
	"math"

	"github.com/egoist/mygo/ui"
)

type logicalRenderer struct {
	p       *ui.Painter
	originX float32
	originY float32
	scale   float32
}

func newLogicalRenderer(p *ui.Painter, content ui.Rect) (*logicalRenderer, bool) {
	values := [...]float32{content.X, content.Y, content.W, content.H}
	for _, value := range values {
		if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
			return nil, false
		}
	}
	if content.W <= 0 || content.H <= 0 {
		return nil, false
	}
	scale := content.H / float32(logicalHeight)
	if scale <= 0 || math.IsNaN(float64(scale)) || math.IsInf(float64(scale), 0) {
		return nil, false
	}
	return &logicalRenderer{p: p, originX: content.X, originY: content.Y, scale: scale}, true
}

func (r *logicalRenderer) x(x float64) float32 {
	return r.originX + float32(x)*r.scale
}

func (r *logicalRenderer) y(y float64) float32 {
	return r.originY + float32(y)*r.scale
}

func (r *logicalRenderer) rect(x, y, width, height float64) ui.Rect {
	return ui.Rect{X: r.x(x), Y: r.y(y), W: float32(width) * r.scale, H: float32(height) * r.scale}
}

func (r *logicalRenderer) Fill(x, y, width, height float64, color ui.Color, radius float64) {
	r.p.Fill(r.rect(x, y, width, height), color, float32(radius)*r.scale)
}

func (r *logicalRenderer) span(text string, size float64, color ui.Color) ui.Span {
	return ui.Span{Text: text, Size: float32(size) * r.scale, Weight: 700, Color: color}
}

func (r *logicalRenderer) Text(text string, x, y, size float64, color ui.Color) {
	r.p.RichText(r.x(x), r.y(y), 0, r.span(text, size, color))
}

func (r *logicalRenderer) CenteredText(text string, centerX, y, size float64, color ui.Color) {
	label := r.span(text, size, color)
	width, _ := r.p.MeasureText(0, label)
	r.p.RichText(r.x(centerX)-width/2, r.y(y), 0, label)
}

func (r *logicalRenderer) CenterText(width float64, text string, y, size float64, color ui.Color) {
	r.CenteredText(text, width/2, y, size, color)
}

func (r *logicalRenderer) Image(source ui.ImageSource, x, y, width, height float64, fit ui.Fit) {
	r.p.Image(source, r.rect(x, y, width, height), fit)
}

func (r *logicalRenderer) Rotate(angle, originX, originY float64, fn func()) {
	r.p.Transform(ui.Rotation(float32(angle)), r.x(originX), r.y(originY), fn)
}
