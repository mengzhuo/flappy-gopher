package main

import (
	"log"
	"math/rand/v2"
	"path/filepath"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

const windowStateKey = "main-flappy-gopher-v1"

func main() {
	mygo.App.WhenReady(func() {
		best := 0
		warning := ""
		store := scoreStore{}
		dir, err := mygo.App.Path(mygo.PathUserData)
		if err != nil {
			warning = "HIGH SCORE UNAVAILABLE"
		} else {
			store.path = filepath.Join(dir, "high-score.json")
			best, err = store.load()
			if err != nil {
				best = 0
				warning = "HIGH SCORE RESET"
			}
		}
		assets, err := loadGameAssets()
		if err != nil {
			log.Fatalf("load game assets: %v", err)
		}
		w := newWorld(best, rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())))
		w.warning = warning
		g := newGameApp(w, store, assets)
		workArea := mygo.Screen.PrimaryDisplay().WorkArea
		windowWidth, windowHeight := initialWindowSize(workArea.Width, workArea.Height)
		mygo.NewWindow(mygo.WindowOptions{
			Title:           appTitle,
			Width:           windowWidth,
			Height:          windowHeight,
			UseContentSize:  true,
			MinWidth:        minWindowWidth,
			MinHeight:       minWindowHeight,
			StateKey:        windowStateKey,
			BackgroundColor: "#F3FBFD",
			Content:         ui.View(g.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
