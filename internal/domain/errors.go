package domain

import "errors"

var (
	ErrEmptyID           = errors.New("id cannot be empty")
	ErrEmptyName         = errors.New("name cannot be empty")
	ErrInvalidMaxScore   = errors.New("maximum score must be positive")
	ErrInvalidWeight     = errors.New("criterion weight must be positive")
	ErrInvalidScore      = errors.New("score must be between zero and maximum")
	ErrInvalidStudentID  = errors.New("student id cannot be empty")
	ErrInvalidCriterion  = errors.New("criterion id cannot be empty")
	ErrNotFound          = errors.New("record not found")
	ErrMissingCriterion  = errors.New("evaluation references unknown criterion")
	ErrNoCriteria        = errors.New("rubric has no criteria")
	ErrDuplicateCriteria = errors.New("duplicate criterion")
	ErrNilStrategy       = errors.New("calculation strategy is not configured")
)
