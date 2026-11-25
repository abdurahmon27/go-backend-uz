// Guru javobi: parol muhit oʻzgaruvchisida emas, secret faylida (/run/secrets/db_parol),
// shuning uchun docker inspect uni koʻrsatmaydi. Qiyin javobidagi qayta urinish ham shu yerda.
// daftar.go va kut.go server/ dagi bilan bir xil (paket mustaqil yigʻilsin deb nusxa).
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// sozlama DATABASE_URL'ni (parolsiz) oʻqiydi va parolni DB_PASSWORD_FILE faylidan qoʻshadi.
func sozlama() (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	if f := os.Getenv("DB_PASSWORD_FILE"); f != "" {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		cfg.ConnConfig.Password = strings.TrimSpace(string(b))
	}
	return cfg, nil
}

// ulan: Qiyin javobidagi kabi, faqat tayyor sozlama bilan.
func ulan(cfg *pgxpool.Config) (*pgxpool.Pool, error) {
	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	oxiri := time.Now().Add(30 * time.Second)
	kutish := 100 * time.Millisecond
	for urinish := 1; ; urinish++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = db.Ping(ctx)
		cancel()
		if err == nil {
			return db, nil
		}
		if time.Now().Add(kutish).After(oxiri) {
			db.Close()
			return nil, err
		}
		fmt.Printf("server: baza hali tayyor emas (%d-urinish), %v kutaman\n", urinish, kutish)
		time.Sleep(kutish)
		kutish = min(kutish*2, 5*time.Second)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8077"
	}
	cfg, err := sozlama()
	if err != nil {
		fmt.Fprintln(os.Stderr, "server: sozlama:", err)
		os.Exit(1)
	}
	db, err := ulan(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "server: baza:", err)
		os.Exit(1)
	}
	defer db.Close()
	if _, err := db.Exec(context.Background(), jadval); err != nil {
		fmt.Fprintln(os.Stderr, "server: jadval:", err)
		os.Exit(1)
	}
	fmt.Println("server: bazaga ulandim")

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "server:", err)
		os.Exit(1)
	}
	srv := &http.Server{Handler: daftar(db)}
	go srv.Serve(ln)
	fmt.Println("server: :" + port + " da tinglayapman")
	kut(srv)
}
