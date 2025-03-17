-- Get all categories linked to a specific link
-- name: GetCategoriesForLink :many
SELECT c.id, c.name, c.description 
FROM category c 
JOIN link_category_map lcm ON c.id = lcm.category_id 
WHERE lcm.link_id = $1;

-- Get all links in a category
-- name: GetLinksForCategory :many
SELECT l.id, l.url, l.title, l.description, l.short_url, l.created_at, l.updated_at 
FROM links l 
JOIN link_category_map lcm ON l.id = lcm.link_id 
WHERE lcm.category_id = $1;

-- Get all uncategorized links
-- name: GetUncategorizedLinks :many
SELECT * 
FROM links l
WHERE NOT EXISTS (
    SELECT 1 
    FROM link_category_map lcm 
    WHERE l.id = lcm.link_id
);

-- Associate a link with a category
-- name: AddLinkToCategory :copyfrom
INSERT INTO link_category_map (link_id, category_id) VALUES ($1, $2);

-- Remove a link from a category
-- name: RemoveLinkFromCategory :execrows
DELETE FROM link_category_map 
WHERE link_id = $1 AND category_id = $2;
