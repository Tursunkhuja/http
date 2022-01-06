UPDATE customers
SET phone = '+992000000011'
WHERE id = 11
RETURNING id, name, active;