// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: category.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO category (name, parent_id, description) 
VALUES ($1, $2, $3) RETURNING id
`

type CreateCategoryParams struct {
	Name        string      `db:"name" json:"name"`
	ParentID    pgtype.UUID `db:"parent_id" json:"parentId"`
	Description *string     `db:"description" json:"description"`
}

// Create a new category
//
//  INSERT INTO category (name, parent_id, description)
//  VALUES ($1, $2, $3) RETURNING id
func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, createCategory, arg.Name, arg.ParentID, arg.Description)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteCategory = `-- name: DeleteCategory :exec
DELETE FROM category WHERE id = $1
`

// Delete category
//
//  DELETE FROM category WHERE id = $1
func (q *Queries) DeleteCategory(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteCategory, id)
	return err
}

const getAllCategories = `-- name: GetAllCategories :many
SELECT id, name, parent_id, description, created_at, updated_at FROM category
`

// Get all categories
//
//  SELECT id, name, parent_id, description, created_at, updated_at FROM category
func (q *Queries) GetAllCategories(ctx context.Context) ([]Category, error) {
	rows, err := q.db.Query(ctx, getAllCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
			&i.Description,
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

const getCategoryByID = `-- name: GetCategoryByID :one
SELECT id, name, parent_id, description FROM category 
WHERE id = $1
`

type GetCategoryByIDRow struct {
	ID          pgtype.UUID `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	ParentID    pgtype.UUID `db:"parent_id" json:"parentId"`
	Description *string     `db:"description" json:"description"`
}

// Get category by ID
//
//  SELECT id, name, parent_id, description FROM category
//  WHERE id = $1
func (q *Queries) GetCategoryByID(ctx context.Context, id pgtype.UUID) (GetCategoryByIDRow, error) {
	row := q.db.QueryRow(ctx, getCategoryByID, id)
	var i GetCategoryByIDRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ParentID,
		&i.Description,
	)
	return i, err
}

const getCategoryByName = `-- name: GetCategoryByName :one
SELECT id, name, parent_id, description FROM category 
WHERE name = $1
`

type GetCategoryByNameRow struct {
	ID          pgtype.UUID `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	ParentID    pgtype.UUID `db:"parent_id" json:"parentId"`
	Description *string     `db:"description" json:"description"`
}

// Get category by name
//
//  SELECT id, name, parent_id, description FROM category
//  WHERE name = $1
func (q *Queries) GetCategoryByName(ctx context.Context, name string) (GetCategoryByNameRow, error) {
	row := q.db.QueryRow(ctx, getCategoryByName, name)
	var i GetCategoryByNameRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ParentID,
		&i.Description,
	)
	return i, err
}

const getSubcategories = `-- name: GetSubcategories :many
SELECT id, name, description, created_at, updated_at 
FROM category 
WHERE parent_id = $1
`

type GetSubcategoriesRow struct {
	ID          pgtype.UUID      `db:"id" json:"id"`
	Name        string           `db:"name" json:"name"`
	Description *string          `db:"description" json:"description"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"createdAt"`
	UpdatedAt   pgtype.Timestamp `db:"updated_at" json:"updatedAt"`
}

// Get subcategories of a category
//
//  SELECT id, name, description, created_at, updated_at
//  FROM category
//  WHERE parent_id = $1
func (q *Queries) GetSubcategories(ctx context.Context, parentID pgtype.UUID) ([]GetSubcategoriesRow, error) {
	rows, err := q.db.Query(ctx, getSubcategories, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubcategoriesRow
	for rows.Next() {
		var i GetSubcategoriesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
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

const updateCategory = `-- name: UpdateCategory :execrows
UPDATE category 
SET name = $1, parent_id = $2, description = $3, updated_at = now() 
WHERE id = $4
`

type UpdateCategoryParams struct {
	Name        string      `db:"name" json:"name"`
	ParentID    pgtype.UUID `db:"parent_id" json:"parentId"`
	Description *string     `db:"description" json:"description"`
	ID          pgtype.UUID `db:"id" json:"id"`
}

// Update category details
//
//  UPDATE category
//  SET name = $1, parent_id = $2, description = $3, updated_at = now()
//  WHERE id = $4
func (q *Queries) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateCategory,
		arg.Name,
		arg.ParentID,
		arg.Description,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
