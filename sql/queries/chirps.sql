-- name: CreateChirps :one
INSERT INTO chirps(created_at,updated_at,body,user_id)
VALUES(
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;


-- name: DeleteAllChirps :exec
DELETE FROM chirps;


-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;


-- name: GetChirp :one
SELECT * 
FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1 AND user_id = $2;