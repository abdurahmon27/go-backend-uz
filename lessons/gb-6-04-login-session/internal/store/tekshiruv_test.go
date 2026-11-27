// Vazifa tekshiruvi (TASK.md). Bu faylni oʻzgartirmang.
package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abdurahmon27/bozor/internal/pgtest"
)

func TestDeleteUserSessions(t *testing.T) {
	ctx := context.Background()
	st := New(pgtest.New(t))
	hour := time.Now().Add(time.Hour)
	for _, s := range []struct {
		hash string
		user int64
	}{{"dilnoza-telefon", 1}, {"dilnoza-noutbuk", 1}, {"jasur", 2}} {
		if err := st.CreateSession(ctx, []byte(s.hash), s.user, hour); err != nil {
			t.Fatal(err)
		}
	}
	n, err := st.DeleteUserSessions(ctx, 1)
	if err != nil {
		t.Fatalf("DeleteUserSessions(1): %v (TASK.md, 1-qadam)", err)
	}
	if n != 2 {
		t.Errorf("DeleteUserSessions(1) = %d, kutilgan 2: Dilnozaning ikki session'i (TASK.md, 1-qadam)", n)
	}
	if _, err := st.UserBySession(ctx, []byte("dilnoza-noutbuk")); !errors.Is(err, ErrNoSession) {
		t.Errorf("Dilnozaning session'i qoldi: err = %v (TASK.md, 1-qadam)", err)
	}
	if _, err := st.UserBySession(ctx, []byte("jasur")); err != nil {
		t.Errorf("Jasurning session'i ham oʻchdi: %v (TASK.md, 1-qadam: WHERE user_id = $1)", err)
	}
}

// 🟡 ixtiyoriy (TASK.md, 3-qadam): muddati oʻtgan qatorlarni tozalash.
func TestDeleteExpiredSessions(t *testing.T) {
	ctx := context.Background()
	st := New(pgtest.New(t))
	if err := st.CreateSession(ctx, []byte("eski"), 1, time.Now().Add(-time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := st.CreateSession(ctx, []byte("tirik"), 1, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	n, err := st.DeleteExpiredSessions(ctx)
	if errors.Is(err, errOptional) {
		t.Skip("🟡 ixtiyoriy: DeleteExpiredSessions hali yozilmagan (TASK.md, 3-qadam)")
	}
	if err != nil || n != 1 {
		t.Errorf("DeleteExpiredSessions = %d, %v; kutilgan 1 (faqat «eski»)", n, err)
	}
	if _, err := st.UserBySession(ctx, []byte("tirik")); err != nil {
		t.Errorf("tirik session oʻchib ketdi: %v", err)
	}
}
