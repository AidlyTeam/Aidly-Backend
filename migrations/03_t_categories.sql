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
('Innovation'),
('Environment'),
('Disability Support'),
('Children'),
('Elderly Care'),
('Arts & Culture'),
('Emergency Relief'),
('Mental Health'),
('Technology Access'),
('Clean Water'),
('Food Aid'),
('Housing Support'),
('Refugee Support'),
('Medical Expenses'),
('Funeral Assistance'),
('Debt Relief'),
('Single Parent Support'),
('Unemployment Aid'),
('Accident Recovery'),
('Rent Support'),
('Utility Bills Assistance'),
('Surgery Funding'),
('Newborn Support');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS t_categories;
-- +goose StatementEnd
