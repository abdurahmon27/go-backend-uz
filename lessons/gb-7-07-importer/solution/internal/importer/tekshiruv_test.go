package importer

// VAZIFA tekshiruvi (TASK.md): Config.BatchTimeout. Testlar synctest ichida: vaqt soxta, natija har safar bir xil.

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// hangSaver: hang-chi chaqiruv javob bermaydi (baza osilib qolgandek) va faqat ctx bekor boʻlganda qaytadi;
// qolgan chaqiruvlar darhol muvaffaqiyatli.
type hangSaver struct {
	hang  int64
	calls atomic.Int64
}

func (s *hangSaver) SaveBatch(ctx context.Context, rows []Product) error {
	if s.calls.Add(1) == s.hang {
		<-ctx.Done()
		return ctx.Err()
	}
	return nil
}

func TestTaskBatchTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		start := time.Now()
		var st *Stats
		var err error
		done := make(chan struct{})
		go func() {
			defer close(done)
			cfg := Config{Workers: 1, BatchTimeout: 5 * time.Second}
			st, err = Run(ctx, catalog(t, 2_000), &hangSaver{hang: 3}, cfg)
		}()
		select {
		case <-done:
		case <-time.After(time.Minute):
			cancel() // osilgan partiyani boʻshatamiz, aks holda test ham osilib qoladi
			<-done
			t.Fatalf("3-partiya osilib qoldi va Run 1 daqiqada ham qaytmadi: BatchTimeout (5s) ishlatilmayapti (TASK.md, 1-qadam)")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v; kutilgan: context.DeadlineExceeded'ni %%w bilan oʻragan xato (TASK.md, 1-qadam)", err)
		}
		if !strings.HasPrefix(err.Error(), "partiya 3: ") {
			t.Errorf("xato matni %q; «partiya 3: …» bilan boshlanishi kerak (TASK.md, 2-qadam)", err)
		}
		if d := time.Since(start); d != 5*time.Second {
			t.Errorf("Run %v da qaytdi; osilgan partiya aynan BatchTimeout = 5s da toʻxtashi kerak", d)
		}
		if st.Batches.Load() != 2 || st.RolledBatches.Load() != 1 {
			t.Errorf("%v; kutilgan: 2 partiya yozilgan, 1 tasi (3-partiya) rollback (TASK.md, 2-qadam)", st)
		}
	})
}

// TestTaskNoTimeout — BatchTimeout 0: cheklov yoʻq, sekin (10 s) partiyalar ham yoziladi.
func TestTaskNoTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		st, err := Run(context.Background(), catalog(t, 1_000), slowSaver{10 * time.Second}, Config{Workers: 2})
		if err != nil || st.Batches.Load() != 2 {
			t.Fatalf("BatchTimeout 0 bilan: %v, err = %v; kutilgan 2 partiya, xatosiz: 0 — «cheklovsiz» (TASK.md, 3-qadam)", st, err)
		}
	})
}

// TestTaskSlowButInTime — BatchTimeout'dan qisqa partiyalar toʻxtatilmaydi.
func TestTaskSlowButInTime(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cfg := Config{Workers: 2, BatchTimeout: 5 * time.Second}
		st, err := Run(context.Background(), catalog(t, 2_000), slowSaver{4 * time.Second}, cfg)
		if err != nil || st.Batches.Load() != 4 {
			t.Fatalf("partiya 4s, BatchTimeout 5s: %v, err = %v; hammasi yozilishi kerak", st, err)
		}
	})
}
