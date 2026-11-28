// GB-8.03: yumshoq toʻxtash. Vazifa tekshiruvi (TASK.md) — TestWaitBG…: waitBG fon ishlarini kutadi, lekin ctx tugasa
// ctx.Err() bilan darhol qaytadi. Postgres kerak emas: go test ./cmd/api -run WaitBG
// TestServerTimeouts va TestReceiptOwned — videodagi kod (ular boshidan oʻtadi). Bu faylni oʻzgartirmang.
package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/abdurahmon27/bozor/internal/pgtest"
	"github.com/abdurahmon27/bozor/internal/store"
)

// Hamma ish 50 ms'da tugaydi, muddat 2 s: waitBG ish tugaguncha kutib, nil qaytarishi kerak.
func TestWaitBGDone(t *testing.T) {
	var bg sync.WaitGroup
	var finished atomic.Bool
	bg.Go(func() {
		time.Sleep(50 * time.Millisecond)
		finished.Store(true)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := waitBG(ctx, &bg); err != nil {
		t.Fatalf("waitBG = %v; fon ishi 50 ms'da tugadi, muddat 2 s — nil kutilgan", err)
	}
	if !finished.Load() {
		t.Fatal("waitBG fon ishi tugashidan oldin qaytdi: hamma ish tugaguncha kutishi kerak")
	}
}

// Bitta chek osilib qoldi (hech qachon tugamaydi), muddat 100 ms: waitBG muddat tugaganda
// context.DeadlineExceeded bilan qaytishi kerak, chekni kutib osilib qolmasligi kerak.
func TestWaitBGDeadline(t *testing.T) {
	var bg sync.WaitGroup
	release := make(chan struct{})
	bg.Go(func() { <-release }) // osilib qolgan chek
	defer func() {
		close(release) // test oxirida chekni qoʻyib yuboramiz: goroutine test'dan keyin qolmasin
		bg.Wait()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	errc := make(chan error, 1)
	go func() { errc <- waitBG(ctx, &bg) }()
	select {
	case err := <-errc:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("waitBG = %v; muddat tugaganda ctx.Err(), yaʼni context.DeadlineExceeded kutilgan", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("waitBG muddat (100 ms) tugagach 2 s oʻtdi, hali ham chekni kutyapti: " +
			"ctx.Done() yopilganda kutishni toʻxtatib, ctx.Err() qaytarishi kerak (TASK.md)")
	}
}

// Videodagi toʻrt timeout: hech biri nol emas (nol — cheksiz kutish).
func TestServerTimeouts(t *testing.T) {
	srv := newServer(":0", nil, quiet, new(sync.WaitGroup))
	for name, d := range map[string]time.Duration{
		"ReadHeaderTimeout": srv.ReadHeaderTimeout,
		"ReadTimeout":       srv.ReadTimeout,
		"WriteTimeout":      srv.WriteTimeout,
		"IdleTimeout":       srv.IdleTimeout,
	} {
		if d <= 0 {
			t.Errorf("%s = %v: nol — cheksiz kutish, musbat muddat kerak", name, d)
		}
	}
}

// createOrder chekni server egalik qiladigan bg orqali boshlaydi: bg.Wait() chek yuborilguncha kutadi.
// «go sendReceipt(…)» boʻlsa, Wait darhol qaytadi va test yiqiladi. Haqiqiy Postgres kerak.
func TestReceiptOwned(t *testing.T) {
	old := receiptDelay
	receiptDelay = 100 * time.Millisecond
	defer func() { receiptDelay = old }()

	var bg sync.WaitGroup
	h := newHandler(store.New(pgtest.New(t)), quiet, &bg)
	c := login(t, h, jasurLogin)
	rec := send(h, "POST", "/orders", `{"items":[{"product_id":1,"quantity":1}]}`, c.Value)
	if rec.Code != 201 {
		t.Fatalf("POST /orders: status %d, kutilgan 201: %s", rec.Code, rec.Body)
	}
	start := time.Now()
	bg.Wait()
	if d := time.Since(start); d < 50*time.Millisecond {
		t.Errorf("bg.Wait() %v'da qaytdi, chek esa 100 ms ketadi: chek a.bg.Go bilan boshlanishi kerak", d)
	}
}
