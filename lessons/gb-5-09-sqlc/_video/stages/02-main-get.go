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

	store := NewStore(pool)
	qctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	p, err := store.Get(qctx, 1)
	fmt.Println(p, err)
}
