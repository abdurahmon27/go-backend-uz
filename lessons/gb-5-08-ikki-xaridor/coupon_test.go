// Promo-kod vazifasini tekshiradi (TASK.md). Uni oʻzgartirmang. Yordamchilar store_test.go da.
package bozor

import (
	"context"
	"errors"
	"sync"
	"testing"
)

// BAHOR promo-kodi bozor.sql da 3 marta ishlatiladi.
func TestUseCouponOneByOne(t *testing.T) {
	s := newStore(t)
	for i := 1; i <= 3; i++ {
		if err := s.UseCoupon(testCtx(t), "BAHOR"); err != nil {
			t.Fatalf("%d-marta UseCoupon(BAHOR) = %v; kutilgan: nil", i, err)
		}
	}
	if err := s.UseCoupon(testCtx(t), "BAHOR"); !errors.Is(err, ErrCouponUsedUp) {
		t.Errorf("4-marta UseCoupon(BAHOR) = %v; kutilgan: ErrCouponUsedUp", err)
	}
	if err := s.UseCoupon(testCtx(t), "YOQ"); !errors.Is(err, ErrCouponNotFound) {
		t.Errorf("UseCoupon(YOQ) = %v; kutilgan: ErrCouponNotFound", err)
	}
}

// 10 kishi bir vaqtda BAHOR'ni ishlatadi: roppa-rosa 3 tasiga oʻtishi kerak. 5 marta takrorlanadi.
func checkCouponConcurrent(t *testing.T, s *Store, use func(context.Context, string) error) {
	t.Helper()
	ctx := testCtx(t)
	for round := 1; round <= 5; round++ {
		if _, err := s.pool.Exec(testCtx(t), "UPDATE coupons SET uses_left = 3 WHERE code = 'BAHOR'"); err != nil {
			t.Fatalf("promo-kod tiklanmadi: %v", err)
		}
		errs := make([]error, 10)
		var wg sync.WaitGroup
		for i := range errs {
			wg.Go(func() {
				errs[i] = use(ctx, "BAHOR")
			})
		}
		wg.Wait()
		ok := 0
		for _, err := range errs {
			switch {
			case err == nil:
				ok++
			case errors.Is(err, ErrCouponUsedUp): // tugagan promo-kod: kutilgan holat
			default:
				t.Fatalf("kutilmagan xato: %v", err)
			}
		}
		if ok != 3 {
			t.Fatalf("%d-urinish: promo-kod 3 marta, 10 kishi bir vaqtda: %d tasiga oʻtdi; kutilgan: 3 (TASK.md)", round, ok)
		}
	}
}

func TestUseCouponConcurrent(t *testing.T) {
	s := newStore(t)
	checkCouponConcurrent(t, s, s.UseCoupon)
}

// Ixtiyoriy vazifa (TASK.md, 2-qadam): UseCouponLocked yozilmagan boʻlsa test oʻtkazib yuboriladi.
func TestUseCouponLockedOptional(t *testing.T) {
	s := newStore(t)
	if err := s.UseCouponLocked(testCtx(t), "YOQ"); errors.Is(err, errUseCouponLockedTODO) {
		t.Skip("ixtiyoriy vazifa: UseCouponLocked hali yozilmagan (TASK.md, 2-qadam)")
	}
	checkCouponConcurrent(t, s, s.UseCouponLocked)
}
