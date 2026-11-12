package main

import (
	"math"
	"math/rand/v2"
	"testing"
)

type fixedRNG int

func (r fixedRNG) IntN(n int) int { return int(r) % n }

func TestDecisionUnitsMatchProductionPhysics(t *testing.T) {
	units, ok := calibrateDecisionUnits()
	if !ok {
		t.Fatal("production physics did not produce valid decision units")
	}
	if units.frames != 15 {
		t.Fatalf("decision frames=%d, want 15", units.frames)
	}
	if math.Abs(units.rise-32.393216) > 1e-6 || math.Abs(units.fall-29.692909) > 1e-6 {
		t.Fatalf("decision units=%+v", units)
	}
}

func TestActionBudgetLimits(t *testing.T) {
	for _, test := range []struct {
		transition uint64
		min, max   int
	}{
		{0, 0, 0}, {1, 1, 1}, {5, 1, 1},
		{6, 1, 2}, {10, 1, 2}, {11, 1, 3},
		{20, 1, 3}, {21, 1, 4}, {1000, 1, 4},
	} {
		min, max := actionBudgetLimits(test.transition)
		if min != test.min || max != test.max {
			t.Fatalf("transition %d limits=%d..%d, want %d..%d", test.transition, min, max, test.min, test.max)
		}
	}
}

func TestPipeTransitionReversesAtBoundariesWithoutReducingBudget(t *testing.T) {
	units, _ := calibrateDecisionUnits()
	up, ok := generatePipeTransition(minPipeCenter, 1, fixedRNG(0), units)
	if !ok || up.actionBudget != 1 || math.Abs(up.center-(minPipeCenter+units.fall)) > 1e-9 {
		t.Fatalf("lower boundary transition=%+v ok=%v", up, ok)
	}
	down, ok := generatePipeTransition(maxPipeCenter, 1, fixedRNG(1), units)
	if !ok || down.actionBudget != 1 || math.Abs(down.center-(maxPipeCenter-units.rise)) > 1e-9 {
		t.Fatalf("upper boundary transition=%+v ok=%v", down, ok)
	}
}

func TestPipeTransitionReflectsWhenBothDirectionsExceedBounds(t *testing.T) {
	units, _ := calibrateDecisionUnits()
	previous := 224.0
	up, ok := generatePipeTransition(previous, 21, fixedRNG(2), units)
	wantUp := 2*minPipeCenter - (previous - 3*units.rise)
	if !ok || up.actionBudget != 3 || math.Abs(up.center-wantUp) > 1e-9 {
		t.Fatalf("double-boundary upward transition=%+v ok=%v want=%v", up, ok, wantUp)
	}
	previous = 245.17784326772997
	down, ok := generatePipeTransition(previous, 21, fixedRNG(3), units)
	wantDown := 2*maxPipeCenter - (previous + 4*units.fall)
	if !ok || down.actionBudget != 4 || math.Abs(down.center-wantDown) > 1e-9 {
		t.Fatalf("double-boundary downward transition=%+v ok=%v want=%v", down, ok, wantDown)
	}
}

func TestPipeTransitionRejectsInvalidInput(t *testing.T) {
	units, _ := calibrateDecisionUnits()
	for _, test := range []struct {
		previous   float64
		transition uint64
		rng        randomSource
		units      decisionUnits
	}{
		{256, 0, fixedRNG(0), units},
		{minPipeCenter - 1, 1, fixedRNG(0), units},
		{math.NaN(), 1, fixedRNG(0), units},
		{256, 1, nil, units},
		{256, 1, fixedRNG(0), decisionUnits{}},
	} {
		if got, ok := generatePipeTransition(test.previous, test.transition, test.rng, test.units); ok {
			t.Fatalf("invalid input accepted: %+v, result=%+v", test, got)
		}
	}
}

func TestPipeTransitionLongSequencesMatchBudgets(t *testing.T) {
	units, _ := calibrateDecisionUnits()
	seen := make(map[int]bool)
	for _, seed := range []uint64{1, 7, 42, 99} {
		rng := rand.New(rand.NewPCG(seed, seed+1))
		previous := logicalHeight / 2
		for transition := uint64(1); transition <= 160; transition++ {
			got, ok := generatePipeTransition(previous, transition, rng, units)
			if !ok {
				t.Fatalf("seed=%d transition=%d previous=%v units=%+v rejected", seed, transition, previous, units)
			}
			min, max := actionBudgetLimits(transition)
			if got.actionBudget < min || got.actionBudget > max {
				t.Fatalf("seed=%d transition=%d budget=%d outside %d..%d", seed, transition, got.actionBudget, min, max)
			}
			if got.center < minPipeCenter || got.center > maxPipeCenter || math.IsNaN(got.center) || math.IsInf(got.center, 0) {
				t.Fatalf("seed=%d transition=%d invalid center=%v", seed, transition, got.center)
			}
			upWant := previous - float64(got.actionBudget)*units.rise
			if upWant < minPipeCenter {
				upWant = 2*minPipeCenter - upWant
			}
			downWant := previous + float64(got.actionBudget)*units.fall
			if downWant > maxPipeCenter {
				downWant = 2*maxPipeCenter - downWant
			}
			if math.Abs(got.center-upWant) > 1e-9 && math.Abs(got.center-downWant) > 1e-9 {
				t.Fatalf("seed=%d transition=%d center=%v want %v or %v", seed, transition, got.center, upWant, downWant)
			}
			seen[got.actionBudget] = true
			previous = got.center
		}
	}
	for budget := 1; budget <= 4; budget++ {
		if !seen[budget] {
			t.Fatalf("budget %d never generated: %v", budget, seen)
		}
	}
}
