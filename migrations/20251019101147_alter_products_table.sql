-- +goose Up
-- +goose StatementBegin

-- Drop existing indexes
DROP INDEX IF EXISTS idx_products_category_lower;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_price_in_pence;
DROP INDEX IF EXISTS idx_products_stock_quantity;
DROP INDEX IF EXISTS idx_products_created_at;

-- Composite index for category + sort columns (most common queries)
-- This covers: WHERE category AND ORDER BY sort_column
CREATE INDEX idx_products_category_stock ON products (LOWER(category), stock_quantity DESC, id);
CREATE INDEX idx_products_category_price ON products (LOWER(category), price_in_pence DESC, id);
CREATE INDEX idx_products_category_created ON products (LOWER(category), created_at DESC, id);
CREATE INDEX idx_products_category_name ON products (LOWER(category), name ASC, id);

-- Fallback single-column indexes for other sort scenarios
CREATE INDEX idx_products_stock_quantity ON products (stock_quantity DESC);
CREATE INDEX idx_products_price_in_pence ON products (price_in_pence DESC);
CREATE INDEX idx_products_created_at ON products (created_at DESC);
CREATE INDEX idx_products_name ON products (name ASC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_products_category_stock;
DROP INDEX IF EXISTS idx_products_category_price;
DROP INDEX IF EXISTS idx_products_category_created;
DROP INDEX IF EXISTS idx_products_category_name;
DROP INDEX IF EXISTS idx_products_stock_quantity;
DROP INDEX IF EXISTS idx_products_price_in_pence;
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_products_name;

-- +goose StatementEnd