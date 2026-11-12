package main

import (
	"github.com/egoist/mygo/plugins/game"
	"github.com/egoist/mygo/ui"
)

type gameApp struct {
	world   *world
	store   scoreStore
	assets  *gameAssets
	effects map[soundEvent]func()
	muted   bool
}

func newGameApp(w *world, store scoreStore, assets *gameAssets) *gameApp {
	return &gameApp{world: w, store: store, assets: assets, effects: assets.effects}
}

func (g *gameApp) view(c *ui.Context) {
	game.View(c, g).Fill().Background(colorBackground).Label(appTitle)
}

func (g *gameApp) Update(in *game.Input) {
	w := g.world
	if in.Pressed(ui.KeyM) {
		g.muted = !g.muted
	}
	w.step(controls{
		flap:    in.Pressed(ui.KeySpace) || in.Pressed(ui.KeyUp) || in.PointerPressed(0),
		pause:   in.Pressed(ui.KeyP) || in.Pressed(ui.KeyEscape),
		restart: in.Pressed(ui.KeyR),
		focused: in.Focused(),
	})
	g.consumeEvents()
	if w.bestDirty {
		w.bestDirty = false
		if g.store.path == "" {
			w.warning = "HIGH SCORE NOT SAVED"
		} else if err := g.store.save(w.best); err != nil {
			w.warning = "HIGH SCORE NOT SAVED"
		}
	}
}

func (g *gameApp) consumeEvents() {
	for _, event := range g.world.drainEvents() {
		if !g.muted && g.effects[event] != nil {
			g.effects[event]()
		}
	}
}

func (g *gameApp) syncViewport(r ui.Rect) bool {
	width, ok := logicalWidthForContent(r.W, r.H)
	if !ok {
		return false
	}
	g.world.resize(width)
	return true
}

func (g *gameApp) drawNative(p *ui.Painter, content ui.Rect, alpha float64) {
	if !g.syncViewport(content) {
		return
	}
	r, ok := newLogicalRenderer(p, content)
	if !ok {
		return
	}
	g.drawLogical(r, alpha)
}

func (g *gameApp) Draw(p *ui.Painter, r ui.Rect, alpha float64) {
	g.drawNative(p, r, alpha)
}
