// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: links.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkIfLinkExistsByURL = `-- name: CheckIfLinkExistsByURL :one
SELECT id, true AS exists FROM links l WHERE l.url = $1
UNION ALL
SELECT NULL, false AS exists WHERE NOT EXISTS (SELECT 1 FROM links WHERE links.url = $1)
LIMIT 1
`

type CheckIfLinkExistsByURLRow struct {
	ID     pgtype.UUID `db:"id" json:"id"`
	Exists bool        `db:"exists" json:"exists"`
}

// Check if link exists by URL
//
//  SELECT id, true AS exists FROM links l WHERE l.url = $1
//  UNION ALL
//  SELECT NULL, false AS exists WHERE NOT EXISTS (SELECT 1 FROM links WHERE links.url = $1)
//  LIMIT 1
func (q *Queries) CheckIfLinkExistsByURL(ctx context.Context, url string) (CheckIfLinkExistsByURLRow, error) {
	row := q.db.QueryRow(ctx, checkIfLinkExistsByURL, url)
	var i CheckIfLinkExistsByURLRow
	err := row.Scan(&i.ID, &i.Exists)
	return i, err
}

const deleteLink = `-- name: DeleteLink :execrows
DELETE FROM links WHERE id = $1
`

// Delete link
//
//  DELETE FROM links WHERE id = $1
func (q *Queries) DeleteLink(ctx context.Context, id pgtype.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteLink, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getAllLinks = `-- name: GetAllLinks :many
SELECT id, url, title, description, short_url, created_at, updated_at FROM links
`

// Get all public links
//
//  SELECT id, url, title, description, short_url, created_at, updated_at FROM links
func (q *Queries) GetAllLinks(ctx context.Context) ([]Link, error) {
	rows, err := q.db.Query(ctx, getAllLinks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Link
	for rows.Next() {
		var i Link
		if err := rows.Scan(
			&i.ID,
			&i.Url,
			&i.Title,
			&i.Description,
			&i.ShortUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLinkByID = `-- name: GetLinkByID :one
SELECT id, url, title, description, short_url, created_at, updated_at 
FROM links 
WHERE id = $1
`

// Get link by ID
//
//  SELECT id, url, title, description, short_url, created_at, updated_at
//  FROM links
//  WHERE id = $1
func (q *Queries) GetLinkByID(ctx context.Context, id pgtype.UUID) (Link, error) {
	row := q.db.QueryRow(ctx, getLinkByID, id)
	var i Link
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Title,
		&i.Description,
		&i.ShortUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLinkByShortURL = `-- name: GetLinkByShortURL :one
SELECT id, url, title, description, short_url FROM links 
WHERE short_url = $1
`

type GetLinkByShortURLRow struct {
	ID          pgtype.UUID `db:"id" json:"id"`
	Url         string      `db:"url" json:"url"`
	Title       string      `db:"title" json:"title"`
	Description string      `db:"description" json:"description"`
	ShortUrl    string      `db:"short_url" json:"shortUrl"`
}

// Get link by short URL
//
//  SELECT id, url, title, description, short_url FROM links
//  WHERE short_url = $1
func (q *Queries) GetLinkByShortURL(ctx context.Context, shortUrl string) (GetLinkByShortURLRow, error) {
	row := q.db.QueryRow(ctx, getLinkByShortURL, shortUrl)
	var i GetLinkByShortURLRow
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Title,
		&i.Description,
		&i.ShortUrl,
	)
	return i, err
}

const getLinkByURL = `-- name: GetLinkByURL :one
SELECT id, url, title, description, short_url FROM links WHERE url = $1
`

type GetLinkByURLRow struct {
	ID          pgtype.UUID `db:"id" json:"id"`
	Url         string      `db:"url" json:"url"`
	Title       string      `db:"title" json:"title"`
	Description string      `db:"description" json:"description"`
	ShortUrl    string      `db:"short_url" json:"shortUrl"`
}

// Get link by URL
//
//  SELECT id, url, title, description, short_url FROM links WHERE url = $1
func (q *Queries) GetLinkByURL(ctx context.Context, url string) (GetLinkByURLRow, error) {
	row := q.db.QueryRow(ctx, getLinkByURL, url)
	var i GetLinkByURLRow
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Title,
		&i.Description,
		&i.ShortUrl,
	)
	return i, err
}

const getLinksPaginated = `-- name: GetLinksPaginated :many
SELECT id, url, title, description, short_url, created_at, updated_at 
FROM links 
ORDER BY created_at DESC 
LIMIT $1 OFFSET $2
`

type GetLinksPaginatedParams struct {
	Limit  int32 `db:"limit" json:"limit"`
	Offset int32 `db:"offset" json:"offset"`
}

// Get paginated links
//
//  SELECT id, url, title, description, short_url, created_at, updated_at
//  FROM links
//  ORDER BY created_at DESC
//  LIMIT $1 OFFSET $2
func (q *Queries) GetLinksPaginated(ctx context.Context, arg GetLinksPaginatedParams) ([]Link, error) {
	rows, err := q.db.Query(ctx, getLinksPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Link
	for rows.Next() {
		var i Link
		if err := rows.Scan(
			&i.ID,
			&i.Url,
			&i.Title,
			&i.Description,
			&i.ShortUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertLink = `-- name: InsertLink :one
INSERT INTO links (url, title, description, short_url) 
VALUES ($1, $2, $3, $4) RETURNING id
`

type InsertLinkParams struct {
	Url         string `db:"url" json:"url"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	ShortUrl    string `db:"short_url" json:"shortUrl"`
}

// Insert a new link
//
//  INSERT INTO links (url, title, description, short_url)
//  VALUES ($1, $2, $3, $4) RETURNING id
func (q *Queries) InsertLink(ctx context.Context, arg InsertLinkParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, insertLink,
		arg.Url,
		arg.Title,
		arg.Description,
		arg.ShortUrl,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const updateLink = `-- name: UpdateLink :execrows
UPDATE links 
SET title = $1, description = $2, updated_at = now() 
WHERE id = $3
`

type UpdateLinkParams struct {
	Title       string      `db:"title" json:"title"`
	Description string      `db:"description" json:"description"`
	ID          pgtype.UUID `db:"id" json:"id"`
}

// Update link details
//
//  UPDATE links
//  SET title = $1, description = $2, updated_at = now()
//  WHERE id = $3
func (q *Queries) UpdateLink(ctx context.Context, arg UpdateLinkParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateLink, arg.Title, arg.Description, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
