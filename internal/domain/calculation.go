package domain

type CalculationStrategy interface {
	Calculate(
		evaluations []Evaluation,
		criteria map[string]Criterion,
	) (float64, error)
}

type WeightedAverageStrategy struct{}

func (s WeightedAverageStrategy) Calculate(
	evaluations []Evaluation,
	criteria map[string]Criterion,
) (float64, error) {
	if len(criteria) == 0 {
		return 0, ErrNoCriteria
	}

	// Отсутствующая оценка считается равной 0.
	scoresByCriterion := make(
		map[string]float64,
		len(evaluations),
	)

	for _, evaluation := range evaluations {
		criterionID := evaluation.CriterionID()

		if _, exists := criteria[criterionID]; !exists {
			return 0, ErrMissingCriterion
		}

		scoresByCriterion[criterionID] = evaluation.Score().Value()
	}

	var weightedSum float64
	var totalWeight float64

	for _, criterion := range criteria {
		score := scoresByCriterion[criterion.ID()]

		normalizedScore := score / criterion.MaxScore()

		weightedSum += normalizedScore * criterion.Weight()
		totalWeight += criterion.Weight()
	}

	if totalWeight <= 0 {
		return 0, ErrInvalidWeight
	}

	result := (weightedSum / totalWeight) * 10

	return result, nil
}

type TotalScoreService struct {
	strategy CalculationStrategy
}

func NewTotalScoreService(
	strategy CalculationStrategy,
) *TotalScoreService {
	return &TotalScoreService{
		strategy: strategy,
	}
}

func (s *TotalScoreService) Calculate(
	evaluations []Evaluation,
	criteria []Criterion,
) (float64, error) {
	if s.strategy == nil {
		return 0, ErrNilStrategy
	}

	criterionMap := make(map[string]Criterion, len(criteria))

	for _, criterion := range criteria {
		if _, exists := criterionMap[criterion.ID()]; exists {
			return 0, ErrDuplicateCriteria
		}

		criterionMap[criterion.ID()] = criterion
	}

	return s.strategy.Calculate(evaluations, criterionMap)
}
