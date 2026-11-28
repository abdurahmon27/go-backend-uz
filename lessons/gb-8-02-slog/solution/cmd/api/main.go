// api — bozor API'si (GB-6.08'dagi himoyalangan bozor), sozlamasi tashqaridan (GB-8.01, config.go):
// port flag yoki muhitdan, DATABASE_URL faqat muhitdan; notoʻgʻri sozlamada server aniq xato bilan toʻxtaydi.
// GB-8.02: log'lar log/slog bilan, har yozuv bitta JSON qator stdout'ga; har soʻrovning qatorlarida bitta request_id.
//
//	docker compose up -d --wait
//	goose up                     (GOOSE_DRIVER, GOOSE_DBSTRING, GOOSE_MIGRATION_DIR: GB-5.06)
//	docker compose exec db psql -f /lab/seed.sql
//	cp .env.example .env         (bir marta; .env gitga tushmaydi)
//	set -a; . ./.env; set +a     (.env'dagi qatorlarni muhitga eksport qiladi)
//	go build ./cmd/api && ./api -port 8782 > api.log   (log faylga; ikkinchi terminalda: cat api.log)
//	./api -h                     → flag'lar roʻyxati
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

// newServer — bozor API serveri (GB-4.04). Handler — newHandler: middleware zanjiri, ichida newMux (GB-6.01).
// logger — asosiy logger (GB-8.02): run uni bir marta yasaydi va aniq parametr bilan beradi.
func newServer(addr string, st *store.Store, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newHandler(st, logger),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func main() {
	err := run(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return // -h: flag yordam matnini chiqardi, bu xato emas
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}

// run — dasturning hammasi shu yerda (GB-5.10): xato qaytaradi, defer'lar ishlaydi, chiqish kodini main hal qiladi.
func run(args []string) error {
	cfg, err := loadConfig(args)
	if err != nil {
		return err
	}
	// GB-8.02: log — JSON qatorlar, stdout'ga (Twelve-Factor App, XI. Logs). nil — standart sozlama: daraja Info.
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	// slog.Default() ham shu logger: loggerFrom'ning zaxirasi va log paketi ham JSON yozadi, stdout'ga.
	slog.SetDefault(logger)
	logger.Info("config", "config", cfg) // cfg.LogValue() chaqiriladi: parol chiqmaydi
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil { // baza ishlamasa server ishga tushmaydi
		return err
	}
	logger.Debug("Postgres ulandi") // Debug < Info: standart sozlamada yozilmaydi
	srv := newServer(cfg.addr(), store.New(pool), logger)
	logger.Info("bozor API ishga tushdi", "addr", cfg.addr())
	return srv.ListenAndServe()
}
