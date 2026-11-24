package application

import (
	"io"

	"pos/internal/domain"
)

type StudentRepository interface {
	Save(student domain.Student) error
	FindByID(id string) (domain.Student, error)
	FindAll() ([]domain.Student, error)
}

type CriterionRepository interface {
	Save(criterion domain.Criterion) error
	FindByID(id string) (domain.Criterion, error)
	FindAll() ([]domain.Criterion, error)
}

type EvaluationRepository interface {
	Save(evaluation domain.Evaluation) error
	FindByStudent(studentID string) ([]domain.Evaluation, error)
	FindAll() ([]domain.Evaluation, error)
	Delete(studentID, criterionID string) error
}

type ResultExporter interface {
	Export(results []StudentResult, writer io.Writer) error
}
