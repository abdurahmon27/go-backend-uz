package importer

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

// catalog — n qatorli sinov katalogi (Generate), xotirada.
func catalog(t *testing.T, n int) io.Reader {
	t.Helper()
	var b strings.Builder
	if err := Generate(&b, n); err != nil {
		t.Fatal(err)
	}
	return strings.NewReader(b.String())
}

// slowSaver: har partiya delay davom etadi; ctx bekor boʻlsa, darhol rollback (ctx xatosi).
// synctest ichida vaqt soxta: 10 ms "kutish" testni sekinlashtirmaydi va natija har safar bir xil (GB-7.06).
type slowSaver struct {
	delay time.Duration
}

func (s slowSaver) SaveBatch(ctx context.Context, rows []Product) error {
	select {
	case <-time.After(s.delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func TestImportAll(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		st, err := Run(context.Background(), catalog(t, 2_250), slowSaver{10 * time.Millisecond}, Config{Workers: 4})
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		if st.Rows.Load() != 2_250 || st.Batches.Load() != 5 {
			t.Errorf("%v; kutilgan: 2250 qator, 5 partiya (oxirgisi 250 qator)", st)
		}
	})
}

// TestCancel — Ctrl+C importning oʻrtasida: 4 worker, har partiya 10 ms, bekor qilish 25 ms'da.
// 10 va 20 ms'da 4 tadan partiya commit boʻladi, 25 ms'da yozilayotgan 4 tasi rollback. Run qaytgandan keyin
// synctest bubble'da birorta goroutine qolsa (leak), test yiqiladi.
func TestCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(25 * time.Millisecond)
			cancel() // Ctrl+C
		}()
		st, err := Run(ctx, catalog(t, 10_000), slowSaver{10 * time.Millisecond}, Config{Workers: 4})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v, kutilgan context.Canceled", err)
		}
		if st.Batches.Load() != 8 || st.RolledBatches.Load() != 4 {
			t.Errorf("%v; kutilgan: 8 partiya yozilgan, 4 tasi rollback", st)
		}
	})
}

// TestBadLine — buzuq qator hammasini toʻxtatadi: xato qator raqami bilan, undan keyingi partiyalar yozilmaydi.
func TestBadLine(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		lines := strings.Split(strings.TrimSpace(readAll(t, catalog(t, 10_000))), "\n")
		lines[1234] = "anor;-5;3" // 1235-qator
		r := strings.NewReader(strings.Join(lines, "\n"))
		st, err := Run(context.Background(), r, slowSaver{10 * time.Millisecond}, Config{Workers: 4})
		if err == nil || !strings.Contains(err.Error(), "1235-qator: narx") {
			t.Fatalf("err = %v, kutilgan «1235-qator: narx: …»", err)
		}
		// 1235-qator 3-partiyada: u hech qachon yuborilmaydi, faqat 1- va 2-partiya (1000 qator) boshlanishi mumkin.
		if got := st.Rows.Load() + st.RolledRows.Load(); got > 1_000 {
			t.Errorf("%v; buzuq qatorli partiya yoki undan keyingisi ham boshlandi", st)
		}
	})
}

func TestBadWorkers(t *testing.T) {
	_, err := Run(context.Background(), catalog(t, 10), slowSaver{}, Config{Workers: 0})
	if !errors.Is(err, ErrBadWorkers) {
		t.Fatalf("err = %v, kutilgan ErrBadWorkers", err)
	}
}

func TestParseLine(t *testing.T) {
	p, err := parseLine("oʻrik;35000;5")
	if err != nil || p != (Product{Name: "oʻrik", Price: 3_500_000, Stock: 5}) {
		t.Errorf("parseLine = %+v, %v", p, err)
	}
	for _, bad := range []string{"olma;12000", "olma;0;4", ";100;1", "non;4000;-1", "choy;abc;1"} {
		if _, err := parseLine(bad); err == nil {
			t.Errorf("parseLine(%q): xato kutilgan edi", bad)
		}
	}
}

func readAll(t *testing.T, r io.Reader) string {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
