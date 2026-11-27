// Video bosqichi: Get va ListActive sqlc orqali (query.sql: 02-query.sql).
package main

import (
	"context"
	"errors"

	"github.com/abdurahmon27/bozor/db"
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

// Store — bozor jadvallariga Go'dan kirish. pool — transaction ochish uchun (Begin),
// q — sqlc query.sql'dan generatsiya qilgan soʻrovlar (db papkasi).
type Store struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

// NewStore pool ustida Store yaratadi: db.New pool'ni sqlc soʻrovlariga ulaydi.
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool, q: db.New(pool)}
}

// Get id boʻyicha bitta mahsulotni oʻqiydi; yoʻq boʻlsa ErrNotFound.
func (s *Store) Get(ctx context.Context, id int64) (Product, error) {
	r, err := s.q.GetProduct(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, err
	}
	return Product{ID: r.ID, Name: r.Name, Price: r.Price, Stock: int(r.Stock)}, nil
}

// ListActive faol mahsulotlarni id tartibida oʻqiydi. Mahsulot yoʻq boʻlsa boʻsh (nil emas) slice qaytaradi.
func (s *Store) ListActive(ctx context.Context) ([]Product, error) {
	rows, err := s.q.ListActiveProducts(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Product, 0, len(rows))
	for _, r := range rows {
		out = append(out, Product{ID: r.ID, Name: r.Name, Price: r.Price, Stock: int(r.Stock)})
	}
	return out, nil
}
