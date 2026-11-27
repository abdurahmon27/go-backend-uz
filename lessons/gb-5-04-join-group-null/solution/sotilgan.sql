SELECT products.name,
       COALESCE(SUM(order_items.quantity), 0) AS sotilgan
FROM products
LEFT JOIN order_items ON order_items.product_id = products.id
GROUP BY products.id
ORDER BY sotilgan DESC, products.name;
