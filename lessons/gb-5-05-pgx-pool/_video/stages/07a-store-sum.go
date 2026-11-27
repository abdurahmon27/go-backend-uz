package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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
		"SELECT SUM(quantity*unit_price) FROM order_items WHERE order_id = $1", orderID).
		Scan(&total)
	return total, err
}

// ByName nom boʻyicha bitta mahsulotni oʻqiydi; yoʻq boʻlsa ErrNotFound. VAZIFA (TASK.md, 1-qadam).
func (s *Store) ByName(ctx context.Context, name string) (Product, error) {
	return Product{}, errors.New("ByName hali yozilmagan")
}

// Cheapest eng arzon n ta faol mahsulotni narx boʻyicha oʻsish tartibida oʻqiydi. VAZIFA (TASK.md, 2-qadam).
func (s *Store) Cheapest(ctx context.Context, n int) ([]Product, error) {
	return nil, errors.New("Cheapest hali yozilmagan")
}
