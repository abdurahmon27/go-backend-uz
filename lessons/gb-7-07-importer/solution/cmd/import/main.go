// import — katta katalog faylini Postgres'ga yozadi (GB-7.07, 7-modul nazorat nuqtasi).
//
//	docker compose up -d --wait
//	goose up                     (GOOSE_DRIVER, GOOSE_DBSTRING, GOOSE_MIGRATION_DIR: GB-5.06)
//	export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5517/bozor
//	go run ./cmd/katalog 200000 > katta.txt
//	go run ./cmd/import -workers 4 -timeout 10s katta.txt
//
// Ctrl+C importni toza toʻxtatadi: yangi partiya boshlanmaydi, Exec paytida uzilgan partiyalar rollback boʻladi,
// hisobot nechta qator yozilgani va nechtasi rollback boʻlganini aytadi (boʻsh jadvalga importda = count(*)).
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/abdurahmon27/bozor/internal/importer"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	workers := flag.Int("workers", 4, "bir vaqtda nechta partiya yoziladi")
	timeout := flag.Duration("timeout", 0, "bitta partiyaga eng koʻp vaqt, masalan 10s (0 — cheklovsiz)")
	flag.Parse()
	cfg := importer.Config{Workers: *workers, BatchTimeout: *timeout}
	if err := run(flag.Arg(0), cfg); err != nil {
		fmt.Fprintln(os.Stderr, "import:", err)
		os.Exit(1)
	}
}

func run(path string, cfg importer.Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	context.AfterFunc(ctx, stop) // birinchi Ctrl+C'dan keyin ushlash oʻchadi: ikkinchisi dasturni darhol toʻxtatadi

	if path == "" {
		return errors.New("fayl nomi kerak: go run ./cmd/import [-workers 4] katta.txt")
	}

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return errors.New("DATABASE_URL yoʻq: export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5517/bozor")
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return err
	}
	defer pool.Close()

	start := time.Now()
	st, err := importer.Run(ctx, f, importer.NewPG(pool), cfg)
	fmt.Printf("%v · vaqt: %v\n", st, time.Since(start).Round(time.Millisecond))
	if errors.Is(err, context.Canceled) { // birinchi xato bekor qilish boʻlsa, uni Ctrl+C qilgan
		return errors.New("toʻxtatildi (Ctrl+C)")
	}
	return err
}
