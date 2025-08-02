-- Insert sports master data (MVP: Basketball and Soccer only)
INSERT INTO
  sports (id, code, name, is_active)
VALUES
  (
    '01234567-89ab-cdef-0123-456789abcdef',
    'soccer',
    'サッカー',
    TRUE
  ),
  (
    '21234567-89ab-cdef-0123-456789abcdef',
    'basketball',
    'バスケットボール',
    TRUE
  );
