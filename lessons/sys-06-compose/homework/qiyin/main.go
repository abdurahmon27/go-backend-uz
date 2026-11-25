// Qiyin javobi: server bazani oʻzi kutadi (healthcheck boʻlmasa ham, masalan Kubernetes'da).
// daftar.go va kut.go server/ dagi bilan bir xil (paket mustaqil yigʻilsin deb nusxa).
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ulan bazaga ulanguncha qayta urinadi: kutish har safar ikki baravar oshadi (backoff), jami 30 soniyagacha.
func ulan(url string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err // manzil notoʻgʻri: qayta urinishdan foyda yoʻq
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
	db, err := ulan(os.Getenv("DATABASE_URL"))
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
