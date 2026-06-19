-- +goose Up
CREATE UNIQUE INDEX credentials_unique_idx_url_forminputid ON credentials (url, form_input_id);

-- Index url with updated_at value with id as tie breaker, sort credentials by url, with update time and id as tie-breaker
CREATE INDEX credentials_idx_url_updatedat_id ON credentials (url ASC, updated_at DESC, id ASC);

-- Sort urls by updated time (latest to oldest) with id as tie breaker
CREATE INDEX urls_idx_updatedat_id ON urls (updated_at DESC, id ASC);

-- +goose Down
DROP INDEX IF EXISTS urls_idx_updatedat_id;
DROP INDEX IF EXISTS credentials_idx_url_updatedat_id;
DROP INDEX IF EXISTS credentials_unique_idx_url_forminputid;