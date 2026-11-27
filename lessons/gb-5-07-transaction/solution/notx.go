package main

import "context"

// PlaceOrderNoTx — videodagi XATO versiya: shu PlaceOrder, lekin transaction'siz.
// Har buyruq pool orqali alohida yuboriladi, shuning uchun xatodan keyin yarim yozilgan buyurtma qoladi.
// Faqat koʻrsatish uchun: ishlatmang.
func (s *Store) PlaceOrderNoTx(ctx context.Context, customer string, items []OrderItem) (int64, error) {
	var orderID int64
	err := s.pool.QueryRow(ctx,
		"INSERT INTO orders (customer) VALUES ($1) RETURNING id", customer).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	for _, it := range items {
		var price int64
		err = s.pool.QueryRow(ctx,
			"UPDATE products SET stock = stock - $1 WHERE id = $2 RETURNING price",
			it.Quantity, it.ProductID).Scan(&price)
		if err != nil {
			return 0, err
		}
		_, err = s.pool.Exec(ctx,
			"INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES ($1, $2, $3, $4)",
			orderID, it.ProductID, it.Quantity, price)
		if err != nil {
			return 0, err
		}
	}
	return orderID, nil
}
