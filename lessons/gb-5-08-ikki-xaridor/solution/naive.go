package bozor

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// buyNaive — videoning birinchi, XATO versiyasi: transaction bor, lekin oʻqish va yozish orasida qulf yoʻq.
// Ikki xaridor bir xil eski qoldiqni oʻqib, ikkalasi ham yozadi (lost update). Faqat koʻrsatish uchun (naive_test.go): ishlatmang.
func (s *Store) buyNaive(ctx context.Context, id int64, qty int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var stock int
	err = tx.QueryRow(ctx,
		"SELECT stock FROM products WHERE id = $1", id).Scan(&stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if stock < qty {
		return ErrNotEnoughStock
	}
	_, err = tx.Exec(ctx,
		"UPDATE products SET stock = $1 WHERE id = $2", stock-qty, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
