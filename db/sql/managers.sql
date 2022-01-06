SELECT m.id, m.name, m.salary*1000 salary, m.plan*1000 plan, ss.total FROM managers m
JOIN(SELECT s.manager_id, 
    SUM((SELECT SUM(sp.price * sp.qty) FROM sale_positions sp
WHERE s.id = sp.sale_id)) total FROM sales s GROUP BY s.manager_id
 ) ss ON m.id = ss.manager_id
 ORDER BY ss.total DESC
