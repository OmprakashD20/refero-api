package errors

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidPayload = errors.New("invalid data")
)

// Category
var (
	ErrCategoryNotFound       = errors.New("category not found")
	ErrCategoryExists         = errors.New("category already exists")
	ErrFailedToCreateCategory = errors.New("failed to create category")
	ErrFailedToUpdateCategory = errors.New("failed to update category")
	ErrFailedToDeleteCategory = errors.New("failed to delete category")
)

// Link
var (
	ErrLinkNotFound       = errors.New("link not found")
	ErrLinkExists         = errors.New("link already exists")
	ErrFailedToCreateLink = errors.New("failed to create link")
	ErrFailedToUpdateLink = errors.New("failed to update link")
	ErrFailedToDeleteLink = errors.New("failed to delete link")
)

// IsErrNoRows checks if the provided error is a pgx.ErrNoRows error.
func IsErrNoRows[T any](err error, value T) (T, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return value, nil
	}
	return value, err
}
