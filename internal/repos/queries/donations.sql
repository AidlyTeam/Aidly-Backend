-- name: GetDonations :many
SELECT
    id, campaign_id, transaction_id, user_id, amount, donation_date 
FROM
    t_donations
WHERE
    (sqlc.narg(id)::UUID IS NULL OR id = sqlc.narg(id)::UUID) AND
    (sqlc.narg(campaign_id)::UUID IS NULL OR campaign_id = sqlc.narg(campaign_id)::UUID) AND
    (sqlc.narg(user_id)::UUID IS NULL OR user_id = sqlc.narg(user_id)::UUID)
LIMIT @lim OFFSET @off;

-- name: GetDonationByID :one
SELECT
    id, campaign_id, transaction_id, user_id, amount, donation_date 
FROM
    t_donations
WHERE
    id = @donation_id;

-- name: CreateDonation :one
INSERT INTO t_donations
    (campaign_id, user_id, amount, donation_date, transaction_id)  
VALUES
    (@campaign_id, @user_id, @amount, NOW(), @transaction_id)
RETURNING id;

-- name: DeleteDonation :exec
DELETE FROM
    t_donations
WHERE
    id = @donation_id;

-- name: CountDonations :one
SELECT 
    COUNT(*) 
FROM
    t_donations 
WHERE
    (sqlc.narg(campaign_id)::UUID IS NULL OR campaign_id = sqlc.narg(campaign_id)::UUID) AND
    (sqlc.narg(user_id)::UUID IS NULL OR user_id = sqlc.narg(user_id)::UUID);