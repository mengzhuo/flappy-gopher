package main

import "math"

const (
	defaultWorldWidth = 288.0
	logicalHeight     = 512.0
	groundY           = 440.0
	birdX             = 72.0
	birdRadius        = 9.0
	gravity           = 1050.0
	flapVelocity      = -330.0
	maxFallVelocity   = 480.0
	airDrag           = gravity / maxFallVelocity
	pipeSpeed         = 105.0
	pipeWidth         = 48.0
	pipeGap           = 128.0
	pipeInterval      = 1.6
	pipeSpacing       = pipeSpeed * pipeInterval
	fixedStep         = 1.0 / 60.0
)

const scoreFeedbackDuration = 1.0

type state uint8

const (
	ready state = iota
	playing
	paused
	gameOver
)

type controls struct {
	flap, pause, restart bool
	focused              bool
}

type pipe struct {
	x, previousX float64
	gapY         float64
	actionBudget int
	scored       bool
}

type randomSource interface{ IntN(int) int }

type replayRandom struct {
	source randomSource
	values []int
	cursor int
}

func (r *replayRandom) IntN(n int) int {
	if r.cursor < len(r.values) {
		value := r.values[r.cursor] % n
		r.cursor++
		return value
	}
	value := r.source.IntN(n)
	r.values = append(r.values, value)
	r.cursor++
	return value
}

func (r *replayRandom) rewind() {
	r.cursor = 0
}

type world struct {
	width                float64
	state                state
	birdY, previousBirdY float64
	velocity             float64
	birdAngle            float64
	previousBirdAngle    float64
	pipes                []pipe
	score, best          int
	feedback             string
	feedbackTime         float64
	bestDirty            bool
	groundTravel         float64
	previousGroundTravel float64
	warning              string
	events               []soundEvent
	animationTime        float64
	rng                  *replayRandom
	nextTransition       uint64
	lastGapCenter        float64
	units                decisionUnits
}

func (w *world) drainEvents() []soundEvent {
	events := w.events
	w.events = nil
	return events
}

func newWorld(best int, rng randomSource) *world {
	units, ok := calibrateDecisionUnits()
	if !ok {
		panic("invalid decision-unit physics")
	}
	w := &world{best: best, rng: &replayRandom{source: rng}, width: defaultWorldWidth, units: units}
	w.reset()
	return w
}

func (w *world) reset() {
	random := w.rng
	random.rewind()
	best := w.best
	width := w.width
	units := w.units
	if width <= 0 {
		width = defaultWorldWidth
	}
	*w = world{best: best, rng: random, width: width, units: units, state: ready}
	w.birdY, w.previousBirdY = logicalHeight/2, logicalHeight/2
	w.resetPipes()
}

func (w *world) pipeHorizon() float64 {
	return w.width + 80 + 3*pipeSpacing
}

func (w *world) resetPipes() {
	w.pipes = nil
	w.nextTransition = 0
	w.lastGapCenter = logicalHeight / 2
	for x := firstPipeX; x <= w.pipeHorizon()+1e-9; x += pipeSpacing {
		w.pipes = append(w.pipes, w.makePipe(x))
	}
}

func (w *world) ensurePipeCoverage() {
	if len(w.pipes) == 0 {
		x := w.width + 80
		if w.nextTransition == 0 {
			x = firstPipeX
		}
		w.pipes = append(w.pipes, w.makePipe(x))
	}
	for w.pipes[len(w.pipes)-1].x < w.pipeHorizon() {
		last := w.pipes[len(w.pipes)-1]
		w.pipes = append(w.pipes, w.makePipe(last.x+pipeSpacing))
	}
}

func (w *world) resize(width float64) {
	if width <= 0 || math.IsNaN(width) || math.IsInf(width, 0) || math.Abs(w.width-width) < 0.001 {
		return
	}
	w.width = width
	if w.state == ready {
		w.resetPipes()
		return
	}
	w.ensurePipeCoverage()
}

func (w *world) makePipe(x float64) pipe {
	center := logicalHeight / 2
	budget := 0
	if w.nextTransition > 0 {
		next, ok := generatePipeTransition(w.lastGapCenter, w.nextTransition, w.rng, w.units)
		if !ok {
			panic("invalid pipe transition")
		}
		center, budget = next.center, next.actionBudget
	}
	w.lastGapCenter = center
	w.nextTransition++
	return pipe{x: x, previousX: x, gapY: center - pipeGap/2, actionBudget: budget}
}
