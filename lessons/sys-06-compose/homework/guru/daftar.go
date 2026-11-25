package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const jadval = `CREATE TABLE IF NOT EXISTS mehmonlar (
	id    bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	vaqt  timestamptz NOT NULL DEFAULT now()
)`

// daftar har bir mehmonni Postgres'ga yozadi va nechanchi ekanini aytadi.
// SYS-0.05 da bu fayl edi, endi bir vaqtdagi soʻrovlarni baza oʻzi tartibga soladi.
func daftar(db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var n int64
		err := db.QueryRow(r.Context(),
			`INSERT INTO mehmonlar DEFAULT VALUES RETURNING id`).Scan(&n)
		if err != nil {
			fmt.Fprintln(os.Stderr, "server: daftar:", err) // tafsilot logga, mijozga emas
			http.Error(w, "bazada xato", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Salom! Siz %d-mehmonsiz\n", n)
	})
}
