// Package importer — katta katalog faylini bazaga yozadigan bosqichlar (GB-7.07, 7-modul nazorat nuqtasi).
//
//	fayl → read (1 goroutine) → batches channel → work (Workers ta goroutine) → Saver (Postgres)
//
// read qatorlarni oʻqiydi, tekshiradi va BatchSize'lik partiyalarga yigʻadi. Har work goroutine partiyani oladi
// va Saver'ga beradi: Postgres'da bitta partiya = bitta transaction (pg.go), yoki toʻliq yoziladi, yoki umuman.
// Hammasi bitta errgroup'da (GB-7.04): birinchi xato (buzuq qator, baza xatosi) yoki tashqi ctx'ning bekor
// boʻlishi (Ctrl+C) qolgan hamma goroutine'ni toʻxtatadi, Run esa ular tugagandan keyingina qaytadi.
package importer

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

// BatchSize — bitta partiyadagi qatorlar soni (bitta transaction).
const BatchSize = 500

// Batch — fayldagi ketma-ket BatchSize ta qator. N — partiya raqami, 1 dan.
type Batch struct {
	N    int
	Rows []Product
}

// Saver partiyani bazaga yozadi: yoki hammasini, yoki hech birini (xato qaytaradi).
// Haqiqiysi — PG (pg.go); testlarda soxtasi (vaqt va xatoni test boshqaradi).
type Saver interface {
	SaveBatch(ctx context.Context, rows []Product) error
}

// Config — import sozlamalari.
type Config struct {
	Workers int // bir vaqtda nechta partiya yoziladi (kamida 1)

	// BatchTimeout — bitta partiyaga beriladigan eng koʻp vaqt; 0 — cheklovsiz (vazifa, TASK.md).
	BatchTimeout time.Duration
}

// Stats — import hisobi. Uni Workers ta goroutine bir vaqtda oshiradi, shuning uchun atomic (GB-7.05).
// Rows — commit boʻlgan partiyalardagi qatorlar (qoʻshilgan yoki upsert bilan yangilangan); faylda nomlar
// takrorlanmasa va jadval boʻsh boʻlsa, import tugagach Rows = SELECT count(*). Rolled* — SaveBatch xato qaytargan
// partiyalar (Begin'ning xatosi ham shu yerda: unda hech narsa yozilmagan). Commit javobi tarmoqda yoʻqolsa,
// partiya bazada boʻlsa ham Rolled'ga tushadi: bu holatni upsert qayta ishga tushirish bilan tuzatadi.
type Stats struct {
	Rows, Batches             atomic.Int64 // commit boʻlgan qatorlar va partiyalar
	RolledRows, RolledBatches atomic.Int64 // xato bilan tugagan (rollback) qatorlar va partiyalar
}

func (s *Stats) String() string {
	return fmt.Sprintf("yozildi: %d qator (%d partiya) · rollback: %d qator (%d partiya)",
		s.Rows.Load(), s.Batches.Load(), s.RolledRows.Load(), s.RolledBatches.Load())
}

// ErrBadWorkers — Workers 1 dan kichik: hech kim partiya olmasdi, read esa abadiy kutardi.
var ErrBadWorkers = errors.New("workers kamida 1 boʻlishi kerak")

// Run r'dagi katalogni s'ga cfg.Workers ta goroutine bilan yozadi. Birinchi xatoni qaytaradi; ctx bekor
// qilinsa (Ctrl+C), ctx xatosini. Qaytganda hech bir goroutine ishlamayotgan boʻladi, Stats esa yakuniy.
func Run(ctx context.Context, r io.Reader, s Saver, cfg Config) (*Stats, error) {
	st := &Stats{}
	if cfg.Workers < 1 {
		return st, ErrBadWorkers
	}
	batches := make(chan Batch)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(batches)
		return read(gctx, r, batches)
	})
	for range cfg.Workers {
		g.Go(func() error {
			return work(gctx, s, batches, st, cfg.BatchTimeout)
		})
	}
	return st, g.Wait()
}

// read qatorlarni oʻqiydi va partiyalarni out'ga yuboradi. Buzuq qatorda raqami bilan xato qaytaradi.
// Boʻsh qatorlar va # bilan boshlangan izohlar tashlab oʻtiladi (GB-3.01).
func read(ctx context.Context, r io.Reader, out chan<- Batch) error {
	sc := bufio.NewScanner(r)
	b := Batch{N: 1}
	line := 0
	for sc.Scan() {
		line++
		text := strings.TrimSpace(sc.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		p, err := parseLine(text)
		if err != nil {
			return fmt.Errorf("%d-qator: %w", line, err)
		}
		b.Rows = append(b.Rows, p)
		if len(b.Rows) == BatchSize {
			if err := send(ctx, out, b); err != nil {
				return err
			}
			b = Batch{N: b.N + 1}
		}
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("%d-qator: %w", line+1, err)
	}
	if len(b.Rows) > 0 {
		return send(ctx, out, b)
	}
	return nil
}

// send partiyani yuboradi yoki ctx bekor boʻlsa, kutmasdan qaytadi (GB-7.03).
func send(ctx context.Context, out chan<- Batch, b Batch) error {
	select {
	case out <- b:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// work partiyalarni in channel'i yopilguncha oladi va s'ga yozadi.
// timeout > 0 boʻlsa, har partiya shu vaqtdan koʻp davom etmaydi (save).
func work(ctx context.Context, s Saver, in <-chan Batch, st *Stats, timeout time.Duration) error {
	for b := range in {
		if ctx.Err() != nil {
			return ctx.Err() // bekor qilingan: yangi partiya boshlanmaydi
		}
		if err := save(ctx, s, b.Rows, timeout); err != nil {
			st.RolledRows.Add(int64(len(b.Rows)))
			st.RolledBatches.Add(1)
			return fmt.Errorf("partiya %d: %w", b.N, err)
		}
		st.Rows.Add(int64(len(b.Rows)))
		st.Batches.Add(1)
	}
	return nil
}

// save bitta partiyani yozadi. timeout > 0 boʻlsa, partiya oʻz muddatli ctx'ini oladi (GB-7.03): muddat tugasa,
// Exec toʻxtaydi, transaction rollback boʻladi, xato esa context.DeadlineExceeded'ni oʻz ichida saqlaydi.
// cancel defer bilan: muddatli ctx'ning taymeri partiya tugashi bilan boʻshaydi (go vet: lostcancel).
func save(ctx context.Context, s Saver, rows []Product, timeout time.Duration) error {
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	return s.SaveBatch(ctx, rows)
}
