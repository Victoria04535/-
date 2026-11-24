package application

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"pos/internal/domain"
)

type POSFacade struct {
	students    StudentRepository
	criteria    CriterionRepository
	evaluations EvaluationRepository
	calculator  *domain.TotalScoreService
	exporter    ResultExporter
	factory     domain.EvaluationFactory
}

func NewPOSFacade(
	students StudentRepository,
	criteria CriterionRepository,
	evaluations EvaluationRepository,
	calculator *domain.TotalScoreService,
	exporter ResultExporter,
) *POSFacade {
	return &POSFacade{
		students:    students,
		criteria:    criteria,
		evaluations: evaluations,
		calculator:  calculator,
		exporter:    exporter,
		factory:     domain.EvaluationFactory{},
	}
}

func (f *POSFacade) AddStudent(id, name string) error {
	student, err := domain.NewStudent(id, name)
	if err != nil {
		return err
	}

	if _, err := f.students.FindByID(student.ID()); err == nil {
		return fmt.Errorf("student %q already exists", student.ID())
	} else if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("find student: %w", err)
	}

	if err := f.students.Save(student); err != nil {
		return fmt.Errorf("save student: %w", err)
	}

	return nil
}

func (f *POSFacade) AddCriterion(
	id string,
	name string,
	maxScore float64,
	weight float64,
) error {
	criterion, err := domain.NewCriterion(
		id,
		name,
		maxScore,
		weight,
	)
	if err != nil {
		return err
	}

	if _, err := f.criteria.FindByID(criterion.ID()); err == nil {
		return fmt.Errorf(
			"criterion %q already exists",
			criterion.ID(),
		)
	} else if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("find criterion: %w", err)
	}

	if err := f.criteria.Save(criterion); err != nil {
		return fmt.Errorf("save criterion: %w", err)
	}

	return nil
}

func (f *POSFacade) ListStudents() ([]domain.Student, error) {
	return f.students.FindAll()
}

func (f *POSFacade) ListCriteria() ([]domain.Criterion, error) {
	return f.criteria.FindAll()
}

func (f *POSFacade) SetEvaluation(
	studentID string,
	criterionID string,
	score float64,
	comment string,
) error {
	studentID = strings.TrimSpace(studentID)
	criterionID = strings.TrimSpace(criterionID)

	if studentID == "" {
		return domain.ErrInvalidStudentID
	}

	if criterionID == "" {
		return domain.ErrInvalidCriterion
	}

	if _, err := f.students.FindByID(studentID); err != nil {
		return fmt.Errorf("find student: %w", err)
	}

	criterion, err := f.criteria.FindByID(criterionID)
	if err != nil {
		return fmt.Errorf("find criterion: %w", err)
	}

	evaluation, err := f.factory.Create(
		studentID,
		criterion,
		score,
		comment,
	)
	if err != nil {
		return err
	}

	if err := f.evaluations.Save(evaluation); err != nil {
		return fmt.Errorf("save evaluation: %w", err)
	}

	return nil
}

func (f *POSFacade) ListEvaluations(
	studentID string,
) ([]domain.Evaluation, error) {
	studentID = strings.TrimSpace(studentID)

	if studentID == "" {
		return nil, domain.ErrInvalidStudentID
	}

	if _, err := f.students.FindByID(studentID); err != nil {
		return nil, fmt.Errorf("find student: %w", err)
	}

	return f.evaluations.FindByStudent(studentID)
}

func (f *POSFacade) CalculateTotal(
	studentID string,
) (float64, error) {
	studentID = strings.TrimSpace(studentID)

	if studentID == "" {
		return 0, domain.ErrInvalidStudentID
	}

	if _, err := f.students.FindByID(studentID); err != nil {
		return 0, fmt.Errorf("find student: %w", err)
	}

	evaluations, err := f.evaluations.FindByStudent(studentID)
	if err != nil {
		return 0, fmt.Errorf("find evaluations: %w", err)
	}

	criteria, err := f.criteria.FindAll()
	if err != nil {
		return 0, fmt.Errorf("find criteria: %w", err)
	}

	return f.calculator.Calculate(evaluations, criteria)
}

func (f *POSFacade) BuildResults() ([]StudentResult, error) {
	students, err := f.students.FindAll()
	if err != nil {
		return nil, fmt.Errorf("find students: %w", err)
	}

	criteria, err := f.criteria.FindAll()
	if err != nil {
		return nil, fmt.Errorf("find criteria: %w", err)
	}

	allEvaluations, err := f.evaluations.FindAll()
	if err != nil {
		return nil, fmt.Errorf("find evaluations: %w", err)
	}

	evaluationsByStudent := make(
		map[string][]domain.Evaluation,
		len(students),
	)

	for _, evaluation := range allEvaluations {
		studentID := evaluation.StudentID()
		evaluationsByStudent[studentID] = append(
			evaluationsByStudent[studentID],
			evaluation,
		)
	}

	results := make([]StudentResult, 0, len(students))

	for _, student := range students {
		evaluations := evaluationsByStudent[student.ID()]

		total, err := f.calculator.Calculate(evaluations, criteria)
		if err != nil {
			return nil, fmt.Errorf(
				"calculate result for student %q: %w",
				student.ID(),
				err,
			)
		}

		results = append(results, StudentResult{
			Student: student,
			Total:   total,
		})
	}

	return results, nil
}

func (f *POSFacade) ExportResults(writer io.Writer) error {
	results, err := f.BuildResults()
	if err != nil {
		return fmt.Errorf("build results: %w", err)
	}

	if err := f.exporter.Export(results, writer); err != nil {
		return fmt.Errorf("export results: %w", err)
	}

	return nil
}

func (f *POSFacade) RemoveEvaluation(
	studentID string,
	criterionID string,
) error {
	studentID = strings.TrimSpace(studentID)
	criterionID = strings.TrimSpace(criterionID)

	if studentID == "" {
		return domain.ErrInvalidStudentID
	}

	if criterionID == "" {
		return domain.ErrInvalidCriterion
	}

	if _, err := f.students.FindByID(studentID); err != nil {
		return fmt.Errorf("find student: %w", err)
	}

	if _, err := f.criteria.FindByID(criterionID); err != nil {
		return fmt.Errorf("find criterion: %w", err)
	}

	if err := f.evaluations.Delete(studentID, criterionID); err != nil {
		return fmt.Errorf("delete evaluation: %w", err)
	}

	return nil
}
