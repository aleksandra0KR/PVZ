-- +goose Up
-- +goose StatementBegin
CREATE TYPE product_type AS ENUM ('электроника', 'одежда', 'обувь');

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    type product_type NOT NULL,
    reception_id UUID NOT NULL REFERENCES receptions(id) ON DELETE CASCADE
    );

CREATE INDEX idx_products_reception_id_date ON products (reception_id, date_time DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
