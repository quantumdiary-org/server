-- +goose Up
-- +goose StatementBegin
CREATE TABLE cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_cache_key ON cache(key);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_cache_expires_at ON cache(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cache;
-- +goose StatementEnd