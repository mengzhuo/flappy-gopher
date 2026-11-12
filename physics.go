package main

import "math"

func verticalAcceleration(velocity float64) float64 {
	return gravity - airDrag*velocity
}

func targetBirdAngle(velocity float64) float64 {
	return math.Max(-12, math.Min(18, velocity/20))
}

func (w *world) updateBirdAngle() {
	const responsePerSecond = 12.0
	blend := 1 - math.Exp(-responsePerSecond*fixedStep)
	w.birdAngle += (targetBirdAngle(w.velocity) - w.birdAngle) * blend
	if math.IsNaN(w.birdAngle) || math.IsInf(w.birdAngle, 0) {
		w.birdAngle = 0
	}
}

func (w *world) step(in controls) {
	w.previousBirdY = w.birdY
	w.previousBirdAngle = w.birdAngle
	w.previousGroundTravel = w.groundTravel
	for i := range w.pipes {
		w.pipes[i].previousX = w.pipes[i].x
	}
	if in.restart {
		w.reset()
		return
	}
	flapHandled := false
	switch w.state {
	case ready:
		if in.flap && in.focused {
			w.state = playing
			w.velocity = flapVelocity
			w.events = append(w.events, eventFlap)
			flapHandled = true
		} else {
			return
		}
	case paused:
		if in.pause && in.focused {
			w.state = playing
		}
		return
	case gameOver:
		if in.flap && in.focused {
			w.reset()
			w.state = playing
			w.velocity = flapVelocity
			w.events = append(w.events, eventFlap)
		}
		return
	}
	if w.state == playing {
		if !in.focused || in.pause {
			w.state = paused
			return
		}
	}
	if in.flap && !flapHandled {
		w.velocity = flapVelocity
		w.events = append(w.events, eventFlap)
	}
	w.animationTime += fixedStep
	if w.feedbackTime > 0 {
		w.feedbackTime = math.Max(0, w.feedbackTime-fixedStep)
		if w.feedbackTime == 0 {
			w.feedback = ""
		}
	}
	w.velocity = math.Min(w.velocity+verticalAcceleration(w.velocity)*fixedStep, maxFallVelocity)
	w.updateBirdAngle()
	w.birdY += w.velocity * fixedStep
	if w.birdY < birdRadius {
		w.birdY = birdRadius
		w.velocity = math.Max(0, w.velocity)
	}
	w.groundTravel += pipeSpeed * fixedStep
	for i := range w.pipes {
		w.pipes[i].x -= pipeSpeed * fixedStep
	}
	for len(w.pipes) > 0 && w.pipes[0].x+pipeWidth < 0 {
		w.pipes = w.pipes[1:]
	}
	w.ensurePipeCoverage()
	for i := range w.pipes {
		p := &w.pipes[i]
		if !p.scored && p.x+pipeWidth < birdX {
			p.scored = true
			w.score++
			w.feedback = "GC OK +1"
			if w.score%5 == 0 {
				w.feedback = "128 MB freed"
			}
			w.feedbackTime = scoreFeedbackDuration
			w.events = append(w.events, eventScore)
		}
		if w.collidesWithPipe(*p) {
			w.finish()
			return
		}
	}
	if w.birdY+birdRadius >= groundY {
		w.birdY = groundY - birdRadius
		w.finish()
	}
}

func (w *world) collidesWithPipe(p pipe) bool {
	return circleHitsRect(birdX, w.birdY, birdRadius, p.x, 0, pipeWidth, p.gapY) ||
		circleHitsRect(birdX, w.birdY, birdRadius, p.x, p.gapY+pipeGap, pipeWidth, groundY-p.gapY-pipeGap)
}

func (w *world) finish() {
	if w.state == gameOver {
		return
	}
	w.feedback = ""
	w.feedbackTime = 0
	w.state = gameOver
	w.events = append(w.events, eventHit)
	if w.score > w.best {
		w.best = w.score
		w.bestDirty = true
	}
}

func circleHitsRect(cx, cy, radius, x, y, width, height float64) bool {
	nearestX := math.Max(x, math.Min(cx, x+width))
	nearestY := math.Max(y, math.Min(cy, y+height))
	dx, dy := cx-nearestX, cy-nearestY
	return dx*dx+dy*dy <= radius*radius
}
