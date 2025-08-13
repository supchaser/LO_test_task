package validate

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/supchaser/LO_test_task/internal/utils/errs"
)

const (
	MaxTaskTitleLength       = 200
	MinTaskTitleLength       = 3
	MaxTaskDescriptionLength = 5000
)

var taskTitleRegex = regexp.MustCompile(`^[A-Za-z0-9А-Яа-я\s.,!?-]+$`)

func CheckTaskTitle(title string) error {
	if title == "" {
		return fmt.Errorf("%w: task title cannot be empty", errs.ErrValidation)
	}

	length := utf8.RuneCountInString(title)
	if length < MinTaskTitleLength {
		return fmt.Errorf("%w: task title must be at least %d characters", errs.ErrValidation, MinTaskTitleLength)
	}

	if length > MaxTaskTitleLength {
		return fmt.Errorf("%w: task title cannot be longer than %d characters", errs.ErrValidation, MaxTaskTitleLength)
	}

	if !taskTitleRegex.MatchString(title) {
		return fmt.Errorf("%w: task title contains invalid characters", errs.ErrValidation)
	}

	return nil
}

func CheckTaskDescription(description string) error {
	length := utf8.RuneCountInString(description)
	if length > MaxTaskDescriptionLength {
		return fmt.Errorf("%w: task description cannot be longer than %d characters", errs.ErrValidation, MaxTaskDescriptionLength)
	}

	return nil
}
