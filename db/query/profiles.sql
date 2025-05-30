-- name: GetProfile :one
SELECT * FROM profiles
WHERE id = $1 LIMIT 1;

-- name: ListProfiles :many
SELECT * FROM profiles
ORDER BY business_name; -- Or any other relevant field for ordering

-- name: CreateProfile :one
INSERT INTO "profiles" (
  id,
  bio,
  phone_no,
  country,
  address,
  experience,
  field,
  business_name,
  roles
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles
  set bio = $2,
  phone_no = $3,
  country = $4,
  address = $5,
  experience = $6,
  field = $7,
  business_name = $8,
  roles = $9
WHERE id = $1
RETURNING *;

-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE id = $1;
