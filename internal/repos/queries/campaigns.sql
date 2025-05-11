-- name: GetCampaigns :many
SELECT 
    c.id, 
    c.user_id, 
    c.title, 
    c.description, 
    c.wallet_address, 
    c.image_path, 
    c.target_amount, 
    c.raised_amount, 
    c.accepted_token_symbol, 
    c.is_verified, 
    c.is_valid, 
    c.status_type, 
    c.start_date, 
    c.end_date, 
    c.created_at,
    u.name AS user_name,
    u.surname AS user_surname
FROM 
    t_campaigns c
JOIN 
    t_users u ON u.id = c.user_id
WHERE
    (sqlc.narg(id)::UUID IS NULL OR c.id = sqlc.narg(id)::UUID) AND
    (sqlc.narg(user_id)::UUID IS NULL OR c.user_id = sqlc.narg(user_id)::UUID) AND
    (sqlc.narg(is_verified)::BOOLEAN IS NULL OR c.is_verified = sqlc.narg(is_verified)::BOOLEAN) AND
    (sqlc.narg(status_type)::TEXT IS NULL OR c.status_type = sqlc.narg(status_type)::TEXT) AND
    (sqlc.narg(title)::TEXT IS NULL OR c.title ILIKE '%' || sqlc.narg(title)::TEXT || '%') AND
    c.is_valid = TRUE
LIMIT @lim OFFSET @off;


-- name: GetCampaignByID :one
SELECT 
    c.id, 
    c.user_id, 
    c.title, 
    c.description, 
    c.wallet_address, 
    c.image_path, 
    c.target_amount, 
    c.raised_amount, 
    c.accepted_token_symbol, 
    c.is_verified, 
    c.is_valid, 
    c.status_type, 
    c.start_date, 
    c.end_date, 
    c.created_at,
    u.name AS user_name,
    u.surname AS user_surname
FROM 
    t_campaigns c
JOIN 
    t_users u ON u.id = c.user_id
WHERE 
    c.id = @campaign_id;

-- name: GetUserCampaign :one
SELECT 
    id, user_id, title, description, wallet_address, image_path, target_amount, raised_amount, 
    accepted_token_symbol, is_verified, is_valid, status_type, start_date, end_date, created_at
FROM 
    t_campaigns
WHERE 
    id = @campaign_id AND user_id = @user_id;

-- name: CreateCampaign :one
INSERT INTO t_campaigns 
    (user_id, wallet_address, title, description, image_path, target_amount, raised_amount, status_type, accepted_token_symbol, start_date, end_date, created_at)
VALUES 
    (@user_id, @wallet_address, @title, @description, @image_path, @target_amount, DEFAULT, @status_type, @accepted_token_symbol, @start_date, @end_date, NOW())
RETURNING id;

-- name: UpdateCampaign :exec
UPDATE
    t_campaigns
SET
    wallet_address = COALESCE(sqlc.narg(wallet_address), wallet_address),
    title = COALESCE(sqlc.narg(title), title),
    description = COALESCE(sqlc.narg(description), description),
    image_path = COALESCE(sqlc.narg(image_path), image_path),
    target_amount = COALESCE(sqlc.narg(target_amount), target_amount),
    raised_amount = COALESCE(sqlc.narg(raised_amount), raised_amount),
    is_valid = COALESCE(sqlc.narg(is_valid), is_valid),
    status_type = COALESCE(sqlc.narg(status_type), status_type),
    accepted_token_symbol = COALESCE(sqlc.narg(accepted_token_symbol), accepted_token_symbol),
    start_date = COALESCE(sqlc.narg(start_date), start_date),
    end_date = COALESCE(sqlc.narg(end_date), end_date)
WHERE
    id = @campaign_id;

-- name: ChangeVerified :exec
UPDATE
    t_campaigns
SET
    is_verified = COALESCE(@is_verified, is_verified)
WHERE
    id = @campaign_id;

-- name: ChangeValid :exec
UPDATE
    t_campaigns
SET
    is_valid = COALESCE(@is_valid, is_valid)
WHERE
    id = @campaign_id;

-- name: DeleteCampaign :exec
DELETE FROM 
    t_campaigns
WHERE
    id = @campaign_id;

-- name: CountCampaigns :one
SELECT 
    COUNT(*) 
FROM 
    t_campaigns 
WHERE 
    (sqlc.narg(user_id)::UUID IS NULL OR user_id = sqlc.narg(user_id)::UUID)
    AND is_valid = FALSE;