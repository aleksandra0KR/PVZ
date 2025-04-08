-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
);
CREATE INDEX idx_receptions_pvz_id_date ON receptions (pvz_id, date_time DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS receptions;
-- +goose StatementEnd
