package store

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	ctx := context.Background()
	s := New()
	if err := s.Seed(ctx); err != nil {
		t.Fatal(err)
	}
	h := []byte("hash-1")
	s.CreateSession(ctx, h, 1, time.Now().Add(time.Hour))
	if u, err := s.UserBySession(ctx, h); err != nil || u.Email != "dilnoza@bozor.test" || u.PasswordHash != "" {
		t.Fatalf("UserBySession: %+v %v", u, err)
	}
	s.DeleteSession(ctx, h)
	if _, err := s.UserBySession(ctx, h); err != ErrNoSession {
		t.Errorf("oʻchirilgan session: err=%v, kutilgan ErrNoSession", err)
	}
	old := []byte("hash-2")
	s.CreateSession(ctx, old, 1, time.Now().Add(-time.Second))
	if _, err := s.UserBySession(ctx, old); err != ErrNoSession {
		t.Errorf("muddati oʻtgan session: err=%v, kutilgan ErrNoSession", err)
	}
}

func TestUsers(t *testing.T) {
	ctx := context.Background()
	s := New()
	if err := s.Seed(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := s.AddUser(ctx, "Dilnoza@Bozor.test", "x"); !errors.Is(err, ErrEmailTaken) {
		t.Errorf("band email (boshqa harflar bilan): err=%v, kutilgan ErrEmailTaken", err)
	}
	if u, err := s.UserByEmail(ctx, "JASUR@bozor.test"); err != nil || u.ID != 2 || u.PasswordHash == "" {
		t.Errorf("UserByEmail: %+v %v", u, err)
	}
	if _, err := s.UserByEmail(ctx, "yoq@bozor.test"); !errors.Is(err, ErrUserNotFound) {
		t.Errorf("yoʻq email: err=%v, kutilgan ErrUserNotFound", err)
	}
	s.CreateSession(ctx, []byte("a"), 1, time.Now().Add(time.Hour))
	s.CreateSession(ctx, []byte("b"), 1, time.Now().Add(time.Hour))
	s.CreateSession(ctx, []byte("c"), 2, time.Now().Add(time.Hour))
	if n, err := s.DeleteUserSessions(ctx, 1); n != 2 || err != nil {
		t.Errorf("DeleteUserSessions(1) = %d, %v; kutilgan 2, nil", n, err)
	}
	if _, err := s.UserBySession(ctx, []byte("c")); err != nil {
		t.Errorf("boshqa foydalanuvchining session'i oʻchmasin: %v", err)
	}
}
