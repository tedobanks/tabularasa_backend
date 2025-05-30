-- name: GetBookedPractitioner :one
SELECT * FROM "bookedPractitioners"
WHERE id = $1 LIMIT 1;

-- name: ListBookedPractitionersByService :many
SELECT * FROM "bookedPractitioners"
WHERE service_id = $1
ORDER BY booked_for;

-- name: ListBookedPractitionersByUser :many
SELECT * FROM "bookedPractitioners"
WHERE booked_by = $1
ORDER BY booked_for;

-- name: CreateBookedPractitioner :one
INSERT INTO "bookedPractitioners" (
  type,
  service_id,
  booked_for,
  booked_by
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateBookedPractitioner :one
UPDATE "bookedPractitioners"
  set type = $2,
  service_id = $3,
  booked_for = $4,
  booked_by = $5
WHERE id = $1
RETURNING *;

-- name: DeleteBookedPractitioner :exec
DELETE FROM "bookedPractitioners"
WHERE id = $1;
