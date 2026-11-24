package domain_test

import (
	"errors"
	"math"
	"testing"

	"pos/internal/domain"
)

func TestNewCriterion(t *testing.T) {
	_, err := domain.NewCriterion("criterion-1", "Практика", 10, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCriterionRejectsEmptyID(t *testing.T) {
	for _, id := range []string{"", "   "} {
		_, err := domain.NewCriterion(id, "Практика", 10, 4)
		if !errors.Is(err, domain.ErrEmptyID) {
			t.Fatalf("id %q: expected ErrEmptyID, got %v", id, err)
		}
	}
}

func TestNewCriterionRejectsEmptyName(t *testing.T) {
	for _, name := range []string{"", "   "} {
		_, err := domain.NewCriterion("criterion-1", name, 10, 4)
		if !errors.Is(err, domain.ErrEmptyName) {
			t.Fatalf("name %q: expected ErrEmptyName, got %v", name, err)
		}
	}
}

func TestNewCriterionRejectsNonPositiveMaxScore(t *testing.T) {
	for _, maxScore := range []float64{0, -1} {
		_, err := domain.NewCriterion("criterion-1", "Практика", maxScore, 4)
		if !errors.Is(err, domain.ErrInvalidMaxScore) {
			t.Fatalf("maxScore %.2f: expected ErrInvalidMaxScore, got %v", maxScore, err)
		}
	}
}

func TestNewCriterionRejectsNaNMaxScore(t *testing.T) {
	_, err := domain.NewCriterion("criterion-1", "Практика", math.NaN(), 4)

	if !errors.Is(err, domain.ErrInvalidMaxScore) {
		t.Fatalf("expected ErrInvalidMaxScore, got %v", err)
	}
}

func TestNewCriterionRejectsInfMaxScore(t *testing.T) {
	_, err := domain.NewCriterion("criterion-1", "Практика", math.Inf(1), 4)

	if !errors.Is(err, domain.ErrInvalidMaxScore) {
		t.Fatalf("expected ErrInvalidMaxScore, got %v", err)
	}
}

func TestNewCriterionRejectsNonPositiveWeight(t *testing.T) {
	for _, weight := range []float64{0, -1} {
		_, err := domain.NewCriterion("criterion-1", "Практика", 10, weight)
		if !errors.Is(err, domain.ErrInvalidWeight) {
			t.Fatalf("weight %.2f: expected ErrInvalidWeight, got %v", weight, err)
		}
	}
}

func TestNewCriterionRejectsNaNWeight(t *testing.T) {
	_, err := domain.NewCriterion("criterion-1", "Практика", 10, math.NaN())

	if !errors.Is(err, domain.ErrInvalidWeight) {
		t.Fatalf("expected ErrInvalidWeight, got %v", err)
	}
}

func TestNewCriterionRejectsInfWeight(t *testing.T) {
	_, err := domain.NewCriterion("criterion-1", "Практика", 10, math.Inf(1))

	if !errors.Is(err, domain.ErrInvalidWeight) {
		t.Fatalf("expected ErrInvalidWeight, got %v", err)
	}
}
