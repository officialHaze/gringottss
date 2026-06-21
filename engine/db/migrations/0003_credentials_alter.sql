-- +goose Up
-- Add a new field on credentials
ALTER TABLE credentials
ADD COLUMN form_input_xpath TEXT NOT NULL DEFAULT '';

-- Drop the existing unique index on credentials
DROP INDEX IF EXISTS credentials_unique_idx_url_forminputid;

-- Create new unique index on credentials
CREATE UNIQUE INDEX credentials_unique_idx_url_inputid_inputxpath ON credentials (url, form_input_id, form_input_xpath);


-- +goose Down
DROP INDEX IF EXISTS credentials_unique_idx_url_inputid_inputxpath;
CREATE UNIQUE INDEX credentials_unique_idx_url_forminputid ON credentials (url, form_input_id);
ALTER TABLE credentials DROP COLUMN form_input_xpath;