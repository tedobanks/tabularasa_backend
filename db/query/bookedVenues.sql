-- name: GetBookedVenue :one
SELECT * FROM "bookedVenues"
WHERE id = $1 LIMIT 1;

-- name: ListBookedVenuesByVenue :many
SELECT * FROM "bookedVenues"
WHERE venue_id = $1
ORDER BY booked_for;

-- name: ListBookedVenuesByUser :many
SELECT * FROM "bookedVenues"
WHERE booked_by = $1
ORDER BY booked_for;

-- name: CreateBookedVenue :one
INSERT INTO "bookedVenues" (
  type,
  venue_id,
  booked_for,
  booked_by
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateBookedVenue :one
UPDATE "bookedVenues"
  set type = $2,
  venue_id = $3,
  booked_for = $4,
  booked_by = $5
WHERE id = $1
RETURNING *;

-- name: DeleteBookedVenue :exec
DELETE FROM "bookedVenues"
WHERE id = $1;
