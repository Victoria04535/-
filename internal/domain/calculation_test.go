package domain_test

import (
	"errors"
	"math"
	"testing"

	"pos/internal/domain"
)

func mustCalcCriterion(t *testing.T, id, name string, maxScore, weight float64) domain.Criterion {
	t.Helper()

	criterion, err := domain.NewCriterion(id, name, maxScore, weight)
	if err != nil {
		t.Fatalf("failed to create criterion %q: %v", id, err)
	}

	return criterion
}

func mustEvaluation(t *testing.T, studentID string, criterion domain.Criterion, score float64) domain.Evaluation {
	t.Helper()

	factory := domain.EvaluationFactory{}

	evaluation, err := factory.Create(studentID, criterion, score, "")
	if err != nil {
		t.Fatalf("failed to create evaluation: %v", err)
	}

	return evaluation
}

func newService(t *testing.T) *domain.TotalScoreService {
	t.Helper()
	return domain.NewTotalScoreService(domain.WeightedAverageStrategy{})
}

func TestWeightedAverageCalculation(t *testing.T) {
	practice := mustCalcCriterion(t, "practice", "Практические работы", 10, 4)
	exam := mustCalcCriterion(t, "exam", "Контрольная", 20, 3)
	activity := mustCalcCriterion(t, "activity", "Активность", 5, 3)

	evaluations := []domain.Evaluation{
		mustEvaluation(t, "student-1", practice, 8),
		mustEvaluation(t, "student-1", exam, 15),
		mustEvaluation(t, "student-1", activity, 5),
	}

	service := newService(t)

	result, err := service.Calculate(
		evaluations,
		[]domain.Criterion{practice, exam, activity},
	)
	if err != nil {
		t.Fatalf("unexpected calculation error: %v", err)
	}

	const expected = 8.45
	const epsilon = 0.000001

	if math.Abs(result-expected) > epsilon {
		t.Errorf("expected %.2f, got %.2f", expected, result)
	}
}

func TestCalculationWithoutEvaluationsReturnsZero(t *testing.T) {
	criterion := mustCalcCriterion(t, "criterion-1", "Практика", 10, 5)

	service := newService(t)

	result, err := service.Calculate(nil, []domain.Criterion{criterion})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 0 {
		t.Errorf("expected 0, got %.2f", result)
	}
}

func TestCalculationTreatsMissingEvaluationAsZero(t *testing.T) {
	criterion1 := mustCalcCriterion(t, "c1", "Критерий 1", 10, 1)
	criterion2 := mustCalcCriterion(t, "c2", "Критерий 2", 10, 1)

	evaluations := []domain.Evaluation{
		mustEvaluation(t, "student-1", criterion1, 10),
	}

	service := newService(t)

	result, err := service.Calculate(
		evaluations,
		[]domain.Criterion{criterion1, criterion2},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const expected = 5.0
	const epsilon = 0.000001

	if math.Abs(result-expected) > epsilon {
		t.Errorf("expected %.2f, got %.2f", expected, result)
	}
}

func TestCalculationRejectsUnknownCriterion(t *testing.T) {
	registered := mustCalcCriterion(t, "registered", "Практика", 10, 5)
	unknown := mustCalcCriterion(t, "unknown", "Неизвестный критерий", 10, 5)

	evaluation := mustEvaluation(t, "student-1", unknown, 8)

	service := newService(t)

	_, err := service.Calculate(
		[]domain.Evaluation{evaluation},
		[]domain.Criterion{registered},
	)

	if !errors.Is(err, domain.ErrMissingCriterion) {
		t.Fatalf("expected ErrMissingCriterion, got %v", err)
	}
}

func TestCalculationRejectsEmptyCriteria(t *testing.T) {
	service := newService(t)

	_, err := service.Calculate(nil, nil)

	if !errors.Is(err, domain.ErrNoCriteria) {
		t.Fatalf("expected ErrNoCriteria, got %v", err)
	}
}

func TestCalculationRejectsDuplicateCriteria(t *testing.T) {
	criterion := mustCalcCriterion(t, "criterion-1", "Практика", 10, 5)

	service := newService(t)

	_, err := service.Calculate(
		nil,
		[]domain.Criterion{criterion, criterion},
	)

	if !errors.Is(err, domain.ErrDuplicateCriteria) {
		t.Fatalf("expected ErrDuplicateCriteria, got %v", err)
	}
}

func TestCalculationRejectsNilStrategy(t *testing.T) {
	service := domain.NewTotalScoreService(nil)

	_, err := service.Calculate(nil, nil)

	if !errors.Is(err, domain.ErrNilStrategy) {
		t.Fatalf("expected ErrNilStrategy, got %v", err)
	}
}
