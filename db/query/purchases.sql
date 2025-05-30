-- name: GetPurchase :one
SELECT * FROM purchases
WHERE id = $1 LIMIT 1;

-- name: ListPurchasesByUser :many
SELECT * FROM purchases
WHERE purchased_by = $1
ORDER BY created_at DESC;

-- name: ListPurchasesByEvent :many
SELECT * FROM purchases
WHERE event_id = $1
ORDER BY created_at DESC;

-- name: ListPurchasesByVenue :many
SELECT * FROM purchases
WHERE venue_id = $1
ORDER BY created_at DESC;

-- name: ListPurchasesByService :many
SELECT * FROM purchases
WHERE service_id = $1
ORDER BY created_at DESC;

-- name: CreatePurchase :one
INSERT INTO "purchases" (
  event_id,
  venue_id,
  service_id,
  purchased_by
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeletePurchase :exec
DELETE FROM purchases
WHERE id = $1;
