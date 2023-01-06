-- name: CreateUser :one
INSERT INTO users (
    username, email, nickname, password, gender
) VALUES(
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1
OR email = $1 LIMIT 1;

-- name: LoginUser :one
SELECT * FROM users
WHERE (username = $1 AND password = $2)
OR (email = $1 AND password = $2) LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET
    email = COALESCE(sqlc.narg(email), email),
    nickname = COALESCE(sqlc.narg(nickname), nickname),
    password = COALESCE(sqlc.narg(password), password),
    gender = COALESCE(sqlc.narg(gender), gender),
    avatar = COALESCE(sqlc.narg(avatar), avatar)
WHERE
    username = sqlc.arg(username);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;