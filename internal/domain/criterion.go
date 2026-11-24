package domain

import (
	"math"
	"strings"
)

type Criterion struct {
	id       string
	name     string
	maxScore float64
	weight   float64
}

func NewCriterion(
	id string,
	name string,
	maxScore float64,
	weight float64,
) (Criterion, error) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)

	if id == "" {
		return Criterion{}, ErrEmptyID
	}

	if name == "" {
		return Criterion{}, ErrEmptyName
	}

	if maxScore <= 0 || math.IsNaN(maxScore) || math.IsInf(maxScore, 0) {
		return Criterion{}, ErrInvalidMaxScore
	}

	if weight <= 0 || math.IsNaN(weight) || math.IsInf(weight, 0) {
		return Criterion{}, ErrInvalidWeight
	}

	return Criterion{
		id:       id,
		name:     name,
		maxScore: maxScore,
		weight:   weight,
	}, nil
}

func (c Criterion) ID() string {
	return c.id
}

func (c Criterion) Name() string {
	return c.name
}

func (c Criterion) MaxScore() float64 {
	return c.maxScore
}

func (c Criterion) Weight() float64 {
	return c.weight
}
