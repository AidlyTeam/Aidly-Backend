-- name: GetBadges :many
SELECT 
    id, name, description, icon_path, donation_threshold, created_at
FROM 
    t_badges
WHERE
    (sqlc.narg(id)::UUID IS NULL OR id = sqlc.narg(id)::UUID) AND
    (sqlc.narg(donation_threshold)::INTEGER IS NULL OR donation_threshold = sqlc.narg(donation_threshold)::INTEGER)
LIMIT @lim OFFSET @off;

-- name: GetBadgeByID :one
SELECT 
    id, name, description, icon_path, donation_threshold, created_at
FROM 
    t_badges
WHERE 
    id = @badge_id;

-- name: GetBadgeByDonationCount :one
SELECT 
    id, name, description, icon_path, donation_threshold, created_at
FROM 
    t_badges
WHERE 
    donation_threshold = @donation_threshold;

-- name: CreateBadge :one
INSERT INTO t_badges 
    (name, description, icon_path, donation_threshold)
VALUES 
    (@name, @description, @icon_path, @donation_threshold)
RETURNING id;

-- name: UpdateBadge :exec
UPDATE
    t_badges
SET
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    icon_path = COALESCE(sqlc.narg(icon_path), icon_path),
    donation_threshold = COALESCE(sqlc.narg(donation_threshold), donation_threshold)
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

-- name: ExistsBadgeByThreshold :one
SELECT 
    EXISTS (
        SELECT 1 
        FROM t_badges 
        WHERE donation_threshold = @donation_threshold 
    ) AS exists;