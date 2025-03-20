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
	DeleteCategoryByID(ctx context.Context, id string) error
	GetLinksForCategory(ctx context.Context, id string) ([]LinkDTO, error)
}

type LinkStore interface {
	AddLinkToCategory(ctx context.Context, mappings []LinkCategoryDTO) error
	RemoveLinkToCategory(ctx context.Context, mappings []LinkCategoryDTO) error
	CheckIfLinkExistsByURL(ctx context.Context, url string) (*string, error)
	CreateLink(ctx context.Context, link validator.CreateLinkPayload, shortUrl string) (*string, error)
	GetLinkByShortURL(ctx context.Context, shortUrl string) (*LinkDTO, error)
	GetCategoriesForLink(ctx context.Context, id string) ([]string, error)
	UpdateLinkByID(ctx context.Context, id string, link validator.UpdateLinkPayload) error
	// DeleteLinkByID(ctx context.Context, id string) error
}

type CategoryDTO struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	ParentID    *string    `json:"parentId"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type LinkDTO struct {
	ID          string     `json:"id"`
	Url         string     `json:"url"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ShortUrl    string     `json:"shortUrl"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type LinkCategoryDTO struct {
	LinkID     string `json:"linkId"`
	CategoryID string `json:"categoryId"`
}
