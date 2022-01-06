SELECT m.id, m.name, m.salary*1000 salary, m.plan*1000 plan, COALESCE(ss.total, 0) total FROM managers m
LEFT JOIN(SELECT s.manager_id, 
    SUM((SELECT SUM(sp.price * sp.qty) FROM sale_positions sp
WHERE s.id = sp.sale_id)) total FROM sales s GROUP BY s.manager_id
 ) ss ON m.id = ss.manager_id
 ORDER BY total DESC
