-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT *
FROM chirps
WHERE ($1::uuid = '00000000-0000-0000-0000-000000000000' OR user_id = $1)
ORDER BY created_at ASC;


-- name: GetChirpByID :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;