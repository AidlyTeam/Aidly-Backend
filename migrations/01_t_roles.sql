-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS t_roles (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(30) NOT NULL
);

INSERT INTO t_roles (name) VALUES ('public');
INSERT INTO t_roles (name) VALUES ('first');
INSERT INTO t_roles (name) VALUES ('user');
INSERT INTO t_roles (name) VALUES ('admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS t_roles;
-- +goose StatementEnd
