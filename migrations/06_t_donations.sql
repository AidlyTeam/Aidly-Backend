-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS t_donations (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL,
    transaction_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    amount DECIMAL(18, 8) NOT NULL,
    donation_date TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT fk_campaign FOREIGN KEY (campaign_id) REFERENCES t_campaigns(id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES t_users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS t_donations;
-- +goose StatementEnd
