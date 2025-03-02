package category

import (
	"context"
	"errors"

	"github.com/OmprakashD20/refero-api/repository"
	validator "github.com/OmprakashD20/refero-api/validations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Store struct {
	db *repository.Queries
}

func NewStore(db *repository.Queries) *Store {
	return &Store{db}
}

func (s *Store) CheckIfCategoryExists(ctx context.Context, name string) (bool, error) {
	_, err := s.db.GetCategoryByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Store) CreateCategory(ctx context.Context, category validator.CreateCategoryPayload) error {
	args := repository.CreateCategoryParams{
		Name:        category.Name,
		Description: category.Description,
		ParentID: func() pgtype.UUID {
			var uuid pgtype.UUID
			uuidStr := category.ParentId
			if err := uuid.Scan(uuidStr); err != nil {
				return pgtype.UUID{}
			}
			return uuid
		}(),
	}

	if categoryId, err := s.db.CreateCategory(ctx, args); !categoryId.Valid {
		return err
	}

	return nil
}
