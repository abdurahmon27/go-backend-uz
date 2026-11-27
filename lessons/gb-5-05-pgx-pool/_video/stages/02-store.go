package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	ID    int64
	Name  string
	Price int64
	Stock int
}

type Store struct {
	pool *pgxpool.Pool
}

func (s *Store) Get(ctx context.Context, id int64) (Product, error) {
	var p Product
	err := s.pool.QueryRow(ctx,
		"SELECT id, name, price, stock FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	return p, err
}
