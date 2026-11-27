package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// firstName: birinchi mahsulotning nomi
func firstName(ctx context.Context, pool *pgxpool.Pool) (string, error) {
	rows, err := pool.Query(ctx, "SELECT name FROM products ORDER BY id")
	if err != nil {
		return "", err
	}
	if !rows.Next() {
		return "", rows.Err()
	}
	var name string
	err = rows.Scan(&name)
	return name, err
}

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "pool:", err)
		os.Exit(1)
	}
	qctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	for i := 1; i <= 2; i++ {
		start := time.Now()
		name, err := firstName(qctx, pool)
		fmt.Println(i, name, err, time.Since(start).Round(time.Second))
	}
}
