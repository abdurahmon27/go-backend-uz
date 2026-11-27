// api — bozor buyurtmalar API'si (GB-5.10, 5-modul nazorat nuqtasi): GB-4.08'dagi katalog API,
// endi maʼlumot Postgres'da. Server toʻxtasa ham mahsulotlar va buyurtmalar saqlanib qoladi.
//
//	docker compose up -d --wait
//	goose up                     (GOOSE_DRIVER, GOOSE_DBSTRING, GOOSE_MIGRATION_DIR: GB-5.06)
//	docker compose exec db psql -f /lab/seed.sql
//	export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5510/bozor
//	go run ./cmd/api             → bozor API: http://localhost:8765
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

// newServer — bozor API serveri (GB-4.04): route'lar newMux'da, sarlavhalar uchun vaqt chegarasi.
func newServer(addr string, st *store.Store) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newMux(st),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}

func run() error {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return errors.New("DATABASE_URL yoʻq: export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5510/bozor")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil { // baza ishlamasa server ishga tushmaydi
		return err
	}
	srv := newServer(":8765", store.New(pool))
	fmt.Println("bozor API: http://localhost:8765")
	return srv.ListenAndServe()
}
