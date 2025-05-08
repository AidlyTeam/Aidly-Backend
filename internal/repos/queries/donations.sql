-- name: GetDonations :many
SELECT
    d.id, 
    d.campaign_id, 
    d.transaction_id, 
    d.user_id, 
    d.amount, 
    d.donation_date,
    c.title AS campaign_title
FROM
    t_donations d
JOIN 
    t_campaigns c ON c.id = d.campaign_id
WHERE
    (sqlc.narg(id)::UUID IS NULL OR d.id = sqlc.narg(id)::UUID) AND
    (sqlc.narg(campaign_id)::UUID IS NULL OR d.campaign_id = sqlc.narg(campaign_id)::UUID) AND
    (sqlc.narg(user_id)::UUID IS NULL OR d.user_id = sqlc.narg(user_id)::UUID)
LIMIT @lim OFFSET @off;

-- name: GetDonationByID :one
SELECT
    d.id, 
    d.campaign_id, 
    d.transaction_id, 
    d.user_id, 
    d.amount, 
    d.donation_date,
    c.title AS campaign_title
FROM
    t_donations d
JOIN 
    t_campaigns c ON c.id = d.campaign_id
WHERE
    d.id = @donation_id;

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