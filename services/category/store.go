package category

import (
	"context"
	"errors"

	"github.com/OmprakashD20/refero-api/repository"
	"github.com/OmprakashD20/refero-api/types"
	"github.com/OmprakashD20/refero-api/utils"
	validator "github.com/OmprakashD20/refero-api/validations"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	db *repository.Queries
}

func NewStore(db *repository.Queries) *Store {
	return &Store{db}
}

func (s *Store) CheckIfCategoryExistsByName(ctx context.Context, name string) (bool, error) {
	_, err := s.db.GetCategoryByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Store) CheckIfCategoryExistsByID(ctx context.Context, id string) (bool, error) {
	_, err := s.db.GetCategoryByID(ctx, utils.ToPgUUID(id))
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
		ParentID:    utils.ToPgUUID(category.ParentId),
	}

	if categoryId, err := s.db.CreateCategory(ctx, args); !categoryId.Valid {
		return err
	}

	return nil
}

func (s *Store) GetAllCategories(ctx context.Context) ([]types.CategoryDTO, error) {
	data, err := s.db.GetAllCategories(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	categories := make([]types.CategoryDTO, len(data))
	for i, category := range data {
		categories[i] = types.CategoryDTO{
			ID:          category.ID.String(),
			Name:        category.Name,
			Description: category.Description,
			ParentID:    utils.PgUUIDToStringPtr(category.ParentID),
			CreatedAt:   &category.CreatedAt.Time,
			UpdatedAt:   &category.UpdatedAt.Time,
		}
	}

	return categories, nil
}

func (s *Store) GetCategoryByID(ctx context.Context, id string) (*types.CategoryDTO, error) {
	categoryId := utils.ToPgUUID(id)

	data, err := s.db.GetCategoryByID(ctx, categoryId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	category := &types.CategoryDTO{
		ID:          data.ID.String(),
		Name:        data.Name,
		Description: data.Description,
		ParentID:    utils.PgUUIDToStringPtr(data.ParentID),
	}

	return category, nil
}

func (s *Store) UpdateCategoryByID(ctx context.Context, id string, category validator.UpdateCategoryPayload) error {
	args := repository.UpdateCategoryParams{
		ID:          utils.ToPgUUID(id),
		Name:        category.Name,
		Description: category.Description,
		ParentID:    utils.ToPgUUID(category.ParentId),
	}

	rows, err := s.db.UpdateCategory(ctx, args)
	if rows == 0 {
		// [Category or Parent Category] does not exists in the database
		return errors.New("category does not exists or no changes applied")
	}

	return err
}

func (s *Store) DeleteCategoryByID(ctx context.Context, id string) error {
	rows, err := s.db.DeleteCategory(ctx, utils.ToPgUUID(id))

	if rows == 0 {
		// Category does not exists in the database
		return errors.New("category does not exists or no changes applied")
	}

	return err
}
