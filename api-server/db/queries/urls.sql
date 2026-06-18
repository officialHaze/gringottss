-- name: AddURL :one
INSERT INTO urls (id, url)
VALUES (@id, @url)
ON CONFLICT(url) DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: ListURLs :many
SELECT * FROM urls
ORDER BY updated_at DESC;

-- name: DeleteURL :exec
DELETE FROM urls
WHERE url = @url;