package main

// Vazifa tekshiruvi (TASK.md). Oʻzgartirmang. Har test cheapestPrice'ni synctest pufagi ichida chaqiradi:
// soat soxta (3 soniyalik taklif bir zumda "keladi"), vaqt esa aniq oʻlchanadi (1s, 300ms…).

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"testing/synctest"
	"time"
)

// natija: cheapestPrice nima qaytardi va qancha (soxta) vaqt ketdi.
type natija struct {
	got  int64
	err  error
	took time.Duration
}

// sinash priceOffers'ni offers bilan almashtiradi, cheapestPrice'ni pufak ichida chaqiradi (parentTimeout > 0 boʻlsa,
// chaqiruvchining ctx'i shuncha vaqtdan keyin tugaydi), keyin ortda qolgan goroutine'larni sanaydi.
func sinash(t *testing.T, offers []priceOffer, parentTimeout time.Duration, tekshir func(t *testing.T, n natija)) {
	t.Helper()
	old := priceOffers
	priceOffers = offers
	defer func() { priceOffers = old }()

	synctest.Test(t, func(t *testing.T) {
		ctx := t.Context()
		if parentTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, parentTimeout)
			defer cancel()
		}
		before := runtime.NumGoroutine()
		start := time.Now()
		got, err := cheapestPrice(ctx)
		n := natija{got, err, time.Since(start)}

		// Bekor qilingan goroutine'lar chiqib ketishiga imkon beramiz (soat yurmaydi), keyin sanaymiz.
		synctest.Wait()
		if left := runtime.NumGoroutine() - before; left > 0 {
			t.Errorf("goroutine leak: cheapestPrice qaytganda %d ta goroutine hali kutib turibdi.\n"+
				"\tQaytishdan oldin ularni toʻxtating (ctx'ni askOffer'ga bering, defer cancel()) va natija channel'iga\n"+
				"\tyuborish hech qachon osilib qolmasin (buffer = takliflar soni). Pastdagi «deadlock» shu goroutine'lar haqida.", left)
		}
		tekshir(t, n)
	})
}

func TestVazifaEngArzon(t *testing.T) {
	// Standart takliflar: 100ms → 1 250 000, 200ms → 1 150 000, 3s → 1 100 000 (quoteTimeout 1s ichida kelmaydi).
	sinash(t, priceOffers, 0, func(t *testing.T, n natija) {
		if n.err != nil || n.got != 1_150_000 || n.took != time.Second {
			t.Errorf("cheapestPrice() = %d, %v (%v oʻtdi); kutilgan 1150000, <nil>, aniq 1s.\n"+
				"\t1 soniya ichida ikki taklif keladi, arzoni 1 150 000; 3 soniyalik taklifni kutmang (quoteTimeout).", n.got, n.err, n.took)
		}
	})
}

func TestVazifaHammasiTez(t *testing.T) {
	fast := []priceOffer{
		{100 * time.Millisecond, 1_300_000},
		{200 * time.Millisecond, 1_150_000},
		{300 * time.Millisecond, 1_200_000},
	}
	sinash(t, fast, 0, func(t *testing.T, n natija) {
		if n.err != nil || n.got != 1_150_000 || n.took != 300*time.Millisecond {
			t.Errorf("hamma taklif 300ms ichida: cheapestPrice() = %d, %v (%v oʻtdi); kutilgan 1150000, <nil>, aniq 300ms.\n"+
				"\tHamma javob kelgach, quoteTimeout'ni kutmasdan qayting; eng arzoni ikkinchisi.", n.got, n.err, n.took)
		}
	})
}

func TestVazifaHechKim(t *testing.T) {
	slow := []priceOffer{
		{2 * time.Second, 1_000_000},
		{5 * time.Second, 900_000},
	}
	sinash(t, slow, 0, func(t *testing.T, n natija) {
		if !errors.Is(n.err, errNoOffer) || n.got != 0 || n.took != time.Second {
			t.Errorf("hech kim 1 soniyada javob bermaydi: cheapestPrice() = %d, %v (%v oʻtdi); kutilgan 0, errNoOffer, aniq 1s.", n.got, n.err, n.took)
		}
	})
}

func TestVazifaChaqiruvchi(t *testing.T) {
	// Chaqiruvchining ctx'i 150ms'da tugaydi: birinchi taklif (100ms) kelgan boʻlsa ham, uning xatosi qaytadi.
	sinash(t, priceOffers, 150*time.Millisecond, func(t *testing.T, n natija) {
		if !errors.Is(n.err, context.DeadlineExceeded) || n.got != 0 || n.took != 150*time.Millisecond {
			t.Errorf("chaqiruvchi ctx'i 150ms'da tugaydi: cheapestPrice() = %d, %v (%v oʻtdi); kutilgan 0, context deadline exceeded, aniq 150ms.\n"+
				"\tquoteTimeout'ni chaqiruvchi ctx'idan yasang (WithTimeout(ctx, quoteTimeout)) va u tugaganda ctx.Err()ni qaytaring.", n.got, n.err, n.took)
		}
	})
}
