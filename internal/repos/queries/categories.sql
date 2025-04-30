-- name: GetCategories :many
SELECT
    id, name
FROM
    t_categories
ORDER BY
    name
LIMIT @lim OFFSET @off;

-- name: GetCategoryByID :one
SELECT
    id, name
FROM
    t_categories
WHERE
    id = @id;

-- name: CreateCategory :one
INSERT INTO t_categories (name)
VALUES (@name)
RETURNING id;

-- name: UpdateCategory :exec
UPDATE t_categories
SET name = COALESCE(sqlc.narg(name), name)
WHERE id = @id;

-- name: DeleteCategory :exec
DELETE FROM t_categories
WHERE id = @id;

-- name: CountCategory :one
SELECT 
    COUNT(*) 
FROM 
    t_categories;