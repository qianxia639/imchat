-- name: AddContact :one
INSERT INTO contact (
    owner_id, target_id, type
) VALUES (
    $1, $2, $3
)
RETURNING *;