-- name: GetPractitioner :one
SELECT * FROM practitioners
WHERE id = $1 LIMIT 1;

-- name: ListPractitioners :many
SELECT * FROM practitioners
ORDER BY name;

-- name: CreatePractitioner :one
INSERT INTO "practitioners" (
  name,
  description,
  image_link,
  is_available,
  created_by,
  opens_at,
  closes_at,
  working_days
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdatePractitioner :one
UPDATE practitioners
  set name = $2,
  description = $3,
  image_link = $4,
  is_available = $5,
  created_by = $6,
  opens_at = $7,
  closes_at = $8,
  working_days = $9
WHERE id = $1
RETURNING *;

-- name: DeletePractitioner :exec
DELETE FROM practitioners
WHERE id = $1;
