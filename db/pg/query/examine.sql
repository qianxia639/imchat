-- name: AddExamine :one
INSERT INTO examine (
    owner_id, target_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetExamine :many
SELECT * FROM examine
WHERE owner_id = $1;

-- name: DeleteExamine :exec
DELETE FROM examine 
WHERE owner_id = $1 AND target_id = $2 AND type = $3;