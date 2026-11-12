package main

import "math"

const (
	firstPipeDelay = 4.0
	firstPipeX     = birdX + pipeSpeed*firstPipeDelay
	minPipeGapY    = 64.0
	maxPipeGapY    = groundY - pipeGap - 64.0
	minPipeCenter  = minPipeGapY + pipeGap/2
	maxPipeCenter  = maxPipeGapY + pipeGap/2
)

type decisionUnits struct {
	rise   float64
	fall   float64
	frames int
}

type pipeTransition struct {
	center       float64
	actionBudget int
}

func calibratedVelocity(velocity float64) float64 {
	return math.Min(velocity+verticalAcceleration(velocity)*fixedStep, maxFallVelocity)
}

func finitePositive(value float64) bool {
	return value > 0 && !math.IsNaN(value) && !math.IsInf(value, 0)
}

func calibrateDecisionUnits() (decisionUnits, bool) {
	velocity := flapVelocity
	rise := 0.0
	for frames := 1; frames <= 60; frames++ {
		velocity = calibratedVelocity(velocity)
		rise -= velocity * fixedStep
		if velocity >= 0 {
			fall, fallVelocity := 0.0, 0.0
			for range frames {
				fallVelocity = calibratedVelocity(fallVelocity)
				fall += fallVelocity * fixedStep
			}
			units := decisionUnits{rise: rise, fall: fall, frames: frames}
			return units, finitePositive(units.rise) && finitePositive(units.fall)
		}
	}
	return decisionUnits{}, false
}

func actionBudgetLimits(transition uint64) (int, int) {
	switch {
	case transition == 0:
		return 0, 0
	case transition <= 5:
		return 1, 1
	case transition <= 10:
		return 1, 2
	case transition <= 20:
		return 1, 3
	default:
		return 1, 4
	}
}

func generatePipeTransition(previous float64, transition uint64, rng randomSource, units decisionUnits) (pipeTransition, bool) {
	minBudget, maxBudget := actionBudgetLimits(transition)
	if minBudget == 0 || rng == nil || !finitePositive(units.rise) || !finitePositive(units.fall) || units.frames <= 0 ||
		previous < minPipeCenter || previous > maxPipeCenter || math.IsNaN(previous) || math.IsInf(previous, 0) {
		return pipeTransition{}, false
	}
	budget := minBudget
	if maxBudget > minBudget {
		budget += rng.IntN(maxBudget - minBudget + 1)
	}
	up := rng.IntN(2) == 0
	center := previous - float64(budget)*units.rise
	if !up {
		center = previous + float64(budget)*units.fall
	}
	if center < minPipeCenter || center > maxPipeCenter {
		opposite := previous - float64(budget)*units.rise
		if up {
			opposite = previous + float64(budget)*units.fall
		}
		if opposite >= minPipeCenter && opposite <= maxPipeCenter {
			center = opposite
		} else if center < minPipeCenter {
			center = 2*minPipeCenter - center
		} else {
			center = 2*maxPipeCenter - center
		}
	}
	result := pipeTransition{center: center, actionBudget: budget}
	if result.center < minPipeCenter || result.center > maxPipeCenter || math.IsNaN(result.center) || math.IsInf(result.center, 0) {
		return pipeTransition{}, false
	}
	return result, true
}
