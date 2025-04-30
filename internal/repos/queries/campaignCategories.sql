-- name: CreateCampaignCategory :one
INSERT INTO t_campaign_categories (campaign_id, category_id)
VALUES (@campaign_id, @category_id)
RETURNING id;

-- name: DeleteCampaignCategory :exec
DELETE FROM t_campaign_categories
WHERE campaign_id = @campaign_id AND category_id = @category_id;
