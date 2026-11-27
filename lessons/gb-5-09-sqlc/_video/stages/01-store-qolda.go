package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Videoning boshidagi holat: GB-5.05'dagi qoʻlda yozilgan Get, ustun nomida xato bilan (stock → qoldiq).
type Product struct {
	ID    int64
	Name  string
	Price int64
	Stock int
}

var ErrNotFound = errors.New("mahsulot topilmadi")

type Store struct {
	pool *pgxpool.Pool
}

func (s *Store) Get(ctx context.Context, id int64) (Product, error) {
	var p Product
	err := s.pool.QueryRow(ctx,
		"SELECT id, name, price, qoldiq FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	return p, err
}
