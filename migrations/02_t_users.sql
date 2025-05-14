-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS t_users (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL,
    wallet_address VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    name VARCHAR(30),
    surname VARCHAR(30),
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES t_roles(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS t_users;
-- +goose StatementEnd
