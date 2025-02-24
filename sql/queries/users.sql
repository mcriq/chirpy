-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
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

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;