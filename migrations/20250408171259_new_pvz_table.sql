-- +goose Up
-- +goose StatementBegin
CREATE TYPE city_name AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');

CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    city city_name NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pvz;
-- +goose StatementEnd
