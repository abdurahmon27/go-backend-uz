// Package bozor: sxema migratsiyalari. SQL fayllar `go:embed` bilan binary ichiga olinadi,
// goose kutubxonasi ularni bazaga qo'llaydi. Videoda goose buyruq qatoridan ishlaydi (goose up),
// serverda esa shu kutubxona (GB-8.06). Bu fayl va migrate_test.go tekshiruv uchun.
package bozor

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

//go:embed migrations/*.sql
var files embed.FS

// NewProvider migrations/ papkasidagi fayllarni db bazasiga qoʻllaydigan goose Provider qaytaradi.
// Postgres'ning sessiya qulfi (advisory lock) bir necha nusxa bir vaqtda ishga tushsa ham,
// migratsiyani faqat bittasi qoʻllashini taʼminlaydi.
func NewProvider(db *sql.DB) (*goose.Provider, error) {
	sub, err := fs.Sub(files, "migrations")
	if err != nil {
		return nil, err
	}
	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return nil, err
	}
	return goose.NewProvider(goose.DialectPostgres, db, sub, goose.WithSessionLocker(locker))
}
