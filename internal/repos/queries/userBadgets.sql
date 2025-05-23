-- name: GetUserBadges :many
SELECT 
    b.id, b.symbol, b.name, b.description, b.seller_fee, b.icon_path, b.donation_threshold, b.uri, b.is_nft, ub.is_minted, b.created_at
FROM 
    t_user_badges ub
JOIN 
    t_badges b ON ub.badge_id = b.id
WHERE 
    ub.user_id = @user_id;

-- name: GetUserBadge :one
SELECT 
    b.id, b.symbol, b.name, b.description, b.seller_fee, b.icon_path, b.donation_threshold, b.uri, b.is_nft, ub.is_minted, b.created_at
FROM 
    t_user_badges ub
JOIN 
    t_badges b ON ub.badge_id = b.id
WHERE 
    ub.user_id = @user_id AND
    b.id = @badge_id;

-- name: ChangeIsMinted :exec
UPDATE
    t_user_badges
SET
    is_minted = COALESCE(TRUE, is_minted)
WHERE
    user_id = @user_id AND badge_id = @badge_id;

-- name: AddUserBadge :one
INSERT INTO t_user_badges 
    (user_id, badge_id)
VALUES 
    (@user_id, @badge_id)
RETURNING id;

-- name: RemoveUserBadge :exec
DELETE FROM 
    t_user_badges
WHERE 
    user_id = @user_id AND badge_id = @badge_id;

-- name: GetUserBadgeExists :one
SELECT 
    EXISTS (
        SELECT 1 
        FROM t_user_badges 
        WHERE user_id = @user_id AND badge_id = @badge_id
    ) AS exists;

-- name: CountUserBadge :one
SELECT
    COUNT(*)
FROM
    t_user_badges
WHERE
    user_id = @user_id;