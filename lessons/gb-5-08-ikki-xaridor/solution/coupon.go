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
// Tekshiruv va kamaytirish bitta UPDATE buyrugʻida (WHERE uses_left > 0): parallel chaqiruvlarda ham roppa-rosa uses_left marta oʻtadi.
func (s *Store) UseCoupon(ctx context.Context, code string) error {
	tag, err := s.pool.Exec(ctx,
		"UPDATE coupons SET uses_left = uses_left - 1 WHERE code = $1 AND uses_left > 0", code)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 1 {
		return nil
	}
	// Nol qator: promo-kod yoʻq yoki tugagan. Qaysi biri ekanini alohida oʻqib bilamiz.
	var left int
	err = s.pool.QueryRow(ctx,
		"SELECT uses_left FROM coupons WHERE code = $1", code).Scan(&left)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrCouponNotFound
	}
	if err != nil {
		return err
	}
	return ErrCouponUsedUp
}

var errUseCouponLockedTODO = errors.New("UseCouponLocked hali yozilmagan")

// UseCouponLocked ham promo-kodni ishlatadi, lekin boshqacha: qatorni FOR UPDATE bilan qulflaydi, Go'da tekshiradi, keyin yozadi.
func (s *Store) UseCouponLocked(ctx context.Context, code string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var left int
	err = tx.QueryRow(ctx,
		"SELECT uses_left FROM coupons WHERE code = $1 FOR UPDATE", code).Scan(&left)
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
