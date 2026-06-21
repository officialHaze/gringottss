-- name: AddCredential :one
INSERT INTO credentials (id, url, form_input_id, form_input_name, form_input_xpath, form_input_type, form_input_val)
VALUES (@id, @url, @form_input_id, @form_input_name, @form_input_xpath, @form_input_type, @form_input_val)
ON CONFLICT (url, form_input_id, form_input_xpath)
DO UPDATE SET
  form_input_name = excluded.form_input_name,
  form_input_type = excluded.form_input_type,
  form_input_val = excluded.form_input_val
RETURNING *;

-- name: GetCredentials :many
SELECT * FROM credentials
WHERE url = @url OR id = @id
ORDER BY url ASC, updated_at DESC, id ASC;

-- name: DeleteCredentials :exec
DELETE FROM credentials
WHERE url = @url OR id = @id;

-- name: UpdateCredential :one
UPDATE credentials
SET
  form_input_id = @form_input_id,
  form_input_name = @form_input_name,
  form_input_xpath = @form_input_xpath,
  form_input_type = @form_input_type,
  form_input_val = @form_input_val
WHERE url = @url
RETURNING *;
