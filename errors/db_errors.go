package errors

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Category
var (
	ErrCategoryNotFound = errors.New("category not found")
)

func IsErrNoRows[T any](err error, value T) (T, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return value, nil
	}
	return value, err
}
