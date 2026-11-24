package memory

import (
	"pos/internal/domain"
)

type StudentRepository struct {
	students map[string]domain.Student
}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{
		students: make(map[string]domain.Student),
	}
}

func (r *StudentRepository) Save(student domain.Student) error {

	r.students[student.ID()] = student

	return nil
}

func (r *StudentRepository) FindByID(
	id string,
) (domain.Student, error) {

	student, exists := r.students[id]
	if !exists {
		return domain.Student{}, domain.ErrNotFound
	}

	return student, nil
}

func (r *StudentRepository) FindAll() ([]domain.Student, error) {

	result := make([]domain.Student, 0, len(r.students))

	for _, student := range r.students {
		result = append(result, student)
	}

	return result, nil
}

type CriterionRepository struct {
	criteria map[string]domain.Criterion
}

func NewCriterionRepository() *CriterionRepository {
	return &CriterionRepository{
		criteria: make(map[string]domain.Criterion),
	}
}

func (r *CriterionRepository) Save(
	criterion domain.Criterion,
) error {
	r.criteria[criterion.ID()] = criterion

	return nil
}

func (r *CriterionRepository) FindByID(
	id string,
) (domain.Criterion, error) {

	criterion, exists := r.criteria[id]
	if !exists {
		return domain.Criterion{}, domain.ErrNotFound
	}

	return criterion, nil
}

func (r *CriterionRepository) FindAll() ([]domain.Criterion, error) {
	result := make([]domain.Criterion, 0, len(r.criteria))

	for _, criterion := range r.criteria {
		result = append(result, criterion)
	}

	return result, nil
}

type EvaluationRepository struct {
	evaluations map[string]domain.Evaluation
}

func NewEvaluationRepository() *EvaluationRepository {
	return &EvaluationRepository{
		evaluations: make(map[string]domain.Evaluation),
	}
}

func evaluationKey(studentID, criterionID string) string {
	return studentID + "\x00" + criterionID
}

func (r *EvaluationRepository) Save(
	evaluation domain.Evaluation,
) error {

	key := evaluationKey(
		evaluation.StudentID(),
		evaluation.CriterionID(),
	)

	// Если оценка уже существует, она заменяется.
	r.evaluations[key] = evaluation

	return nil
}

func (r *EvaluationRepository) FindByStudent(
	studentID string,
) ([]domain.Evaluation, error) {

	result := make([]domain.Evaluation, 0)

	for _, evaluation := range r.evaluations {
		if evaluation.StudentID() == studentID {
			result = append(result, evaluation)
		}
	}

	return result, nil
}

func (r *EvaluationRepository) FindAll() (
	[]domain.Evaluation,
	error,
) {

	result := make(
		[]domain.Evaluation,
		0,
		len(r.evaluations),
	)

	for _, evaluation := range r.evaluations {
		result = append(result, evaluation)
	}

	return result, nil
}

func (r *EvaluationRepository) Delete(
	studentID, criterionID string,
) error {

	key := evaluationKey(studentID, criterionID)

	if _, exists := r.evaluations[key]; !exists {
		return domain.ErrNotFound
	}

	delete(r.evaluations, key)

	return nil
}
