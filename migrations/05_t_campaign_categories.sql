-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS t_campaign_categories (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL,
    category_id UUID NOT NULL,

    CONSTRAINT fk_campaign FOREIGN KEY (campaign_id) REFERENCES t_campaigns(id) ON DELETE CASCADE,
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES t_categories(id) ON DELETE CASCADE,
    CONSTRAINT uq_campaign_category UNIQUE (campaign_id, category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS t_campaign_categories;
-- +goose StatementEnd
