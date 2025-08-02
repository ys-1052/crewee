-- name: GetAllRegions :many
SELECT
  jis_code,
  name,
  name_kana,
  region_type,
  parent_jis_code,
  created_at,
  updated_at
FROM
  regions
ORDER BY
  jis_code;


-- name: GetPrefectures :many
SELECT
  jis_code,
  name,
  name_kana,
  region_type,
  parent_jis_code,
  created_at,
  updated_at
FROM
  regions
WHERE
  region_type = 'prefecture'
ORDER BY
  jis_code;


-- name: GetMunicipalitiesByPrefecture :many
SELECT
  jis_code,
  name,
  name_kana,
  region_type,
  parent_jis_code,
  created_at,
  updated_at
FROM
  regions
WHERE
  region_type = 'municipality'
  AND parent_jis_code = $1
ORDER BY
  jis_code;


-- name: GetRegionByCode :one
SELECT
  jis_code,
  name,
  name_kana,
  region_type,
  parent_jis_code,
  created_at,
  updated_at
FROM
  regions
WHERE
  jis_code = $1;


-- name: GetRegionHierarchy :many
SELECT
  r.jis_code,
  r.name,
  r.name_kana,
  r.region_type,
  r.parent_jis_code,
  p.name AS parent_name,
  r.created_at,
  r.updated_at
FROM
  regions r
  LEFT JOIN regions p ON r.parent_jis_code = p.jis_code
ORDER BY
  r.region_type,
  r.jis_code;
