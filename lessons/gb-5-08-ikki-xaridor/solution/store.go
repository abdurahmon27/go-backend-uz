package bozor

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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

// BuyLocked ham sotadi, lekin boshqacha: avval qatorni FOR UPDATE bilan qulflaydi, qoldiqni Go'da tekshiradi,
// keyin yozadi. Qulf transaction tugaguncha turadi, shuning uchun boshqa xaridorning FOR UPDATE'i yoki UPDATE'i shu qatorni kutadi (oddiy SELECT kutmaydi).
func (s *Store) BuyLocked(ctx context.Context, id int64, qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var stock int
	err = tx.QueryRow(ctx,
		"SELECT stock FROM products WHERE id = $1 FOR UPDATE", id).Scan(&stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if stock < qty {
		return ErrNotEnoughStock
	}
	_, err = tx.Exec(ctx,
		"UPDATE products SET stock = $1 WHERE id = $2", stock-qty, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
