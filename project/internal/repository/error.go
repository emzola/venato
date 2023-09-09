package repository

import "errors"

var (
	// ErrNotFound is returned when a requested record is not found.
	ErrNotFound = errors.New("the requested resource could not be found")
	// ErrEditConflict is returned when there ia an edit conflict due to a race condition.
	ErrEditConflict = errors.New("unable to update the record due to an edit conflict, please try again")
)
