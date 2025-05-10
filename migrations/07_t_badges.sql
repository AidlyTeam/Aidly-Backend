-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS t_badges (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    symbol VARCHAR(50) DEFAULT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    seller_fee INT DEFAULT NULL,
    icon_path TEXT,
    donation_threshold INTEGER UNIQUE NOT NULL,
    uri VARCHAR(130) DEFAULT NULL,
    is_nft BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS t_badges;
-- +goose StatementEnd
