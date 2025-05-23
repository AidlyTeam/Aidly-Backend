// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: campaignCategories.sql

package repo

import (
	"context"

	"github.com/google/uuid"
)

const createCampaignCategory = `-- name: CreateCampaignCategory :one
INSERT INTO t_campaign_categories (campaign_id, category_id)
VALUES ($1, $2)
RETURNING id
`

type CreateCampaignCategoryParams struct {
	CampaignID uuid.UUID
	CategoryID uuid.UUID
}

func (q *Queries) CreateCampaignCategory(ctx context.Context, arg CreateCampaignCategoryParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createCampaignCategory, arg.CampaignID, arg.CategoryID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteCampaignCategory = `-- name: DeleteCampaignCategory :exec
DELETE FROM t_campaign_categories
WHERE campaign_id = $1 AND category_id = $2
`

type DeleteCampaignCategoryParams struct {
	CampaignID uuid.UUID
	CategoryID uuid.UUID
}

func (q *Queries) DeleteCampaignCategory(ctx context.Context, arg DeleteCampaignCategoryParams) error {
	_, err := q.db.ExecContext(ctx, deleteCampaignCategory, arg.CampaignID, arg.CategoryID)
	return err
}

const getCampaignCategories = `-- name: GetCampaignCategories :many
SELECT
    tc.category_id,
    tc.campaign_id,
    cat.name AS category_name
FROM t_campaign_categories tc
JOIN t_categories cat ON cat.id = tc.category_id
WHERE tc.campaign_id = $1
LIMIT $3 OFFSET $2
`

type GetCampaignCategoriesParams struct {
	CampaignID uuid.UUID
	Off        int32
	Lim        int32
}

type GetCampaignCategoriesRow struct {
	CategoryID   uuid.UUID
	CampaignID   uuid.UUID
	CategoryName string
}

func (q *Queries) GetCampaignCategories(ctx context.Context, arg GetCampaignCategoriesParams) ([]GetCampaignCategoriesRow, error) {
	rows, err := q.db.QueryContext(ctx, getCampaignCategories, arg.CampaignID, arg.Off, arg.Lim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCampaignCategoriesRow
	for rows.Next() {
		var i GetCampaignCategoriesRow
		if err := rows.Scan(&i.CategoryID, &i.CampaignID, &i.CategoryName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCampaignCategoriesOne = `-- name: GetCampaignCategoriesOne :one
SELECT
    id, campaign_id, category_id
FROM
    t_campaign_categories
WHERE
    campaign_id = $1 AND
    category_id = $2
`

type GetCampaignCategoriesOneParams struct {
	CampaignID uuid.UUID
	CategoryID uuid.UUID
}

func (q *Queries) GetCampaignCategoriesOne(ctx context.Context, arg GetCampaignCategoriesOneParams) (TCampaignCategory, error) {
	row := q.db.QueryRowContext(ctx, getCampaignCategoriesOne, arg.CampaignID, arg.CategoryID)
	var i TCampaignCategory
	err := row.Scan(&i.ID, &i.CampaignID, &i.CategoryID)
	return i, err
}
