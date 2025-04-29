-- name: GetCampaigns :many
SELECT 
    id, user_id, title, description, wallet_address, image_path, target_amount, raised_amount, start_date, end_date, created_at
FROM 
    t_campaigns
WHERE
    (sqlc.narg(id)::UUID IS NULL OR id = sqlc.narg(id)::UUID) AND
    (sqlc.narg(user_id)::UUID IS NULL OR user_id = sqlc.narg(user_id)::UUID)
LIMIT @lim OFFSET @off;

-- name: GetCampaignByID :one
SELECT 
    id, user_id, title, description, wallet_address, image_path, target_amount, raised_amount, start_date, end_date, created_at
FROM 
    t_campaigns
WHERE 
    id = @campaign_id;

-- name: GetUserCampaign :one
SELECT 
    id, user_id, title, description, wallet_address, image_path, target_amount, raised_amount, start_date, end_date, created_at
FROM 
    t_campaigns
WHERE 
    id = @campaign_id AND user_id = @user_id;

-- name: CreateCampaign :one
INSERT INTO t_campaigns 
    (user_id, wallet_address, title, description, image_path, target_amount, raised_amount, start_date, end_date, created_at)
VALUES 
    (@user_id, @wallet_address, @title, @description, @image_path, @target_amount, DEFAULT, @start_date, @end_date, NOW())
RETURNING id;

-- name: UpdateCampaign :exec
UPDATE
    t_campaigns
SET
    wallet_address = COALESCE(@wallet_address, wallet_address),
    title = COALESCE(@title, title),
    description = COALESCE(@description, description),
    image_path = COALESCE(@image_path, image_path),
    target_amount = COALESCE(sqlc.narg(target_amount), target_amount),
    raised_amount = COALESCE(@raised_amount, raised_amount),
    start_date = COALESCE(@start_date, start_date),
    end_date = COALESCE(@end_date, end_date)
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
    (sqlc.narg(user_id)::UUID IS NULL OR user_id = sqlc.narg(user_id)::UUID);