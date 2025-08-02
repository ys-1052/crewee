-- name: GetAllSports :many
SELECT
  id,
  code,
  name,
  is_active,
  created_at,
  updated_at
FROM
  sports
ORDER BY
  name;


-- name: GetActiveSports :many
SELECT
  id,
  code,
  name,
  is_active,
  created_at,
  updated_at
FROM
  sports
WHERE
  is_active = TRUE
ORDER BY
  name;


-- name: GetSportByID :one
SELECT
  id,
  code,
  name,
  is_active,
  created_at,
  updated_at
FROM
  sports
WHERE
  id = $1;


-- name: GetSportByCode :one
SELECT
  id,
  code,
  name,
  is_active,
  created_at,
  updated_at
FROM
  sports
WHERE
  code = $1;
