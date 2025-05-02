-- name: GetBadges :many
SELECT 
    id, name, description, icon_path, created_at
FROM 
    t_badges
WHERE
    (sqlc.narg(id)::UUID IS NULL OR id = sqlc.narg(id)::UUID)
LIMIT @lim OFFSET @off;

-- name: GetBadgeByID :one
SELECT 
    id, name, description, icon_path, created_at
FROM 
    t_badges
WHERE 
    id = @badge_id;

-- name: CreateBadge :one
INSERT INTO t_badges 
    (name, description, icon_path)
VALUES 
    (@name, @description, @icon_path)
RETURNING id;

-- name: UpdateBadge :exec
UPDATE
    t_badges
SET
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    icon_path = COALESCE(sqlc.narg(icon_path), icon_path)
WHERE
    id = @badge_id;

-- name: DeleteBadge :exec
DELETE FROM 
    t_badges
WHERE
    id = @badge_id;

-- name: CountBadge :one
SELECT
    COUNT(*)
FROM
    t_badges
WHERE
    (sqlc.narg(id)::UUID IS NULL OR id = sqlc.narg(id)::UUID);