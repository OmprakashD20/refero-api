package types

import (
	"context"

	validator "github.com/OmprakashD20/refero-api/validations"
)

type CategoryStore interface {
	CheckIfCategoryExists(ctx context.Context, name string) (bool, error)
	CreateCategory(ctx context.Context, category validator.CreateCategoryPayload) error
}
