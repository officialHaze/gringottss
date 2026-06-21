-- +goose Up
CREATE TABLE urls (
    id  TEXT PRIMARY KEY NOT NULL,
    url TEXT UNIQUE NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE credentials (
    id  TEXT PRIMARY KEY NOT NULL,
    url TEXT NOT NULL REFERENCES urls(url) ON DELETE CASCADE,
    form_input_id TEXT NOT NULL,
    form_input_name TEXT,
    form_input_type TEXT NOT NULL DEFAULT 'text',
    form_input_val TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementBegin
CREATE TRIGGER urls_updated_at
    BEFORE UPDATE ON urls
    FOR EACH ROW
    BEGIN UPDATE urls SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER credentials_updated_at
    BEFORE UPDATE ON credentials
    FOR EACH ROW
    BEGIN UPDATE credentials SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS credentials_updated_at;
DROP TRIGGER IF EXISTS urls_updated_at;
DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS urls;