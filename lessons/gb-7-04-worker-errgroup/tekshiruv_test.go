package main

// Vazifani tekshiradi (TASK.md). Oʻzgartirmang.
//
// Har chaqiruv vaqt chegarasi bilan ishlaydi: SetLimit(0) bilan g.Go hech qachon qaytmaydi. Oddiy dasturda runtime
// buni «all goroutines are asleep - deadlock!» bilan toʻxtatadi, testda esa time.After taymeri tirik, shuning uchun
// test osilib qolardi. Shu sababli 3 soniyadan keyin oʻzbekcha xabar bilan yiqiladi (select va time.After: GB-7.03).

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type result struct {
	prices []int64
	err    error
}

// run fetchPrices'ni chaqiradi; 3 soniyada tugamasa, ok = false.
func run(ids []int64, limit int) (result, bool) {
	resetStats()
	done := make(chan result, 1)
	go func() {
		p, err := fetchPrices(context.Background(), ids, limit)
		done <- result{p, err}
	}()
	select {
	case r := <-done:
		return r, true
	case <-time.After(3 * time.Second):
		return result{}, false
	}
}

func TestVazifaLimit(t *testing.T) {
	for _, limit := range []int{5, 2} {
		r, ok := run(catalogIDs(200), limit)
		if !ok {
			t.Fatalf("fetchPrices(…, %d) 3 soniyada tugamadi", limit)
		}
		if stats.peak > limit {
			t.Fatalf("limit %d edi, lekin yetkazib beruvchiga bir vaqtda %d ta soʻrov ketdi.\n"+
				"\tHar ID uchun alohida goroutine ochilyapti: errgroup'ning SetLimit(limit)'i kerak.", limit, stats.peak)
		}
		if r.err != nil {
			t.Fatalf("200 ta toʻgʻri ID bilan fetchPrices xato qaytardi: %s", firstLine(r.err))
		}
		if len(r.prices) != 200 {
			t.Fatalf("200 ta narx kutilgan, keldi %d ta", len(r.prices))
		}
		for i, id := range catalogIDs(200) {
			if want := id * 1000 * tiyinPerSom; r.prices[i] != want {
				t.Fatalf("prices[%d] = %d, kutilgan %d (ID %d narxi): natija ids tartibida, prices[i]ga yozilsin",
					i, r.prices[i], want, id)
			}
		}
	}
}

func TestVazifaBirinchiXato(t *testing.T) {
	ids := catalogIDs(200)
	ids[49] = 9999 // yetkazib beruvchi bilmaydigan ID
	r, ok := run(ids, 5)
	if !ok {
		t.Fatal("fetchPrices 3 soniyada tugamadi")
	}
	if r.err == nil {
		t.Fatal("ID 9999 boʻlsa fetchPrices xato qaytarishi kerak, qaytardi nil")
	}
	if stats.calls >= 100 {
		t.Fatalf("xatodan keyin ham soʻrovlar ketaverdi: jami %d ta (200 tadan), qolganlari bekor qilinmadi.\n"+
			"\terrgroup.WithContext(ctx) bergan ctx'ni g.Go ichida checkPrice'ga bering.", stats.calls)
	}
	if !errors.Is(r.err, ErrNotFound) {
		t.Fatalf("xato ErrNotFound'ni oʻrab olishi kerak (%%w), keldi: %s", firstLine(r.err))
	}
	if !strings.HasPrefix(r.err.Error(), "narx 9999: ") {
		t.Fatalf("xato «narx 9999: mahsulot topilmadi» boʻlishi kerak, keldi «%v»", firstLine(r.err))
	}
	if r.prices != nil {
		t.Errorf("xato boʻlsa narxlar nil boʻlsin, keldi %d ta", len(r.prices))
	}
}

func TestVazifaNolLimit(t *testing.T) {
	for _, limit := range []int{0, -1} {
		r, ok := run(catalogIDs(10), limit)
		if !ok {
			t.Fatalf("fetchPrices(…, %d) 3 soniyada tugamadi: SetLimit(0) bilan g.Go hech qachon qaytmaydi.\n"+
				"\tlimit 1 dan kichik boʻlsa, errgroup'ga yetmasdan errBadLimit qaytaring.", limit)
		}
		if !errors.Is(r.err, errBadLimit) {
			t.Fatalf("fetchPrices(…, %d): errBadLimit kutilgan, keldi %v", limit, r.err)
		}
	}
}
