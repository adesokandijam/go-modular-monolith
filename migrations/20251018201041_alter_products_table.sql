-- +goose Up
-- +goose StatementBegin
ALTER TABLE products ADD COLUMN category VARCHAR(100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN category;
-- +goose StatementEnd