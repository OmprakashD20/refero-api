package links

import (
	"context"

	errs "github.com/OmprakashD20/refero-api/errors"
	"github.com/OmprakashD20/refero-api/repository"
	"github.com/OmprakashD20/refero-api/types"
	"github.com/OmprakashD20/refero-api/utils"
	validator "github.com/OmprakashD20/refero-api/validations"
)

type Store struct {
	db *repository.Queries
}

func NewStore(db *repository.Queries) *Store {
	return &Store{db}
}

func (s *Store) CheckIfLinkExistsByURL(ctx context.Context, url string) (*string, error) {
	link, err := s.db.CheckIfLinkExistsByURL(ctx, url)
	if err != nil {
		// Link doesn't exists in the database
		return errs.IsErrNoRows[*string](err, nil)
	}

	return utils.PgUUIDToStringPtr(link.ID), nil
}

func (s *Store) CreateLink(ctx context.Context, link validator.CreateLinkPayload, shortUrl string) (*string, error) {
	args := repository.CreateLinkParams{
		Title:       link.Title,
		Description: *link.Description,
		Url:         link.URL,
		ShortUrl:    shortUrl,
	}

	linkId, err := s.db.CreateLink(ctx, args)
	if !linkId.Valid {
		return nil, err
	}

	return utils.PgUUIDToStringPtr(linkId), nil
}

func (s *Store) GetCategoriesForLink(ctx context.Context, id string) ([]string, error) {
	categories, err := s.db.GetCategoriesForLink(ctx, utils.ToPgUUID(id))
	if err != nil {
		// Link doesn't exists in the database
		return errs.IsErrNoRows[[]string](err, nil)
	}

	categoryIds := make([]string, len(categories))
	for i, category := range categories {
		categoryIds[i] = category.ID.String()
	}
	return categoryIds, nil
}

func (s *Store) AddLinkToCategory(ctx context.Context, mappings []types.LinkCategoryDTO) error {
	var args []repository.AddLinkToCategoryParams

	for _, obj := range mappings {
		args = append(args, repository.AddLinkToCategoryParams{
			LinkID:     utils.ToPgUUID(obj.LinkID),
			CategoryID: utils.ToPgUUID(obj.CategoryID),
		})
	}

	_, err := s.db.AddLinkToCategory(ctx, args)
	return err
}

func (s *Store) GetLinkByShortURL(ctx context.Context, shortUrl string) (*types.LinkDTO, error) {
	link, err := s.db.GetLinkByShortURL(ctx, shortUrl)
	if err != nil {
		// Link doesn't exists in the database
		return errs.IsErrNoRows[*types.LinkDTO](err, nil)
	}

	data := &types.LinkDTO{
		ID:          *utils.PgUUIDToStringPtr(link.ID),
		Url:         link.Url,
		Title:       link.Title,
		Description: link.Description,
		ShortUrl:    link.ShortUrl,
	}

	return data, nil
}

func (s *Store) UpdateLinkByID(ctx context.Context, id string, link validator.UpdateLinkPayload) error {
	args := repository.UpdateLinkParams{
		ID:          utils.ToPgUUID(id),
		Title:       link.Title,
		Description: *link.Description,
	}

	rows, err := s.db.UpdateLink(ctx, args)
	if rows == 0 {
		// Link does not exists in the database
		return errs.ErrLinkNotFound
	}

	return err
}

func (s *Store) RemoveLinkToCategory(ctx context.Context, mappings []types.LinkCategoryDTO) error {
	var args []repository.RemoveLinkFromCategoryParams

	for _, obj := range mappings {
		args = append(args, repository.RemoveLinkFromCategoryParams{
			LinkID:     utils.ToPgUUID(obj.LinkID),
			CategoryID: utils.ToPgUUID(obj.CategoryID),
		})
	}

	// Execute batch deletion
	deleteBatch := s.db.RemoveLinkFromCategory(ctx, args)

	// Execute each deletion in the batch
	var batchErr error
	deleteBatch.Exec(func(i int, err error) {
		if err != nil {
			batchErr = errs.InternalServerError(errs.WithCause(err))
		}
	})

	// Close the batch
	if err := deleteBatch.Close(); err != nil {
		return err
	}

	return batchErr
}

func (s *Store) DeleteLinkByID(ctx context.Context, id string) error {
	rows, err := s.db.DeleteLink(ctx, utils.ToPgUUID(id))

	if rows == 0 {
		// Link does not exists in the database
		return errs.ErrLinkNotFound
	}

	return err
}
