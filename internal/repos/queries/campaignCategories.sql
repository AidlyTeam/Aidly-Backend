-- name: GetCampaignCategories :many
SELECT
    tc.category_id,
    tc.campaign_id,
    cat.name AS category_name
FROM t_campaign_categories tc
JOIN t_categories cat ON cat.id = tc.category_id
WHERE tc.campaign_id = @campaign_id
LIMIT @lim OFFSET @off;

-- name: GetCampaignCategoriesOne :one
SELECT
    id, campaign_id, category_id
FROM
    t_campaign_categories
WHERE
    campaign_id = @campaign_id AND
    category_id = @category_id;

-- name: CreateCampaignCategory :one
INSERT INTO t_campaign_categories (campaign_id, category_id)
VALUES (@campaign_id, @category_id)
RETURNING id;

-- name: DeleteCampaignCategory :exec
DELETE FROM t_campaign_categories
WHERE campaign_id = @campaign_id AND category_id = @category_id;
