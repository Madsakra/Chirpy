-- name: CreateUser :one
INSERT INTO users(created_at,updated_at,email,hashed_password)
VALUES(
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserHash :one
SELECT *
FROM users
WHERE email = $1;


-- name: UpdateUser :exec
UPDATE users
SET email = $1,
    hashed_password = $2
WHERE id = $3;

-- name: UpgradeChirp :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;