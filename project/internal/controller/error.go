package controller

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrNotFound is returned when a requested record is not found.
	ErrNotFound = errors.New("not found")
	// ErrFailedValidation is returned when there is a validation error.
	ErrFailedValidation = errors.New("failed validation")
	// ErrEditConflict is returned when there is an edit conflict error.
	ErrEditConflict = errors.New("edit conflict")
)

// FailedValidation loops through a validation error map and
// returns an error string with the key and value of the error map.
func FailedValidation(errorMap map[string]string) error {
	var s strings.Builder
	for k, v := range errorMap {
		s.WriteString(fmt.Sprintf("%s: %s; ", k, v))
	}
	err := fmt.Errorf("%s", s.String())
	return err
}
