-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS t_campaigns (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    wallet_address TEXT NOT NULL,
    image_path TEXT, 
    target_amount DECIMAL(18, 8) NOT NULL,
    raised_amount DECIMAL(18, 8) DEFAULT 0,
    accepted_token_symbol VARCHAR(10) NOT NULL DEFAULT 'SOL',
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_valid BOOLEAN NOT NULL DEFAULT FALSE,
    status_type VARCHAR(50) NOT NULL DEFAULT 'normal',
    start_date TIMESTAMPTZ DEFAULT NOW(),
    end_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES t_users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS t_campaigns;
-- +goose StatementEnd
