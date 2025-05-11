-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS t_user_badges (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    badge_id UUID NOT NULL,
    is_minted BOOLEAN NOT NULL DEFAULT FALSE,
    awarded_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES t_users(id) ON DELETE CASCADE,
    CONSTRAINT fk_badge FOREIGN KEY (badge_id) REFERENCES t_badges(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_badge UNIQUE (user_id, badge_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS t_user_badges;

-- +goose StatementEnd
