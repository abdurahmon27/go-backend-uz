package importer

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/abdurahmon27/bozor/internal/pgtest"
	"github.com/jackc/pgx/v5/pgxpool"
)

func countProducts(t *testing.T, pool *pgxpool.Pool) int64 {
	t.Helper()
	var n int64
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM products").Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

// TestPostgresImport — haqiqiy Postgres: hamma qator yoziladi; ikkinchi import (upsert) qator qoʻshmaydi.
func TestPostgresImport(t *testing.T) {
	pool := pgtest.New(t)
	for i := range 2 {
		st, err := Run(context.Background(), catalog(t, 2_000), NewPG(pool), Config{Workers: 4})
		if err != nil {
			t.Fatalf("%d-import: %v", i+1, err)
		}
		if st.Rows.Load() != 2_000 {
			t.Errorf("%d-import: %v; kutilgan 2000 qator", i+1, st)
		}
		if n := countProducts(t, pool); n != 2_000 {
			t.Errorf("%d-importdan keyin bazada %d qator, kutilgan 2000", i+1, n)
		}
	}
}

// cancelAfter — PG'ni oʻrab, after-chi partiya commit boʻlgach ctx'ni bekor qiladi (Ctrl+C oʻrnida).
// Shu paytda boshqa worker'lar oʻz partiyasini yozayotgan boʻladi.
type cancelAfter struct {
	pg     *PG
	after  int64
	done   atomic.Int64
	cancel context.CancelFunc
}

func (c *cancelAfter) SaveBatch(ctx context.Context, rows []Product) error {
	err := c.pg.SaveBatch(ctx, rows)
	if err == nil && c.done.Add(1) == c.after {
		c.cancel()
	}
	return err
}

// TestPostgresCancel — bekor qilingan importdan keyin bazadagi qatorlar soni hisobotdagi «yozildi» bilan teng:
// yarim partiya qolmaydi.
func TestPostgresCancel(t *testing.T) {
	pool := pgtest.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := &cancelAfter{pg: NewPG(pool), after: 3, cancel: cancel}
	st, err := Run(ctx, catalog(t, 20_000), s, Config{Workers: 4})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, kutilgan context.Canceled", err)
	}
	n := countProducts(t, pool)
	t.Logf("%v · bazada: %d qator", st, n)
	if n != st.Rows.Load() || n%BatchSize != 0 {
		t.Errorf("bazada %d qator, hisobot: %v; teng va %d ga karrali boʻlishi kerak", n, st, BatchSize)
	}
}
