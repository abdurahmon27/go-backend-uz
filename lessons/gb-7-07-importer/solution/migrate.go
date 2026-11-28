// Package bozor: sxema migratsiyalari (GB-5.06, GB-5.10'dagi uchta fayl). SQL fayllar `go:embed` bilan binary
// ichiga olinadi. Videoda migratsiyani goose buyrugʻi qoʻllaydi (goose up), testlar esa shu kutubxona orqali har
// test uchun yangi sxemaga qoʻllaydi (internal/pgtest).
package bozor

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var files embed.FS

// NewProvider migrations/ papkasidagi fayllarni db bazasiga qoʻllaydigan goose Provider qaytaradi
// (sessiya qulfisiz: testlar har biri oʻz sxemasida parallel migratsiya qiladi, GB-5.10).
func NewProvider(db *sql.DB) (*goose.Provider, error) {
	sub, err := fs.Sub(files, "migrations")
	if err != nil {
		return nil, err
	}
	return goose.NewProvider(goose.DialectPostgres, db, sub)
}
