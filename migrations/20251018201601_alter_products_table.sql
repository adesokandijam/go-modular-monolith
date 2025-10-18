-- +goose Up
-- +goose StatementBegin
ALTER TABLE products ALTER COLUMN category SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products ALTER COLUMN category DROP NOT NULL;
-- +goose StatementEnd