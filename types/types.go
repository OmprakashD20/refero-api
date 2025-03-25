package types

import (
	"context"
	"time"

	"github.com/OmprakashD20/refero-api/repository"
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
	AddLinkToCategory(ctx context.Context, mappings []LinkCategoryDTO, txn *repository.Queries) error
	RemoveLinkFromCategory(ctx context.Context, mappings []LinkCategoryDTO, txn *repository.Queries) error
	CheckIfLinkExistsByURL(ctx context.Context, url string, txn *repository.Queries) (*string, error)
	CreateLink(ctx context.Context, link validator.CreateLinkPayload, shortUrl string, txn *repository.Queries) (*string, error)
	GetLinkByShortURL(ctx context.Context, shortUrl string, txn *repository.Queries) (*LinkDTO, error)
	GetCategoriesForLink(ctx context.Context, id string, txn *repository.Queries) ([]string, error)
	UpdateLinkByID(ctx context.Context, id string, link validator.UpdateLinkPayload, txn *repository.Queries) error
	DeleteLinkByID(ctx context.Context, id string, txn *repository.Queries) error
}

type TransactionStore interface {
	Exec(ctx context.Context, fn func(q *repository.Queries) error) error
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
