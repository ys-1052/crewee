-- name: GetUserByID :one
SELECT
  id,
  email,
  name,
  email_verified_at,
  home_region_code,
  created_at,
  updated_at
FROM
  users
WHERE
  id = $1;


-- name: GetUserByEmail :one
SELECT
  id,
  email,
  name,
  email_verified_at,
  home_region_code,
  created_at,
  updated_at
FROM
  users
WHERE
  email = $1;


-- name: CreateUser :one
INSERT INTO
  users (
    id,
    email,
    name,
    home_region_code,
    email_verified_at
  )
VALUES
  ($1, $2, $3, $4, $5)
RETURNING
  id,
  email,
  name,
  email_verified_at,
  home_region_code,
  created_at,
  updated_at;


-- name: UpdateUser :one
UPDATE users
SET
  name = $2,
  home_region_code = $3,
  updated_at = current_timestamp
WHERE
  id = $1
RETURNING
  id,
  email,
  name,
  email_verified_at,
  home_region_code,
  created_at,
  updated_at;


-- name: UpdateUserEmailVerified :one
UPDATE users
SET
  email_verified_at = $2,
  updated_at = current_timestamp
WHERE
  id = $1
RETURNING
  id,
  email,
  name,
  email_verified_at,
  home_region_code,
  created_at,
  updated_at;
