-- name: CreateChirp :one
INSERT INTO chirps (id , created_at, updated_at , body , user_id)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirpsOrderByCreatedAtAsc :many
SELECT * 
FROM chirps AS c
ORDER BY c.created_at ASC;

-- name: GetChirp :one
SELECT *
FROM chirps AS c
WHERE c.id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE chirps.id = $1;