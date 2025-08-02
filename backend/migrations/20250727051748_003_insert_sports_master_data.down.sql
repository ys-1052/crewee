-- Delete sports master data (MVP: Basketball and Soccer only)
DELETE FROM sports
WHERE
  code IN ('soccer', 'basketball');
