-- name: GetFavourite :one
SELECT * FROM favourites
WHERE id = $1 LIMIT 1;

-- name: ListFavouritesByEvent :many
SELECT * FROM favourites
WHERE event_id = $1
ORDER BY created_at DESC;

-- name: ListFavouritesByUser :many
SELECT * FROM favourites
WHERE added_by = $1
ORDER BY created_at DESC;

-- name: CreateFavourite :one
INSERT INTO "favourites" (
  event_id,
  added_by
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteFavourite :exec
DELETE FROM favourites
WHERE id = $1;

-- name: DeleteFavouriteByUserAndEvent :exec
DELETE FROM favourites
WHERE event_id = $1 AND added_by = $2;
