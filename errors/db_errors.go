package errors

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Category
var (
	ErrCategoryNotFound       = errors.New("category not found")
	ErrCategoryExists         = errors.New("category already exists")
	ErrInvalidCategory        = errors.New("invalid category")
	ErrFailedToCreateCategory = errors.New("failed to create category")
	ErrFailedToUpdateCategory = errors.New("failed to update category")
	ErrFailedToDeleteCategory = errors.New("failed to delete category")
)

// IsErrNoRows checks if the provided error is a pgx.ErrNoRows error.
// If so, it returns the provided default value and a nil error.
//
// This function is useful for handling database queries where an absence
// of a row should not be treated as an error but rather as a valid case.
//
// Parameters:
//   - err (error): The error to check.
//   - value (T): The default value to return if the error is pgx.ErrNoRows.
//
// Returns:
//   - T: The default value if pgx.ErrNoRows, otherwise the original value.
//   - error: nil if pgx.ErrNoRows, otherwise the original error.
func IsErrNoRows[T any](err error, value T) (T, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return value, nil
	}
	return value, err
}
