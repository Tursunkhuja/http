SELECT p.id, p.name, ss.total FROM products p
JOIN (SELECT sp.product_id, SUM(sp.price * sp.qty) total FROM sale_positions sp
    GROUP BY sp.product_id
) ss ON p.id = ss.product_id
ORDER BY ss.total DESC
LIMIT 3;