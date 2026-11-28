// api — bozor API'sining login qismi, ikki usulda (GB-6.05): POST /login parol toʻgʻri boʻlsa
// ham session cookie'si (GB-6.04), ham JWT beradi — ikkalasini bir xil soʻrovlarda solishtirish uchun.
// Ishda bitta usul tanlanadi; bozor uchun tanlov — session (README).
//
// Store xotirada (internal/store): bu darsda gap bazada emas. Server qayta ishga tushsa, session'lar
// ham, JWT kaliti ham yangidan boshlanadi: eski cookie'lar ham, eski JWT'lar ham oʻladi.
//
//	go run ./cmd/api             → bozor API: http://localhost:8765
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
	"github.com/abdurahmon27/bozor/internal/token"
)

// newServer — bozor API serveri (GB-4.04). Handler — newHandler: middleware zanjiri, ichida newMux (GB-6.01).
func newServer(addr string, st *store.Store, jwtKey []byte) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newHandler(st, jwtKey),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}

func run() error {
	st := store.New()
	if err := st.Seed(context.Background()); err != nil {
		return err
	}
	// JWT kaliti: har ishga tushishda yangi 32 bayt. Ishda u secret sifatida tashqaridan beriladi (GB-8.01).
	srv := newServer(":8765", st, token.NewKey())
	fmt.Println("bozor API: http://localhost:8765")
	return srv.ListenAndServe()
}
