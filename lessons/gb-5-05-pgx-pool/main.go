package main

import (
	"context"
	"fmt"
	"os"
	"time"

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

	qctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	p, err := store.Get(qctx, 1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	fmt.Println(p.ID, p.Name, p.Price, p.Stock)

	products, err := store.ListActive(qctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	for _, p := range products {
		fmt.Println(p.ID, p.Name, p.Price, p.Stock)
	}

	fmt.Println(store.OrderTotal(qctx, 1))
	fmt.Println(store.OrderTotal(qctx, 99))
}
