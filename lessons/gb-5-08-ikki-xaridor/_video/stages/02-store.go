package bozor

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound — jadvalda bunday mahsulot yoʻq (GB-2.08, GB-5.05).
var ErrNotFound = errors.New("mahsulot topilmadi")

// ErrNotEnoughStock — omborda yetarli qoldiq yoʻq (GB-5.07).
var ErrNotEnoughStock = errors.New("qoldiq yetarli emas")

// ErrInvalidQuantity — miqdor noldan katta boʻlishi kerak (GB-5.07: manfiy miqdor qoldiqni oshirib yuborar edi).
var ErrInvalidQuantity = errors.New("miqdor noldan katta boʻlishi kerak")

// Store — Postgres'dagi bozor jadvallariga Go'dan kirish (GB-5.05).
type Store struct {
	pool *pgxpool.Pool
}

// Buy id mahsulotdan qty dona sotadi. Tekshiruv va kamaytirish bitta UPDATE buyrugʻida:
// "WHERE stock >= $1" shart bajarilmasa qator oʻzgarmaydi va RowsAffected() nol boʻladi.
// Bu darsda id mavjud deb faraz qilinadi: nol qator "qoldiq yetmaydi" degani (mavjud boʻlmagan id ham nol qator beradi;
// ikkisini ajratish uchun alohida SELECT kerak, BuyLocked esa ErrNoRows orqali ErrNotFound beradi).
func (s *Store) Buy(ctx context.Context, id int64, qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}
	tag, err := s.pool.Exec(ctx,
		"UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1", qty, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotEnoughStock
	}
	return nil
}
