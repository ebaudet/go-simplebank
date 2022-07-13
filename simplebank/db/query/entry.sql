-- CREATE --

-- name: CreateEntry :one
INSERT INTO entries (
  account_id, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- READ --

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- UPDATE --

-- name: UpdateEntry :one
UPDATE entries
set amount = $2
WHERE id = $1
RETURNING *;

-- DELETE --

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;
