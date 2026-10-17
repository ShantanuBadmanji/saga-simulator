-- name: GetCategoryByID :one
SELECT * FROM category WHERE id = ?1;

-- name: GetCategoryBySlug :one
SELECT * FROM category WHERE slug = ?1;

-- name: ListCategories :many
SELECT * FROM category order by slug;

-- name: ListSubCategories :many
SELECT * FROM category
where parent_id = ?1;

