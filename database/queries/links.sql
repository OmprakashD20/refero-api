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

-- Check if link exists by URL
-- name: CheckIfLinkExistsByURL :one
SELECT id, true AS exists FROM links l WHERE l.url = $1
UNION ALL
SELECT NULL, false AS exists WHERE NOT EXISTS (SELECT 1 FROM links WHERE links.url = $1)
LIMIT 1;

-- Get link by short URL
-- name: GetLinkByShortURL :one
SELECT id, url, title, description, short_url FROM links 
WHERE short_url = $1;

-- Create a new link
-- name: CreateLink :one
INSERT INTO links (url, title, description, short_url) 
VALUES ($1, $2, $3, $4) RETURNING id;

-- Update link details
-- name: UpdateLink :execrows
UPDATE links 
SET title = $1, description = $2, updated_at = now() 
WHERE id = $3;

-- Delete link
-- name: DeleteLink :execrows
DELETE FROM links WHERE id = $1;

-- Get paginated links
-- name: GetLinksPaginated :many
SELECT id, url, title, description, short_url, created_at, updated_at 
FROM links 
ORDER BY created_at DESC 
LIMIT $1 OFFSET $2;
