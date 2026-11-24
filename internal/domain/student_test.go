package domain_test

import (
	"errors"
	"testing"

	"pos/internal/domain"
)

func TestNewStudent(t *testing.T) {
	student, err := domain.NewStudent("student-1", "Иван")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if student.ID() != "student-1" {
		t.Errorf("expected ID student-1, got %s", student.ID())
	}

	if student.Name() != "Иван" {
		t.Errorf("expected name Иван, got %s", student.Name())
	}
}

func TestNewStudentRejectsEmptyID(t *testing.T) {
	_, err := domain.NewStudent("", "Иван")

	if !errors.Is(err, domain.ErrEmptyID) {
		t.Fatalf("expected ErrEmptyID, got %v", err)
	}
}

func TestNewStudentRejectsEmptyName(t *testing.T) {
	_, err := domain.NewStudent("student-1", "")

	if !errors.Is(err, domain.ErrEmptyName) {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}
}
