package types

import (
	"context"
	"time"

	validator "github.com/OmprakashD20/refero-api/validations"
)

type CategoryStore interface {
	CheckIfCategoryExistsByName(ctx context.Context, name string) (bool, error)
	CheckIfCategoryExistsByID(ctx context.Context, id string) (bool, error)
	CreateCategory(ctx context.Context, category validator.CreateCategoryPayload) error
	GetAllCategories(ctx context.Context) ([]CategoryDTO, error)
	GetCategoryByID(ctx context.Context, id string) (*CategoryDTO, error)
	UpdateCategoryByID(ctx context.Context, id string, category validator.UpdateCategoryPayload) error
}

type CategoryDTO struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	ParentID    *string    `json:"parentId"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}
