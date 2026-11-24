package domain

import "strings"

type Student struct {
	id   string
	name string
}

func NewStudent(id, name string) (Student, error) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)

	if id == "" {
		return Student{}, ErrEmptyID
	}

	if name == "" {
		return Student{}, ErrEmptyName
	}

	return Student{
		id:   id,
		name: name,
	}, nil
}

func (s Student) ID() string {
	return s.id
}

func (s Student) Name() string {
	return s.name
}
