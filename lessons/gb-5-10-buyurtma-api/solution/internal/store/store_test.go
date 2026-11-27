// Store testlari haqiqiy Postgres'da (mock yoʻq). Postgres uchun: docker compose up -d --wait (shu papkada).
// Har test oʻz sxemasida (pgtest.New): migratsiyalar + seed.sql. Boshlangʻich holat:
// olma(1) 40, non(2) 3, choy(3) 8, oʻrik(4) 2; 3 ta buyurtma, 4 ta buyurtma qatori.
package store_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/abdurahmon27/bozor/internal/pgtest"
	"github.com/abdurahmon27/bozor/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// count — jadvaldagi (yoki shartga mos) qatorlar soni.
func count(t *testing.T, pool *pgxpool.Pool, from string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(testCtx(t), "SELECT count(*) FROM "+from).Scan(&n); err != nil {
		t.Fatalf("%s sanalmadi: %v", from, err)
	}
	return n
}

func stockOf(t *testing.T, pool *pgxpool.Pool, id int64) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(testCtx(t), "SELECT stock FROM products WHERE id = $1", id).Scan(&n); err != nil {
		t.Fatalf("mahsulot %d qoldigʻi oʻqilmadi: %v", id, err)
	}
	return n
}

func TestPlaceOrder(t *testing.T) {
	pool := pgtest.New(t)
	st := store.New(pool)
	o, err := st.PlaceOrder(testCtx(t), "Aziza", []store.OrderItem{
		{ProductID: 1, Quantity: 2}, {ProductID: 4, Quantity: 1},
	})
	if err != nil {
		t.Fatalf("PlaceOrder: %v", err)
	}
	if o.ID != 4 || o.Status != "yangi" || o.Total != 2*1200000+3500000 || len(o.Items) != 2 {
		t.Errorf("buyurtma = %+v; kutilgan: id 4, status yangi, 2 qator, jami 5900000", o)
	}
	if got := count(t, pool, "order_items WHERE order_id = 4"); got != 2 {
		t.Errorf("4-buyurtmada %d qator, kutilgan 2", got)
	}
	if got := stockOf(t, pool, 1); got != 38 {
		t.Errorf("olma qoldigʻi %d, kutilgan 38", got)
	}
}

// Ikkinchi qatorda qoldiq yetmaydi: birinchi qatorning yozuvi ham bekor boʻlishi kerak.
func TestPlaceOrderRollsBack(t *testing.T) {
	pool := pgtest.New(t)
	_, err := store.New(pool).PlaceOrder(testCtx(t), "Jasur", []store.OrderItem{
		{ProductID: 1, Quantity: 2}, {ProductID: 4, Quantity: 5},
	})
	var se *store.StockError
	if !errors.As(err, &se) {
		t.Fatalf("PlaceOrder(oʻrik x5) = %v; kutilgan *StockError", err)
	}
	if got := count(t, pool, "orders"); got != 3 {
		t.Errorf("orders'da %d qator, kutilgan 3: yarim buyurtma qoldi", got)
	}
	if got := count(t, pool, "order_items"); got != 4 {
		t.Errorf("order_items'da %d qator, kutilgan 4", got)
	}
	if got := stockOf(t, pool, 1); got != 40 {
		t.Errorf("olma qoldigʻi %d, kutilgan 40: birinchi qator bekor boʻlmadi", got)
	}
}

func TestPlaceOrderUnknownProduct(t *testing.T) {
	pool := pgtest.New(t)
	_, err := store.New(pool).PlaceOrder(testCtx(t), "Aziza", []store.OrderItem{{ProductID: 99, Quantity: 1}})
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("PlaceOrder(99) = %v; kutilgan ErrNotFound", err)
	}
	if got := count(t, pool, "orders"); got != 3 {
		t.Errorf("orders'da %d qator, kutilgan 3", got)
	}
}

// 20 xaridor bir vaqtda bittadan choy oladi, omborda 8 ta: aynan 8 tasi oʻtadi, qoldiq 0, manfiy emas.
func TestConcurrentOrders(t *testing.T) {
	pool := pgtest.New(t)
	st := store.New(pool)
	ctx := testCtx(t)
	errs := make([]error, 20)
	var wg sync.WaitGroup
	for i := range errs {
		wg.Go(func() {
			_, errs[i] = st.PlaceOrder(ctx, "xaridor", []store.OrderItem{{ProductID: 3, Quantity: 1}})
		})
	}
	wg.Wait()
	ok := 0
	for _, err := range errs {
		var se *store.StockError
		switch {
		case err == nil:
			ok++
		case !errors.As(err, &se):
			t.Errorf("kutilmagan xato: %v", err)
		}
	}
	if ok != 8 {
		t.Errorf("%d ta buyurtma oʻtdi, kutilgan 8 (omborda 8 ta choy)", ok)
	}
	if got := stockOf(t, pool, 3); got != 0 {
		t.Errorf("choy qoldigʻi %d, kutilgan 0", got)
	}
	if got := count(t, pool, "orders"); got != 3+8 {
		t.Errorf("orders'da %d qator, kutilgan 11", got)
	}
}

// Ikki xil tartibdagi savatlar ([olma, choy] va [choy, olma]) bir vaqtda: PlaceOrder qatorlarni ID tartibida
// qulflagani uchun bir-birini kutib qotib qolmaydi (aks holda Postgres 40P01 deadlock detected beradi).
func TestOppositeCarts(t *testing.T) {
	pool := pgtest.New(t)
	st := store.New(pool)
	ctx := testCtx(t)
	errs := make([]error, 16)
	var wg sync.WaitGroup
	for i := range errs {
		items := []store.OrderItem{{ProductID: 1, Quantity: 1}, {ProductID: 3, Quantity: 1}}
		if i%2 == 1 {
			items = []store.OrderItem{{ProductID: 3, Quantity: 1}, {ProductID: 1, Quantity: 1}}
		}
		wg.Go(func() {
			_, errs[i] = st.PlaceOrder(ctx, "xaridor", items)
		})
	}
	wg.Wait()
	ok := 0
	for _, err := range errs {
		var se *store.StockError
		switch {
		case err == nil:
			ok++
		case !errors.As(err, &se):
			t.Errorf("kutilmagan xato (deadlock?): %v", err)
		}
	}
	if ok != 8 {
		t.Errorf("%d ta buyurtma oʻtdi, kutilgan 8 (choy 8 ta)", ok)
	}
	if got := stockOf(t, pool, 1); got != 32 {
		t.Errorf("olma qoldigʻi %d, kutilgan 32 (40 − 8)", got)
	}
}

func TestEmptyCart(t *testing.T) {
	pool := pgtest.New(t)
	if _, err := store.New(pool).PlaceOrder(testCtx(t), "Aziza", nil); !errors.Is(err, store.ErrBadQuantity) {
		t.Fatalf("PlaceOrder(boʻsh savat) = %v; kutilgan ErrBadQuantity", err)
	}
	if got := count(t, pool, "orders"); got != 3 {
		t.Errorf("orders'da %d qator, kutilgan 3", got)
	}
}

func TestAddDuplicate(t *testing.T) {
	st := store.New(pgtest.New(t))
	if _, err := st.Add(testCtx(t), "olma", 100, 1); !errors.Is(err, store.ErrDuplicate) {
		t.Fatalf("Add(olma) = %v; kutilgan ErrDuplicate", err)
	}
}

func TestDelete(t *testing.T) {
	st := store.New(pgtest.New(t))
	p, err := st.Add(testCtx(t), "anor", 2000000, 5)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := st.Delete(testCtx(t), p.ID); err != nil {
		t.Errorf("Delete(%d) = %v; kutilgan nil", p.ID, err)
	}
	if err := st.Delete(testCtx(t), p.ID); !errors.Is(err, store.ErrNotFound) {
		t.Errorf("ikkinchi Delete = %v; kutilgan ErrNotFound", err)
	}
	if err := st.Delete(testCtx(t), 1); !errors.Is(err, store.ErrInUse) {
		t.Errorf("Delete(olma) = %v; kutilgan ErrInUse: olma buyurtmalarda bor", err)
	}
}
