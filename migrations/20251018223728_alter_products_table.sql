-- +goose Up
-- +goose StatementBegin

-- Index for filtering by category (case-insensitive)
-- This function-based index perfectly matches your "LOWER(category)" WHERE clause.
CREATE INDEX idx_products_category_lower ON products (LOWER(category));

-- Indexes for sorting
CREATE INDEX idx_products_name ON products (name);
CREATE INDEX idx_products_price_in_pence ON products (price_in_pence);
CREATE INDEX idx_products_stock_quantity ON products (stock_quantity);
CREATE INDEX idx_products_created_at ON products (created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_products_category_lower;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_price_in_pence;
DROP INDEX IF EXISTS idx_products_stock_quantity;
DROP INDEX IF EXISTS idx_products_created_at;

-- +goose StatementEnd