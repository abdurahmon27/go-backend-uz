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
