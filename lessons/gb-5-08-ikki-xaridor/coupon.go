package bozor

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// ErrCouponUsedUp — promo-kodning ishlatish soni tugagan.
var ErrCouponUsedUp = errors.New("promo-kod tugagan")

// ErrCouponNotFound — bunday promo-kod yoʻq.
var ErrCouponNotFound = errors.New("promo-kod topilmadi")

// UseCoupon promo-kodni bir marta ishlatadi: coupons.uses_left bittaga kamayadi.
// VAZIFA (TASK.md): bu versiya parallel chaqiruvlarda xato ishlaydi. Uni tuzating.
func (s *Store) UseCoupon(ctx context.Context, code string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var left int
	err = tx.QueryRow(ctx,
		"SELECT uses_left FROM coupons WHERE code = $1", code).Scan(&left)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrCouponNotFound
	}
	if err != nil {
		return err
	}
	if left == 0 {
		return ErrCouponUsedUp
	}
	_, err = tx.Exec(ctx,
		"UPDATE coupons SET uses_left = $1 WHERE code = $2", left-1, code)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

var errUseCouponLockedTODO = errors.New("UseCouponLocked hali yozilmagan")

// UseCouponLocked — IXTIYORIY vazifa (TASK.md, 2-qadam): UseCoupon bilan bir xil natija, lekin SELECT … FOR UPDATE bilan.
func (s *Store) UseCouponLocked(ctx context.Context, code string) error {
	return errUseCouponLockedTODO
}
