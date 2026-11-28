package importer

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PG — Saver: partiyani Postgres'ga bitta transaction'da yozadi (GB-5.07).
type PG struct {
	pool *pgxpool.Pool
}

// NewPG pool ustidagi Saver qaytaradi. Pool bir vaqtda koʻp goroutine uchun xavfsiz (GB-5.05).
func NewPG(pool *pgxpool.Pool) *PG {
	return &PG{pool: pool}
}

// upsert: shu nomli mahsulot yoʻq boʻlsa qoʻshadi, bor boʻlsa narxi va qoldigʻini yangilaydi.
// excluded — qoʻshilmoqchi boʻlgan (UNIQUE'ga urilgan) qator. Shuning uchun importni qayta ishga tushirish xavfsiz.
// Eslatma: qator yangilansa ham identity ketma-ketligi bitta qiymatni sarflaydi (ID'larda boʻshliq boʻladi).
const upsert = `
INSERT INTO products (name, price, stock) VALUES ($1, $2, $3)
ON CONFLICT (name) DO UPDATE SET price = excluded.price, stock = excluded.stock`

// endTimeout — Commit yoki Rollback javobini kutishning yuqori chegarasi (baza yoki tarmoq javob bermay qolsa).
const endTimeout = 5 * time.Second

// SaveBatch rows'ni bitta transaction'da yozadi. ctx bekor boʻlsa (Ctrl+C yoki BatchTimeout), Exec toʻxtaydi va
// transaction rollback boʻladi. Commit va Rollback esa alohida end ctx bilan: context.WithoutCancel ctx'ning bekor
// qilinishini (va muddatini) olib tashlaydi, WithTimeout oʻz chegarasini qoʻyadi. Natija: Exec'lari tugagan
// partiyaning Commit'i Ctrl+C'dan yarmida uzilmaydi, lekin abadiy ham kutmaydi; partiya yo toʻliq, yo yoʻq.
func (s *PG) SaveBatch(ctx context.Context, rows []Product) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	end, cancel := context.WithTimeout(context.WithoutCancel(ctx), endTimeout)
	defer cancel()
	defer tx.Rollback(end) // Commit'dan keyin hech narsa qilmaydi
	for _, p := range rows {
		if _, err := tx.Exec(ctx, upsert, p.Name, p.Price, p.Stock); err != nil {
			return err
		}
	}
	return tx.Commit(end)
}
