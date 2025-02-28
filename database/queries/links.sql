-- Get all public links
-- name: GetAllLinks :many
SELECT id, url, title, description, short_url, created_at, updated_at FROM links;

-- Get link by ID
-- name: GetLinkByID :one
SELECT id, url, title, description, short_url, created_at, updated_at 
FROM links 
WHERE id = $1;

-- Get link by URL
-- name: GetLinkByURL :one
SELECT id, url, title, description, short_url FROM links WHERE url = $1;

-- Get link by short URL
-- name: GetLinkByShortURL :one
SELECT id, url, title, description, short_url FROM links 
WHERE short_url = $1;

-- Insert a new link
-- name: InsertLink :one
INSERT INTO links (url, title, description, short_url) 
VALUES ($1, $2, $3, $4) RETURNING id;

-- Update link details
-- name: UpdateLink :exec
UPDATE links 
SET title = $1, description = $2, updated_at = now() 
WHERE id = $3;

-- Delete link
-- name: DeleteLink :exec
DELETE FROM links WHERE id = $1;

-- Get paginated links
-- name: GetLinksPaginated :many
SELECT id, url, title, description, short_url, created_at, updated_at 
FROM links 
ORDER BY created_at DESC 
LIMIT $1 OFFSET $2;
