// api — bozor API'si (GB-6.08'dagi himoyalangan bozor), sozlamasi tashqaridan (GB-8.01), log'lari slog bilan (GB-8.02),
// endi signal kelganda yumshoq toʻxtaydi (GB-8.03, graceful shutdown):
//
//	SIGINT (Ctrl+C) yoki SIGTERM (kill, docker stop, systemctl stop) → yangi ulanish qabul qilinmaydi,
//	ishlayotgan soʻrovlar oxirigacha bajariladi, fon ishlari (chek) kutiladi, keyin Postgres pool'i yopiladi.
//	Hammasi shutdownTimeout ichida; ulgurmasa — xato va 1 kodi.
//
//	docker compose up -d --wait
//	goose up                     (GOOSE_DRIVER, GOOSE_DBSTRING, GOOSE_MIGRATION_DIR: GB-5.06)
//	docker compose exec db psql -f /lab/seed.sql
//	cp .env.example .env         (bir marta; .env gitga tushmaydi)
//	set -a; . ./.env; set +a     (.env'dagi qatorlarni muhitga eksport qiladi)
//	go build ./cmd/api && ./api  → JSON log qatorlari stdout'ga (GB-8.02); Ctrl+C → "toʻxtash boshlandi", "toʻxtadi"
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

// shutdownTimeout — signaldan keyin toʻxtashga beriladigan eng koʻp vaqt: Shutdown va fon ishlari birgalikda.
// docker stop standart holda 10 soniya kutadi, keyin SIGKILL yuboradi (uni ushlab boʻlmaydi): muddat undan qisqa.
const shutdownTimeout = 8 * time.Second

// newServer — bozor API serveri (GB-4.04). Handler — newHandler: middleware zanjiri, ichida newMux (GB-6.01).
// logger — asosiy logger (GB-8.02). bg — handler'lar boshlagan fon ishlari (chek, receipt.go); ularni run kutadi.
// Toʻrt timeout (GB-8.03): har biri bitta ulanish serverni qancha band qila olishini cheklaydi.
func newServer(addr string, st *store.Store, logger *slog.Logger, bg *sync.WaitGroup) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newHandler(st, logger, bg),
		ReadHeaderTimeout: 5 * time.Second,  // soʻrov qatori va sarlavhalar
		ReadTimeout:       10 * time.Second, // butun soʻrov, body bilan
		WriteTimeout:      30 * time.Second, // sarlavhalar oʻqilgandan javob yozilib boʻlguncha
		IdleTimeout:       time.Minute,      // keep-alive ulanish keyingi soʻrovni kutishi
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
// Toʻxtash tartibi defer'lar bilan: run Shutdown va fon ishlari tugagandagina qaytadi, keyin defer pool.Close() —
// Postgres'ni ishlatayotgan soʻrov qolmaganda.
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

	// GB-8.03: SIGINT (Ctrl+C) yoki SIGTERM kelganda sigCtx bekor boʻladi. Signal endi dasturni oʻldirmaydi.
	// Eng boshida: Postgres'ga ulanish paytida kelgan signal ham ulanishni bekor qiladi (ctx sigCtx'dan).
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithTimeout(sigCtx, 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.databaseURL)
	if err != nil {
		return err
	}
	// defer'lar teskari tartibda ishlaydi: pool.Close() eng oxirida, server va fon ishlari tugagandan keyin.
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil { // baza ishlamasa server ishga tushmaydi
		return err
	}
	logger.Debug("Postgres ulandi") // Debug < Info: standart sozlamada yozilmaydi

	// bg — handler'lar boshlagan fon ishlari. bg.Go faqat handler ichida chaqiriladi, Shutdown esa hamma handler
	// tugashini kutadi: waitBG boshlanganda yangi bg.Go boʻlmaydi (sync.WaitGroup: Wait bilan parallel Add mumkin emas).
	var bg sync.WaitGroup
	srv := newServer(cfg.addr(), store.New(pool), logger, &bg)
	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()
	logger.Info("bozor API ishga tushdi", "addr", cfg.addr())

	select {
	case err := <-errc:
		return err // server signalsiz toʻxtadi: masalan, port band
	case <-sigCtx.Done():
	}
	stop() // signal kuzatuvi oʻchdi: ikkinchi Ctrl+C dasturni darhol toʻxtatadi
	logger.Info("toʻxtash boshlandi")

	shCtx, shCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shCancel()
	if err := srv.Shutdown(shCtx); err != nil {
		// Muddat tugadi: qolgan ulanishlarni uzamiz. Ularning r.Context()'i bekor boʻladi, r.Context() bilan ketgan
		// SQL soʻrovlar qaytib, pool ulanishlarini qaytaradi — shunda defer pool.Close() osilib qolmaydi.
		// Fon ishlari bu yoʻlda kutilmaydi: muddat allaqachon tugagan.
		srv.Close()
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := waitBG(shCtx, &bg); err != nil {
		return fmt.Errorf("fon ishlari: %w", err)
	}
	logger.Info("toʻxtadi")
	return nil
}
