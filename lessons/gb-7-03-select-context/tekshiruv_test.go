package main

// Vazifani tekshiradi (TASK.md). Oʻzgartirmang.

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// Tekshiruvda yetkazib beruvchi 3 soniya emas, 100 millisekundda javob beradi.
const delay = 100 * time.Millisecond

func TestMain(m *testing.M) {
	supplierDelay = delay
	os.Exit(m.Run())
}

func TestVazifa1SumStock(t *testing.T) {
	start := time.Now()
	got, err := sumStock(context.Background(), []int64{1, 2, 3})
	if err != nil {
		t.Fatalf("sumStock(ctx, [1 2 3]) xato qaytardi: %v\n\tHar id uchun checkStock(ctx, id) chaqirib, qoldiqlarni qoʻshing (VAZIFA 1).", err)
	}
	if got != 165 {
		t.Errorf("sumStock(ctx, [1 2 3]) = %d, kutilgan 165 (40 + 120 + 5)", got)
	}
	if took := time.Since(start); took < 3*delay {
		t.Errorf("sumStock %v da tugadi: uchala mahsulot checkStock bilan soʻralishi kerak (har biri %v)", took, delay)
	}
}

func TestVazifa1Muddat(t *testing.T) {
	// Muddat ikkinchi mahsulotning oʻrtasida tugaydi: sumStock shu zahoti toʻxtashi kerak.
	ctx, cancel := context.WithTimeout(context.Background(), delay+delay/2)
	defer cancel()
	start := time.Now()
	got, err := sumStock(ctx, []int64{1, 2, 3})
	took := time.Since(start)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("muddati %v boʻlgan ctx bilan sumStock xatosi «%v», kutilgan context.DeadlineExceeded (got = %d).\n"+
			"\tctx'ni checkStock'ka uzating: checkStock(ctx, id), context.Background() emas; xatoni oʻzgartirmay qaytaring.", delay+delay/2, err, got)
	}
	if got != 0 {
		t.Errorf("xato boʻlganda sumStock 0 qaytarishi kerak, hozir %d", got)
	}
	if took > 5*delay/2 {
		t.Errorf("muddat %v da tugadi, lekin sumStock %v ishladi: muddat tugashi bilan toʻxtashi kerak", delay+delay/2, took.Round(time.Millisecond))
	}
}

func TestVazifa1Bekor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // mijoz allaqachon ketgan
	start := time.Now()
	_, err := sumStock(ctx, []int64{1, 2, 3})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("bekor qilingan ctx bilan sumStock xatosi «%v», kutilgan context.Canceled", err)
	}
	if took := time.Since(start); took > delay {
		t.Errorf("ctx oldindan bekor qilingan, lekin sumStock %v kutdi: darhol qaytishi kerak", took.Round(time.Millisecond))
	}
}

func get(t *testing.T, ctx context.Context) (*httptest.ResponseRecorder, time.Duration) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/report", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	start := time.Now()
	newMux().ServeHTTP(rec, req)
	return rec, time.Since(start)
}

func TestVazifa2Report(t *testing.T) {
	orig := reportTimeout
	defer func() { reportTimeout = orig }()
	reportTimeout = 2 * time.Second
	rec, _ := get(t, context.Background())
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "jami: 165" {
		t.Fatalf("GET /report: status %d, tana «%s»; kutilgan 200 va «jami: 165»", rec.Code, strings.TrimSpace(rec.Body.String()))
	}
}

func TestVazifa2Muddat(t *testing.T) {
	orig := reportTimeout
	defer func() { reportTimeout = orig }()
	reportTimeout = delay + delay/2
	rec, took := get(t, context.Background())
	if rec.Code != http.StatusGatewayTimeout {
		t.Errorf("reportTimeout %v, yetkazib beruvchi 3 × %v: status %d, kutilgan 504 Gateway Timeout.\n"+
			"\tctx, cancel := context.WithTimeout(r.Context(), reportTimeout); defer cancel(); "+
			"errors.Is(err, context.DeadlineExceeded) boʻlsa http.StatusGatewayTimeout (VAZIFA 2).", reportTimeout, delay, rec.Code)
	}
	if took > 5*delay/2 {
		t.Errorf("muddat %v, lekin handler %v ishladi: sumStock'ka muddatli ctx berilsin", reportTimeout, took.Round(time.Millisecond))
	}
}

func TestVazifa2MijozKetdi(t *testing.T) {
	orig := reportTimeout
	defer func() { reportTimeout = orig }()
	reportTimeout = 2 * time.Second
	// Mijoz delay/2 dan keyin ulanishni uzadi: server r.Context()'ni bekor qiladi (context.Canceled).
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(delay / 2)
		cancel()
	}()
	rec, took := get(t, ctx)
	if took > 2*delay {
		t.Fatalf("mijoz %v dan keyin ketdi, lekin handler %v ishladi: ctx r.Context()'dan yasalsin, context.Background()'dan emas", delay/2, took.Round(time.Millisecond))
	}
	if rec.Code == http.StatusGatewayTimeout {
		t.Errorf("mijoz ketdi (r.Context() bekor), lekin 504 yozildi: 504 faqat reportTimeout tugaganda (context.DeadlineExceeded)")
	}
}
