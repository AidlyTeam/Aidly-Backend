-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS t_badges (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    description TEXT,
    icon_path TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS t_badges;
-- +goose StatementEnd
