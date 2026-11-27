// Package bozor: sxema migratsiyalari va namunaviy maʼlumot (GB-5.06). SQL fayllar `go:embed` bilan
// binary ichiga olinadi. Videoda migratsiyani goose buyrugʻi qoʻllaydi (goose up), testlar esa shu
// kutubxona orqali har test uchun yangi sxemaga qoʻllaydi (internal/pgtest).
package bozor

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var files embed.FS

// SeedSQL — seed.sql matni: toʻrt mahsulot va uch buyurtma (faqat laboratoriya va testlar uchun).
//
//go:embed seed.sql
var SeedSQL string

// NewProvider migrations/ papkasidagi fayllarni db bazasiga qoʻllaydigan goose Provider qaytaradi.
// GB-5.06'dagidan farqi: sessiya qulfi (advisory lock) yoʻq. U butun bazaga bitta qulf, testlar esa
// har biri oʻz sxemasida parallel migratsiya qiladi va qulfni kutib sekinlashardi. Bir nechta server
// nusxasi bir vaqtda migratsiya qiladigan joyda (GB-8.06) qulf qaytadi.
func NewProvider(db *sql.DB) (*goose.Provider, error) {
	sub, err := fs.Sub(files, "migrations")
	if err != nil {
		return nil, err
	}
	return goose.NewProvider(goose.DialectPostgres, db, sub)
}
