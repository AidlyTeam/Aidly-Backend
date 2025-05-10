-- name: GetBadges :many
SELECT 
    id, symbol, name, description, seller_fee, icon_path, donation_threshold, uri, is_nft, created_at
FROM 
    t_badges
WHERE
    (sqlc.narg(id)::UUID IS NULL OR id = sqlc.narg(id)::UUID) AND
    (sqlc.narg(is_nft)::BOOLEAN IS NULL OR is_nft = sqlc.narg(is_nft)::BOOLEAN) AND
    (sqlc.narg(donation_threshold)::INTEGER IS NULL OR donation_threshold = sqlc.narg(donation_threshold)::INTEGER)
LIMIT @lim OFFSET @off;

-- name: GetBadgeByID :one
SELECT 
    id, symbol, name, description, seller_fee, icon_path, donation_threshold, uri, is_nft, created_at
FROM 
    t_badges
WHERE 
    id = @badge_id;

-- name: GetBadgeByDonationCount :one
SELECT 
    id, symbol, name, description, seller_fee, icon_path, donation_threshold, uri, is_nft, created_at
FROM 
    t_badges
WHERE 
    donation_threshold = @donation_threshold;

-- name: CreateBadge :one
INSERT INTO t_badges 
    (symbol, name, description, seller_fee, icon_path, uri, is_nft, donation_threshold)
VALUES 
    (@symbol, @name, @description, @seller_fee, @icon_path, @uri, @is_nft, @donation_threshold)
RETURNING id;

-- name: UpdateBadge :exec
UPDATE
    t_badges
SET
    symbol = COALESCE(sqlc.narg(symbol), symbol),
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    seller_fee = COALESCE(sqlc.narg(seller_fee), seller_fee),
    icon_path = COALESCE(sqlc.narg(icon_path), icon_path),
    uri = COALESCE(sqlc.narg(uri), uri),
    is_nft = COALESCE(sqlc.narg(is_nft), is_nft),
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