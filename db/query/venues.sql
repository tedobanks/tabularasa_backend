-- name: GetVenue :one
SELECT * FROM venues
WHERE id = $1 LIMIT 1;

-- name: Listvenues :many
SELECT * FROM venues
ORDER BY name;

-- name: CreateVenue :one
INSERT INTO "venues" (
  image_links,
  name,
  type,
  description,
  location,
  dimension,
  capacity,
  facilities,
  has_accomodation,
  room_type,
  no_of_rooms,
  sleeps,
  bed_type,
  rent,
  owned_by,
  is_available,
  opens_at,
  closes_at,
  rental_days,
  booking_price
) VALUES (
  $1::varchar[], $2, $3, $4, $5, $6, $7, $8::varchar[], $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
)
RETURNING *;

-- name: UpdateVenue :exec
UPDATE venues
  set name = $2,
  image_links = $3::varchar[],
  type = $4,
  description = $5,
  location = $6,
  dimension = $7,
  capacity = $8,
  facilities = $9::varchar[],
  has_accomodation = $10,
  room_type = $11,
  no_of_rooms = $12,
  sleeps = $13,
  bed_type = $14,
  rent = $15,
  owned_by = $16,
  is_available = $17,
  opens_at = $18,
  closes_at = $19,
  rental_days = $20,
  booking_price = $21
WHERE id = $1;

-- name: DeleteVenue :exec
DELETE FROM venues
WHERE id = $1;