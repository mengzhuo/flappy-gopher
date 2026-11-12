package main

const (
	appTitle  = "Flappy Gopher"
	readyHint = "SPACE: go func()"
)

func stateOverlay(current state) (title, hint string) {
	switch current {
	case paused:
		return "runtime.Gosched()", "P / ESC TO RESUME"
	case gameOver:
		return "STOP THE WORLD", "SPACE: restart goroutine"
	default:
		return "", ""
	}
}
