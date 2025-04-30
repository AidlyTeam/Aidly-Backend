-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS t_categories (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(30) NOT NULL
);

INSERT INTO t_categories (name) VALUES 
('Education'),
('Health'),
('Earthquake'),
('Animal Support'),
('Social Support'),
('Innovation');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS t_categories;
-- +goose StatementEnd
