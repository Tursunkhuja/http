SELECT id, name, qty FROM products 
WHERE active AND qty > 0
ORDER BY qty DESC
LIMIT 10