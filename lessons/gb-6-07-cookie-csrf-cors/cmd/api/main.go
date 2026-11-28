// api — bozor API'si, GB-6.07 laboratoriyasi: login va session'lar (GB-6.04), buyurtmalar va egalik
// tekshiruvi (GB-6.06), endi brauzer tomoni ham: session cookie atributlari, CSRF tekshiruvi,
// CORS allow-list va POST /login'ga rate limit.
//
// Store bu laboratoriyada xotirada (internal/store): mavzu brauzer va server orasidagi qoidalar,
// Postgres emas. Docker kerak emas:
//
//	go run ./cmd/api             → bozor API: http://localhost:8765
//	go run ./cmd/begona          → «begona sayt»: http://localhost:9090 va http://127.0.0.1:9090
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
)

// newServer — bozor API serveri (GB-4.04). Handler — newHandler: middleware zanjiri, ichida newMux (GB-6.01).
func newServer(addr string, st *store.Store) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newHandler(st),
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
	srv := newServer(":8765", st)
	fmt.Println("bozor API: http://localhost:8765")
	return srv.ListenAndServe()
}
