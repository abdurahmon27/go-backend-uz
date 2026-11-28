// Store testlari: haqiqiy Postgres, har test oʻz sxemasida (pgtest.New), seed.sql'dagi ikki foydalanuvchi.
package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abdurahmon27/bozor/internal/pgtest"
)

func TestUserByEmail(t *testing.T) {
	st := New(pgtest.New(t))
	u, err := st.UserByEmail(context.Background(), "dilnoza@bozor.test")
	if err != nil || u.ID != 1 || u.PasswordHash == "" {
		t.Fatalf("UserByEmail(dilnoza) = %v, %v", u, err)
	}
	if _, err := st.UserByEmail(context.Background(), "yoq@bozor.test"); !errors.Is(err, ErrUserNotFound) {
		t.Errorf("yoʻq email: err = %v, kutilgan ErrUserNotFound", err)
	}
}

func TestSessionLifecycle(t *testing.T) {
	ctx := context.Background()
	st := New(pgtest.New(t))
	hash := []byte("32-baytli-hash-oʻrnida-sinov-uchun")
	if err := st.CreateSession(ctx, hash, 1, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	u, err := st.UserBySession(ctx, hash)
	if err != nil || u.Email != "dilnoza@bozor.test" {
		t.Fatalf("UserBySession = %v, %v; kutilgan dilnoza", u, err)
	}
	if err := st.DeleteSession(ctx, hash); err != nil {
		t.Fatal(err)
	}
	if _, err := st.UserBySession(ctx, hash); !errors.Is(err, ErrNoSession) {
		t.Errorf("oʻchirilgan session: err = %v, kutilgan ErrNoSession", err)
	}
	old := []byte("eski")
	if err := st.CreateSession(ctx, old, 1, time.Now().Add(-time.Minute)); err != nil {
		t.Fatal(err)
	}
	if _, err := st.UserBySession(ctx, old); !errors.Is(err, ErrNoSession) {
		t.Errorf("muddati oʻtgan session: err = %v, kutilgan ErrNoSession", err)
	}
}

func TestUserByEmailCase(t *testing.T) {
	st := New(pgtest.New(t))
	u, err := st.UserByEmail(context.Background(), "Dilnoza@Bozor.TEST")
	if err != nil || u.ID != 1 {
		t.Errorf("katta harfli email: %v, %v; kutilgan dilnoza (ID 1)", u, err)
	}
}

// GB-6.06: Order va UpdateNote egasini SQL shartida tekshiradi. seed.sql: 41 — Dilnozaniki (1), 42 — Jasurniki (2).
func TestOrderOwner(t *testing.T) {
	ctx := context.Background()
	st := New(pgtest.New(t))
	o, err := st.Order(ctx, 41, 1)
	if err != nil || o.UserID != 1 || len(o.Items) != 2 || o.Total != 3600000 || o.Note == "" {
		t.Fatalf("Order(41, dilnoza) = %+v, %v; kutilgan Dilnozaning buyurtmasi, 2 qator, 3600000 tiyin", o, err)
	}
	for _, c := range []struct{ id, user int64 }{{42, 1}, {999, 1}, {41, 2}} {
		if _, err := st.Order(ctx, c.id, c.user); !errors.Is(err, ErrOrderNotFound) {
			t.Errorf("Order(%d, foydalanuvchi %d): err = %v, kutilgan ErrOrderNotFound", c.id, c.user, err)
		}
	}
}

func TestUsersRole(t *testing.T) {
	st := New(pgtest.New(t))
	users, err := st.Users(context.Background())
	if err != nil || len(users) != 3 || users[2].Role != RoleAdmin || users[0].Role != RoleCustomer {
		t.Fatalf("Users() = %+v, %v; kutilgan 3 ta, uchinchisi admin", users, err)
	}
	for _, u := range users {
		if u.PasswordHash != "" {
			t.Errorf("Users(): %s uchun parol hash'i qaytdi", u.Email)
		}
	}
}

func TestPlaceOrderOwner(t *testing.T) {
	ctx := context.Background()
	st := New(pgtest.New(t))
	o, err := st.PlaceOrder(ctx, 2, "", []OrderItem{{ProductID: 3, Quantity: 2}})
	if err != nil || o.ID != 43 || o.UserID != 2 || o.Total != 5000000 {
		t.Fatalf("PlaceOrder(jasur) = %+v, %v; kutilgan 43-buyurtma, egasi 2", o, err)
	}
	if _, err := st.Order(ctx, o.ID, 1); !errors.Is(err, ErrOrderNotFound) {
		t.Errorf("Jasurning yangi buyurtmasi Dilnozaga koʻrindi: err = %v", err)
	}
	var se *StockError
	if _, err := st.PlaceOrder(ctx, 1, "", []OrderItem{{ProductID: 4, Quantity: 5}}); !errors.As(err, &se) {
		t.Errorf("oʻrik 5 dona: err = %v, kutilgan *StockError", err)
	}
}
