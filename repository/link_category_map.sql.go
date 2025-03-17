// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: link_category_map.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type AddLinkToCategoryParams struct {
	LinkID     pgtype.UUID `db:"link_id" json:"linkId"`
	CategoryID pgtype.UUID `db:"category_id" json:"categoryId"`
}

const getCategoriesForLink = `-- name: GetCategoriesForLink :many
SELECT c.id, c.name, c.description 
FROM category c 
JOIN link_category_map lcm ON c.id = lcm.category_id 
WHERE lcm.link_id = $1
`

type GetCategoriesForLinkRow struct {
	ID          pgtype.UUID `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description *string     `db:"description" json:"description"`
}

// Get all categories linked to a specific link
//
//  SELECT c.id, c.name, c.description
//  FROM category c
//  JOIN link_category_map lcm ON c.id = lcm.category_id
//  WHERE lcm.link_id = $1
func (q *Queries) GetCategoriesForLink(ctx context.Context, linkID pgtype.UUID) ([]GetCategoriesForLinkRow, error) {
	rows, err := q.db.Query(ctx, getCategoriesForLink, linkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCategoriesForLinkRow
	for rows.Next() {
		var i GetCategoriesForLinkRow
		if err := rows.Scan(&i.ID, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLinksForCategory = `-- name: GetLinksForCategory :many
SELECT l.id, l.url, l.title, l.description, l.short_url, l.created_at, l.updated_at 
FROM links l 
JOIN link_category_map lcm ON l.id = lcm.link_id 
WHERE lcm.category_id = $1
`

// Get all links in a category
//
//  SELECT l.id, l.url, l.title, l.description, l.short_url, l.created_at, l.updated_at
//  FROM links l
//  JOIN link_category_map lcm ON l.id = lcm.link_id
//  WHERE lcm.category_id = $1
func (q *Queries) GetLinksForCategory(ctx context.Context, categoryID pgtype.UUID) ([]Link, error) {
	rows, err := q.db.Query(ctx, getLinksForCategory, categoryID)
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

const getUncategorizedLinks = `-- name: GetUncategorizedLinks :many
SELECT id, url, title, description, short_url, created_at, updated_at 
FROM links l
WHERE NOT EXISTS (
    SELECT 1 
    FROM link_category_map lcm 
    WHERE l.id = lcm.link_id
)
`

// Get all uncategorized links
//
//  SELECT id, url, title, description, short_url, created_at, updated_at
//  FROM links l
//  WHERE NOT EXISTS (
//      SELECT 1
//      FROM link_category_map lcm
//      WHERE l.id = lcm.link_id
//  )
func (q *Queries) GetUncategorizedLinks(ctx context.Context) ([]Link, error) {
	rows, err := q.db.Query(ctx, getUncategorizedLinks)
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

const removeLinkFromCategory = `-- name: RemoveLinkFromCategory :execrows
DELETE FROM link_category_map 
WHERE link_id = $1 AND category_id = $2
`

type RemoveLinkFromCategoryParams struct {
	LinkID     pgtype.UUID `db:"link_id" json:"linkId"`
	CategoryID pgtype.UUID `db:"category_id" json:"categoryId"`
}

// Remove a link from a category
//
//  DELETE FROM link_category_map
//  WHERE link_id = $1 AND category_id = $2
func (q *Queries) RemoveLinkFromCategory(ctx context.Context, arg RemoveLinkFromCategoryParams) (int64, error) {
	result, err := q.db.Exec(ctx, removeLinkFromCategory, arg.LinkID, arg.CategoryID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
