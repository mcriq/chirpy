-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: UserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE id = $1
);