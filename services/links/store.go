package links

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	errs "github.com/OmprakashD20/refero-api/errors"
	"github.com/OmprakashD20/refero-api/repository"
	"github.com/OmprakashD20/refero-api/types"
	"github.com/OmprakashD20/refero-api/utils"
	validator "github.com/OmprakashD20/refero-api/validations"
)

type Store struct {
	conn *pgxpool.Pool
	db   *repository.Queries
}

func NewStore(conn *pgxpool.Pool) *Store {
	return &Store{conn: conn, db: repository.New(conn)}
}

func (s *Store) CheckIfLinkExistsByURL(ctx context.Context, url string, txn *repository.Queries) (*string, error) {
	if txn == nil {
		txn = s.db
	}
	link, err := txn.CheckIfLinkExistsByURL(ctx, url)
	if err != nil {
		// Link doesn't exists in the database
		return errs.IsErrNoRows[*string](err, nil)
	}

	return utils.PgUUIDToStringPtr(link.ID), nil
}

func (s *Store) CreateLink(ctx context.Context, link validator.CreateLinkPayload, shortUrl string, txn *repository.Queries) (*string, error) {
	if txn == nil {
		txn = s.db
	}
	args := repository.CreateLinkParams{
		Title:       link.Title,
		Description: *link.Description,
		Url:         link.URL,
		ShortUrl:    shortUrl,
	}

	linkID, err := txn.CreateLink(ctx, args)
	if !linkID.Valid {
		return nil, err
	}

	return utils.PgUUIDToStringPtr(linkID), nil
}

func (s *Store) GetAllLinks(ctx context.Context) ([]types.LinkDTO, error) {
	data, err := s.db.GetAllLinks(ctx)
	if err != nil {
		return errs.IsErrNoRows[[]types.LinkDTO](err, nil)
	}

	links := make([]types.LinkDTO, len(data))
	for i, link := range data {
		links[i] = types.LinkDTO{
			ID:          link.ID.String(),
			Title:       link.Title,
			Description: link.Description,
			Url:         link.Url,
			ShortUrl:    link.ShortUrl,
			CreatedAt:   &link.CreatedAt.Time,
			UpdatedAt:   &link.UpdatedAt.Time,
		}
	}

	return links, nil
}

func (s *Store) GetLinkByID(ctx context.Context, id string) (*types.LinkDTO, error) {
	linkID := utils.ToPgUUID(id)
	
	data, err := s.db.GetLinkByID(ctx, linkID)
	if err != nil {
		return errs.IsErrNoRows[*types.LinkDTO](err, nil)
	}

	link := &types.LinkDTO{
		ID:          data.ID.String(),
		Title:       data.Title,
		Description: data.Description,
		Url:         data.Url,
		ShortUrl:    data.ShortUrl,
	}

	return link, nil
}

func (s *Store) GetCategoriesForLink(ctx context.Context, id string, txn *repository.Queries) ([]string, error) {
	if txn == nil {
		txn = s.db
	}
	categories, err := txn.GetCategoriesForLink(ctx, utils.ToPgUUID(id))
	if err != nil {
		// Link doesn't exists in the database
		return errs.IsErrNoRows[[]string](err, nil)
	}

	categoryIDs := make([]string, len(categories))
	for i, category := range categories {
		categoryIDs[i] = category.ID.String()
	}
	return categoryIDs, nil
}

func (s *Store) AddLinkToCategory(ctx context.Context, mappings []types.LinkCategoryDTO, txn *repository.Queries) error {
	if txn == nil {
		txn = s.db
	}
	var args []repository.AddLinkToCategoryParams

	for _, obj := range mappings {
		args = append(args, repository.AddLinkToCategoryParams{
			LinkID:     utils.ToPgUUID(obj.LinkID),
			CategoryID: utils.ToPgUUID(obj.CategoryID),
		})
	}

	_, err := txn.AddLinkToCategory(ctx, args)
	return err
}

func (s *Store) GetLinkByShortURL(ctx context.Context, shortUrl string, txn *repository.Queries) (*types.LinkDTO, error) {
	if txn == nil {
		txn = s.db
	}
	link, err := txn.GetLinkByShortURL(ctx, shortUrl)
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

func (s *Store) UpdateLinkByID(ctx context.Context, id string, link validator.UpdateLinkPayload, txn *repository.Queries) error {
	if txn == nil {
		txn = s.db
	}
	args := repository.UpdateLinkParams{
		ID:          utils.ToPgUUID(id),
		Title:       link.Title,
		Description: *link.Description,
	}

	rows, err := txn.UpdateLink(ctx, args)
	if rows == 0 {
		// Link does not exists in the database
		return errs.ErrLinkNotFound
	}

	return err
}

func (s *Store) RemoveLinkFromCategory(ctx context.Context, mappings []types.LinkCategoryDTO, txn *repository.Queries) error {
	if txn == nil {
		txn = s.db
	}
	var args []repository.RemoveLinkFromCategoryParams

	for _, obj := range mappings {
		args = append(args, repository.RemoveLinkFromCategoryParams{
			LinkID:     utils.ToPgUUID(obj.LinkID),
			CategoryID: utils.ToPgUUID(obj.CategoryID),
		})
	}

	// Execute batch deletion
	deleteBatch := txn.RemoveLinkFromCategory(ctx, args)

	// Execute each deletion in the batch
	var batchErr error
	deleteBatch.Exec(func(i int, err error) {
		if err != nil {
			batchErr = err
		}
	})

	// Close the batch
	if err := deleteBatch.Close(); err != nil {
		return err
	}

	return batchErr
}

func (s *Store) DeleteLinkByID(ctx context.Context, id string, txn *repository.Queries) error {
	if txn == nil {
		txn = s.db
	}
	rows, err := txn.DeleteLink(ctx, utils.ToPgUUID(id))

	if rows == 0 {
		// Link does not exists in the database
		return errs.ErrLinkNotFound
	}

	return err
}
