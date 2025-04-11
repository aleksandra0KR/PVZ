-- +goose Up
-- +goose StatementBegin
CREATE TYPE status AS ENUM ('in_progress', 'close');

CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    status status NOT NULL DEFAULT status('in_progress')
);
CREATE INDEX idx_receptions_pvz_id_date ON receptions (pvz_id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS receptions;
-- +goose StatementEnd
