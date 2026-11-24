package domain_test

import (
	"errors"
	"math"
	"testing"

	"pos/internal/domain"
)

func TestNewScore(t *testing.T) {
	score, err := domain.NewScore(8, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if score.Value() != 8 {
		t.Errorf("expected value 8, got %.2f", score.Value())
	}

	if score.Max() != 10 {
		t.Errorf("expected max 10, got %.2f", score.Max())
	}
}

func TestNewScoreAllowsBoundaryValues(t *testing.T) {
	for _, value := range []float64{0, 10} {
		score, err := domain.NewScore(value, 10)
		if err != nil {
			t.Fatalf("value %.2f: unexpected error: %v", value, err)
		}
		if score.Value() != value {
			t.Errorf("expected value %.2f, got %.2f", value, score.Value())
		}
	}
}

func TestNewScoreRejectsNegativeValue(t *testing.T) {
	_, err := domain.NewScore(-1, 10)

	if !errors.Is(err, domain.ErrInvalidScore) {
		t.Fatalf("expected ErrInvalidScore, got %v", err)
	}
}

func TestNewScoreRejectsValueAboveMaximum(t *testing.T) {
	_, err := domain.NewScore(12, 10)

	if !errors.Is(err, domain.ErrInvalidScore) {
		t.Fatalf("expected ErrInvalidScore, got %v", err)
	}
}

func TestNewScoreRejectsNaNValue(t *testing.T) {
	_, err := domain.NewScore(math.NaN(), 10)

	if !errors.Is(err, domain.ErrInvalidScore) {
		t.Fatalf("expected ErrInvalidScore, got %v", err)
	}
}

func TestNewScoreRejectsInfValue(t *testing.T) {
	_, err := domain.NewScore(math.Inf(1), 10)

	if !errors.Is(err, domain.ErrInvalidScore) {
		t.Fatalf("expected ErrInvalidScore, got %v", err)
	}
}

func TestNewScoreRejectsNonPositiveMax(t *testing.T) {
	for _, max := range []float64{0, -1} {
		_, err := domain.NewScore(0, max)
		if !errors.Is(err, domain.ErrInvalidMaxScore) {
			t.Fatalf("max %.2f: expected ErrInvalidMaxScore, got %v", max, err)
		}
	}
}

func TestNewScoreRejectsNaNMax(t *testing.T) {
	_, err := domain.NewScore(0, math.NaN())

	if !errors.Is(err, domain.ErrInvalidMaxScore) {
		t.Fatalf("expected ErrInvalidMaxScore, got %v", err)
	}
}

func TestNewScoreRejectsInfMax(t *testing.T) {
	_, err := domain.NewScore(0, math.Inf(1))

	if !errors.Is(err, domain.ErrInvalidMaxScore) {
		t.Fatalf("expected ErrInvalidMaxScore, got %v", err)
	}
}

func mustCriterion(t *testing.T) domain.Criterion {
	t.Helper()

	criterion, err := domain.NewCriterion("criterion-1", "Практика", 10, 4)
	if err != nil {
		t.Fatalf("failed to create criterion: %v", err)
	}

	return criterion
}

func TestEvaluationFactoryCreate(t *testing.T) {
	criterion := mustCriterion(t)
	factory := domain.EvaluationFactory{}

	evaluation, err := factory.Create("student-1", criterion, 8, "  хорошо  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if evaluation.StudentID() != "student-1" {
		t.Errorf("expected studentID student-1, got %q", evaluation.StudentID())
	}

	if evaluation.CriterionID() != criterion.ID() {
		t.Errorf("expected criterionID %q, got %q", criterion.ID(), evaluation.CriterionID())
	}

	if evaluation.Score().Value() != 8 {
		t.Errorf("expected score value 8, got %.2f", evaluation.Score().Value())
	}

	if evaluation.Score().Max() != criterion.MaxScore() {
		t.Errorf("expected score max %.2f, got %.2f", criterion.MaxScore(), evaluation.Score().Max())
	}

	if evaluation.Comment() != "хорошо" {
		t.Errorf("expected trimmed comment %q, got %q", "хорошо", evaluation.Comment())
	}
}

func TestEvaluationFactoryRejectsEmptyStudentID(t *testing.T) {
	criterion := mustCriterion(t)
	factory := domain.EvaluationFactory{}

	for _, id := range []string{"", "   "} {
		_, err := factory.Create(id, criterion, 8, "")
		if !errors.Is(err, domain.ErrInvalidStudentID) {
			t.Fatalf("id %q: expected ErrInvalidStudentID, got %v", id, err)
		}
	}
}

func TestEvaluationFactoryRejectsEmptyCriterion(t *testing.T) {
	factory := domain.EvaluationFactory{}

	_, err := factory.Create("student-1", domain.Criterion{}, 8, "")

	if !errors.Is(err, domain.ErrInvalidCriterion) {
		t.Fatalf("expected ErrInvalidCriterion, got %v", err)
	}
}

func TestEvaluationFactoryRejectsInvalidScore(t *testing.T) {
	criterion := mustCriterion(t)
	factory := domain.EvaluationFactory{}

	_, err := factory.Create("student-1", criterion, 15, "")

	if !errors.Is(err, domain.ErrInvalidScore) {
		t.Fatalf("expected ErrInvalidScore, got %v", err)
	}
}

func TestEvaluationFactoryRejectsNegativeScore(t *testing.T) {
	criterion := mustCriterion(t)
	factory := domain.EvaluationFactory{}

	_, err := factory.Create("student-1", criterion, -1, "")

	if !errors.Is(err, domain.ErrInvalidScore) {
		t.Fatalf("expected ErrInvalidScore, got %v", err)
	}
}

func TestEvaluationFactoryRejectsNaNScore(t *testing.T) {
	criterion := mustCriterion(t)
	factory := domain.EvaluationFactory{}

	_, err := factory.Create("student-1", criterion, math.NaN(), "")

	if !errors.Is(err, domain.ErrInvalidScore) {
		t.Fatalf("expected ErrInvalidScore, got %v", err)
	}
}

func TestEvaluationFactoryRejectsInfScore(t *testing.T) {
	criterion := mustCriterion(t)
	factory := domain.EvaluationFactory{}

	_, err := factory.Create("student-1", criterion, math.Inf(1), "")

	if !errors.Is(err, domain.ErrInvalidScore) {
		t.Fatalf("expected ErrInvalidScore, got %v", err)
	}
}
