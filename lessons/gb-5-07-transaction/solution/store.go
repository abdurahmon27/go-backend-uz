package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Product — products jadvalidagi bitta qator (narx tiyinda, GB-1.03).
type Product struct {
	ID    int64
	Name  string
	Price int64
	Stock int
}

// ErrNotFound — jadvalda bunday mahsulot yoʻq (GB-2.08: chaqiruvchi errors.Is bilan tekshiradi).
var ErrNotFound = errors.New("mahsulot topilmadi")

// Store — products jadvaliga Go'dan kirish. pool bir vaqtda koʻp goroutine'dan ishlatilishi mumkin.
type Store struct {
	pool *pgxpool.Pool
}

// Get id boʻyicha bitta mahsulotni oʻqiydi; yoʻq boʻlsa ErrNotFound.
func (s *Store) Get(ctx context.Context, id int64) (Product, error) {
	var p Product
	err := s.pool.QueryRow(ctx,
		"SELECT id, name, price, stock FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	return p, err
}

// ListActive faol mahsulotlarni id tartibida oʻqiydi.
func (s *Store) ListActive(ctx context.Context) ([]Product, error) {
	rows, err := s.pool.Query(ctx,
		"SELECT id, name, price, stock FROM products WHERE active ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// OrderTotal bitta buyurtmaning jami summasini (tiyinda) hisoblaydi. Buyurtma yoʻq boʻlsa SUM NULL beradi: COALESCE uni nolga aylantiradi.
func (s *Store) OrderTotal(ctx context.Context, orderID int64) (int64, error) {
	var total int64
	err := s.pool.QueryRow(ctx,
		"SELECT COALESCE(SUM(quantity*unit_price), 0) FROM order_items WHERE order_id = $1", orderID).
		Scan(&total)
	return total, err
}

// ByName nom boʻyicha bitta mahsulotni oʻqiydi; yoʻq boʻlsa ErrNotFound.
func (s *Store) ByName(ctx context.Context, name string) (Product, error) {
	var p Product
	err := s.pool.QueryRow(ctx,
		"SELECT id, name, price, stock FROM products WHERE name = $1", name).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	return p, err
}

// Cheapest eng arzon n ta faol mahsulotni narx boʻyicha oʻsish tartibida oʻqiydi.
func (s *Store) Cheapest(ctx context.Context, n int) ([]Product, error) {
	rows, err := s.pool.Query(ctx,
		"SELECT id, name, price, stock FROM products WHERE active ORDER BY price, id LIMIT $1", n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ErrNotEnoughStock — omborda yetarli qoldiq yoʻq: baza CHECK (stock >= 0) bilan rad etdi.
var ErrNotEnoughStock = errors.New("qoldiq yetarli emas")

// ErrInvalidQuantity — miqdor noldan katta boʻlishi kerak.
var ErrInvalidQuantity = errors.New("miqdor noldan katta boʻlishi kerak")

// OrderItem — buyurtmadagi bitta qator: qaysi mahsulot va nechta.
type OrderItem struct {
	ProductID int64
	Quantity  int
}

// PlaceOrder buyurtmani bitta transaction'da yozadi: orders qatori, har mahsulot uchun products.stock kamayishi
// va order_items qatori. Biror qadam xato bersa, hech narsa yozilmaydi.
func (s *Store) PlaceOrder(ctx context.Context, customer string, items []OrderItem) (int64, error) {
	for _, it := range items {
		if it.Quantity <= 0 {
			return 0, ErrInvalidQuantity
		}
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var orderID int64
	err = tx.QueryRow(ctx,
		"INSERT INTO orders (customer) VALUES ($1) RETURNING id", customer).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	for _, it := range items {
		var price int64
		err = tx.QueryRow(ctx,
			"UPDATE products SET stock = stock - $1 WHERE id = $2 RETURNING price",
			it.Quantity, it.ProductID).Scan(&price)
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return 0, ErrNotEnoughStock
		}
		if err != nil {
			return 0, err
		}
		_, err = tx.Exec(ctx,
			"INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES ($1, $2, $3, $4)",
			orderID, it.ProductID, it.Quantity, price)
		if err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return orderID, nil
}

// Restock tovar qabul qilishni yozadi: har mahsulotning qoldigʻi Quantity ga oshadi.
// Hammasi bitta transaction'da: bitta mahsulot topilmasa, hech birining qoldigʻi oʻzgarmaydi.
func (s *Store) Restock(ctx context.Context, items []OrderItem) error {
	for _, it := range items {
		if it.Quantity <= 0 {
			return ErrInvalidQuantity
		}
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, it := range items {
		var stock int
		err = tx.QueryRow(ctx,
			"UPDATE products SET stock = stock + $1 WHERE id = $2 RETURNING stock",
			it.Quantity, it.ProductID).Scan(&stock)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
