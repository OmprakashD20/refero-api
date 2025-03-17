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
func IsErrNoRows[T any](err error, value T) (T, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return value, nil
	}
	return value, err
}
