package main

import (
	"math"
	"testing"
)

func activeControls() controls { return controls{focused: true} }

func TestAirResistanceOpposesVelocityAndSetsTerminalSpeed(t *testing.T) {
	if got := verticalAcceleration(0); got != gravity {
		t.Fatalf("acceleration at rest=%v, want %v", got, gravity)
	}
	if got := verticalAcceleration(-100); got <= gravity {
		t.Fatalf("upward drag did not oppose motion: acceleration=%v", got)
	}
	if got := verticalAcceleration(100); got >= gravity {
		t.Fatalf("downward drag did not oppose motion: acceleration=%v", got)
	}
	if got := verticalAcceleration(maxFallVelocity); math.Abs(got) > 1e-9 {
		t.Fatalf("terminal acceleration=%v, want 0", got)
	}
	if got := verticalAcceleration(maxFallVelocity + 100); got >= 0 {
		t.Fatalf("above-terminal acceleration=%v, want negative", got)
	}
}

func TestFlapGravityAndCeiling(t *testing.T) {
	w := testWorld()
	in := activeControls()
	in.flap = true
	w.step(in)
	wantFirst := flapVelocity + verticalAcceleration(flapVelocity)*fixedStep
	if w.state != playing || math.Abs(w.velocity-wantFirst) > 1e-9 {
		t.Fatalf("first flap: state=%v velocity=%v, want %v", w.state, w.velocity, wantFirst)
	}
	w.step(activeControls())
	want := wantFirst + verticalAcceleration(wantFirst)*fixedStep
	if math.Abs(w.velocity-want) > 1e-9 {
		t.Fatalf("velocity=%v, want %v", w.velocity, want)
	}
	w.birdY, w.velocity = birdRadius, -100
	w.step(activeControls())
	if w.birdY < birdRadius || w.velocity < 0 {
		t.Fatalf("escaped ceiling: y=%v velocity=%v", w.birdY, w.velocity)
	}
}

func TestPauseFocusAndRestart(t *testing.T) {
	w := testWorld()
	w.step(controls{flap: true, focused: true})
	in := activeControls()
	in.pause = true
	w.step(in)
	y := w.birdY
	w.step(activeControls())
	if w.state != paused || w.birdY != y {
		t.Fatalf("pause advanced world: state=%v y=%v", w.state, w.birdY)
	}
	w.step(in)
	if w.state != playing {
		t.Fatalf("did not resume: %v", w.state)
	}
	w.step(controls{})
	if w.state != paused {
		t.Fatalf("focus loss did not pause: %v", w.state)
	}
	w.step(activeControls())
	if w.state != paused {
		t.Fatal("focus return auto-resumed")
	}
	w.step(controls{restart: true, focused: true})
	if w.state != ready || w.score != 0 {
		t.Fatalf("restart failed: %+v", w)
	}
}

func TestGroundCollisionEndsGame(t *testing.T) {
	w := testWorld()
	w.state, w.birdY = playing, groundY-birdRadius-0.1
	w.velocity = maxFallVelocity
	w.step(activeControls())
	if w.state != gameOver || w.birdY != groundY-birdRadius {
		t.Fatalf("ground collision: state=%v y=%v", w.state, w.birdY)
	}
}

func TestPipeScoreOnlyOnce(t *testing.T) {
	w := testWorld()
	w.state = playing
	w.birdY = 40
	w.velocity = -100
	w.pipes = []pipe{{x: birdX - pipeWidth + 1, previousX: birdX - pipeWidth + 1, gapY: 200}}
	w.step(activeControls())
	if w.score != 1 || !w.pipes[0].scored {
		t.Fatalf("first crossing: score=%d pipe=%+v", w.score, w.pipes[0])
	}
	w.step(activeControls())
	if w.score != 1 {
		t.Fatalf("pipe scored twice: %d", w.score)
	}
}

func TestCircleRectangleCollision(t *testing.T) {
	if !circleHitsRect(5, 5, 5, 10, 0, 20, 20) {
		t.Fatal("touching circle and rectangle must collide")
	}
	if circleHitsRect(0, 0, 4, 10, 10, 20, 20) {
		t.Fatal("distant circle collided")
	}
}

func TestPipeRecycling(t *testing.T) {
	w := testWorld()
	w.state, w.birdY, w.velocity = playing, 40, -100
	previousLastX := w.pipes[len(w.pipes)-1].x
	w.pipes[0].x = -pipeWidth - 1
	w.step(activeControls())
	if len(w.pipes) != 3 {
		t.Fatalf("pipe count=%d", len(w.pipes))
	}
	want := previousLastX - pipeSpeed*fixedStep + pipeSpacing
	if math.Abs(w.pipes[2].x-want) > 1e-9 {
		t.Fatalf("new pipe x=%v, want %v", w.pipes[2].x, want)
	}
}

func TestRandomSequenceReproducible(t *testing.T) {
	a, b := testWorld(), testWorld()
	for range 600 {
		a.step(activeControls())
		b.step(activeControls())
		for i := range a.pipes {
			if a.pipes[i].gapY != b.pipes[i].gapY {
				t.Fatalf("pipe %d differs: %v != %v", i, a.pipes[i].gapY, b.pipes[i].gapY)
			}
		}
	}
}
