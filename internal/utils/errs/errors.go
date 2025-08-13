package errs

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidID    = errors.New("invalid task ID")
	ErrValidation   = errors.New("validation error")
)
