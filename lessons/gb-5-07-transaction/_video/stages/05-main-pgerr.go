package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "pool:", err)
		os.Exit(1)
	}
	defer pool.Close()

	store := &Store{pool: pool}
	qctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = store.PlaceOrder(qctx, "Aziza", []OrderItem{{ProductID: 2, Quantity: 5}})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		fmt.Println(pgErr.Code, pgErr.ConstraintName)
	}
}
