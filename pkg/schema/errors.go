package schema

import "fmt"

type ErrFieldNotFound struct {
	Name string
}

func (e *ErrFieldNotFound) Error() string {
	return fmt.Sprintf("field %s not found", e.Name)
}

type ErrEmptyField struct {
	Name string
}

func (e *ErrEmptyField) Error() string {
	return fmt.Sprintf("field %s is empty", e.Name)
}
