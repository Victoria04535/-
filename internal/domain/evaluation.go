package domain

import (
	"math"
	"strings"
	"time"
)

type Score struct {
	value float64
	max   float64
}

func NewScore(value, max float64) (Score, error) {
	if max <= 0 || math.IsNaN(max) || math.IsInf(max, 0) {
		return Score{}, ErrInvalidMaxScore
	}

	if value < 0 || value > max || math.IsNaN(value) || math.IsInf(value, 0) {
		return Score{}, ErrInvalidScore
	}

	return Score{
		value: value,
		max:   max,
	}, nil
}

func (s Score) Value() float64 {
	return s.value
}

func (s Score) Max() float64 {
	return s.max
}

type Evaluation struct {
	studentID   string
	criterionID string
	score       Score
	comment     string
	createdAt   time.Time
}

type EvaluationFactory struct{}

func (f EvaluationFactory) Create(
	studentID string,
	criterion Criterion,
	value float64,
	comment string,
) (Evaluation, error) {
	studentID = strings.TrimSpace(studentID)

	if studentID == "" {
		return Evaluation{}, ErrInvalidStudentID
	}

	if criterion.ID() == "" {
		return Evaluation{}, ErrInvalidCriterion
	}

	score, err := NewScore(value, criterion.MaxScore())
	if err != nil {
		return Evaluation{}, err
	}

	return Evaluation{
		studentID:   studentID,
		criterionID: criterion.ID(),
		score:       score,
		comment:     strings.TrimSpace(comment),
		createdAt:   time.Now(),
	}, nil
}

func (e Evaluation) StudentID() string {
	return e.studentID
}

func (e Evaluation) CriterionID() string {
	return e.criterionID
}

func (e Evaluation) Score() Score {
	return e.score
}

func (e Evaluation) Comment() string {
	return e.comment
}

func (e Evaluation) CreatedAt() time.Time {
	return e.createdAt
}
