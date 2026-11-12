package main

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
)

const gateCapHeight = 22.0

var (
	colorBackground = ui.Hex("#F3FBFD")
	colorPanel      = ui.Hex("#E8F8FB")
	colorGo         = ui.Hex("#00ADD8")
	colorText       = ui.Hex("#00758D")
	colorMuted      = ui.Hex("#555759")
	colorDanger     = ui.Hex("#CE3262")
	colorSuccess    = ui.Hex("#00A29C")
	colorLightBlue  = ui.Hex("#5DC9E2")
	colorGround     = ui.Hex("#DBD9D6")
	colorOnAccent   = ui.Hex("#FFFFFF")
)

func withAlpha(color ui.Color, alpha float32) ui.Color {
	return ui.RGBA(color.R, color.G, color.B, alpha)
}

func (g *gameApp) drawLogical(r *logicalRenderer, alpha float64) {
	w := g.world
	g.drawBackground(r)
	groundTravel := w.previousGroundTravel + (w.groundTravel-w.previousGroundTravel)*alpha
	drawMountains(r, groundTravel, w.width)
	for _, current := range w.pipes {
		x := current.previousX + (current.x-current.previousX)*alpha
		drawTimeGate(r, x, current.gapY)
	}
	drawGround(r, groundTravel, w.width)
	birdY := w.previousBirdY + (w.birdY-w.previousBirdY)*alpha
	birdAngle := w.previousBirdAngle + (w.birdAngle-w.previousBirdAngle)*alpha
	g.drawGopher(r, birdY, birdAngle)
	if w.feedback != "" && w.feedbackTime > 0 {
		y := math.Max(100, birdY-28)
		r.CenterText(w.width, w.feedback, y, 10, withAlpha(colorSuccess, float32(math.Min(1, w.feedbackTime*2))))
	}
	r.CenterText(w.width, fmt.Sprintf("GC %d", w.score), 28, 30, colorText)
	r.CenterText(w.width, fmt.Sprintf("BEST TRACE %d", w.best), 65, 10, colorText)
	if g.muted {
		r.CenterText(w.width, "MUTED", 82, 9, colorText)
	}
	if w.state == ready {
		x, y, width, height := readyCardRect(w.width)
		r.Fill(x, y, width, height, withAlpha(colorOnAccent, 0.94), 8)
		r.CenterText(w.width, appTitle, 150, 30, colorText)
		r.CenterText(w.width, readyHint, 192, 11, colorText)
		r.CenterText(w.width, "P / ESC: pause    M: mute", 214, 8, colorMuted)
	}
	if title, hint := stateOverlay(w.state); title != "" {
		titleColor := colorText
		if w.state == gameOver {
			titleColor = colorDanger
		}
		drawOverlay(r, w.width, title, hint, titleColor)
	}
	if w.warning != "" {
		r.CenterText(w.width, w.warning, 503, 7, colorText)
	}
}

func (g *gameApp) drawBackground(r *logicalRenderer) {
	width := g.world.width
	r.Fill(0, 0, width, groundY, colorBackground, 0)
	r.Fill(0, 0, width, 3, colorGo, 0)
	keywords := [...]string{"go", "chan", "defer", "map", "interface{}", "func", "select", "range"}
	span := width + 120
	for i, keyword := range keywords {
		x := math.Mod(float64(i)*83-g.world.animationTime*12, span)
		if x < 0 {
			x += span
		}
		x -= 60
		y := 105 + float64(i%4)*72
		r.Text(keyword, x, y, 11, withAlpha(colorLightBlue, 0.55))
	}
}

func drawTimeGate(r *logicalRenderer, x, gapY float64) {
	maxMS, minMS := gateTimings(gapY)
	drawGateHalf(r, x, 0, gapY, true, maxMS)
	drawGateHalf(r, x, gapY+pipeGap, groundY-gapY-pipeGap, false, minMS)
	r.CenteredText("GC BUDGET", x+pipeWidth/2, gapY+pipeGap/2, 7, colorMuted)
}

func drawGateHalf(r *logicalRenderer, x, y, height float64, top bool, timingMS int) {
	if height <= 0 {
		return
	}
	capY := y
	if top {
		capY = y + height - gateCapHeight
	}
	r.Fill(x+5, y, pipeWidth-10, height, colorPanel, 0)
	r.Fill(x+5, y, 2, height, colorGo, 0)
	r.Fill(x+pipeWidth-7, y, 2, height, colorGo, 0)
	r.Fill(x-4, capY-1, pipeWidth+8, gateCapHeight+2, withAlpha(colorGo, 0.16), 4)
	r.Fill(x-2, capY+1, pipeWidth+4, gateCapHeight-2, colorText, 3)
	label := "STW MIN"
	if top {
		label = "STW MAX"
	}
	r.CenteredText(label, x+pipeWidth/2, capY+3, 7, colorOnAccent)
	r.CenteredText(fmt.Sprintf("%02dms", timingMS), x+pipeWidth/2, capY+13, 6, colorDanger)
}

func (g *gameApp) drawGopher(r *logicalRenderer, y, angle float64) {
	if g.world.state == playing {
		for i := range 3 {
			r.Fill(birdX-34-float64(i)*5, y-8+float64(i)*7, 6, 1, colorMuted, 0)
		}
	}
	r.Rotate(angle, birdX, y, func() {
		r.Image(g.assets.gopherBody, birdX-planeSpriteWidth/2, y-planeSpriteHeight/2, planeSpriteWidth, planeSpriteHeight, ui.FillBox)
		pivotX, pivotY := propellerPivot(y)
		r.Rotate(propellerAngle(g.world.animationTime), pivotX, pivotY, func() {
			r.Image(g.assets.gopherPropeller, birdX-planeSpriteWidth/2, y-planeSpriteHeight/2, planeSpriteWidth, planeSpriteHeight, ui.FillBox)
		})
	})
}

func readyCardRect(width float64) (x, y, cardWidth, cardHeight float64) {
	return (width - 320) / 2, 132, 320, 108
}

func drawOverlay(r *logicalRenderer, width float64, title, hint string, titleColor ui.Color) {
	r.Fill((width-248)/2, 175, 248, 100, withAlpha(colorOnAccent, 0.94), 8)
	r.CenterText(width, title, 196, 19, titleColor)
	r.CenterText(width, hint, 238, 9, colorText)
}
