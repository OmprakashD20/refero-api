-- Get all categories
-- name: GetAllCategories :many
SELECT id, name, parent_id, description, created_at, updated_at FROM category;

-- Get category by ID
-- name: GetCategoryByID :one
SELECT id, name, parent_id, description FROM category 
WHERE id = $1;

-- Get category by name
-- name: GetCategoryByName :one
SELECT id, name, parent_id, description FROM category 
WHERE name = $1;

-- Get subcategories of a category
-- name: GetSubcategories :many
SELECT id, name, description, created_at, updated_at 
FROM category 
WHERE parent_id = $1;

-- Create a new category
-- name: CreateCategory :one
INSERT INTO category (name, parent_id, description) 
VALUES ($1, $2, $3) RETURNING id;

-- Update category details
-- name: UpdateCategory :execrows
UPDATE category 
SET name = $1, parent_id = $2, description = $3, updated_at = now() 
WHERE id = $4;

-- Delete category
-- name: DeleteCategory :exec
DELETE FROM category WHERE id = $1;
