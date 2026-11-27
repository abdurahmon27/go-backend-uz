// api — bozor katalog API'si (GB-4.08, 4-modul nazorat nuqtasi): GB-3.06'dagi internal/store
// katalogi HTTP orqali. Katalog ishga tushishda bozor.json'dan oʻqiladi (CLI yozgan fayl) va keyin
// xotirada turadi: API'dagi oʻzgarishlar faylga yozilmaydi, server toʻxtasa yoʻqoladi.
// Saqlash Postgres'ga 5-modulda koʻchadi (GB-5.01: nega baza kerak).
//
//	go run ./cmd/api        → bozor API: http://localhost:8765
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
)

const dataFile = "bozor.json"

// newServer — bozor API serveri (GB-4.04): route'lar newMux'da (GB-4.05), sarlavhalar uchun vaqt chegarasi.
func newServer(addr string, st *store.Store) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newMux(st),
		ReadHeaderTimeout: 5 * time.Second, // soʻrov sarlavhalari 5 s'da kelmasa, ulanish yopiladi
	}
}

func main() {
	st, err := store.Load(dataFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	srv := newServer(":8765", st)
	fmt.Println("bozor API: http://localhost:8765")
	err = srv.ListenAndServe()
	fmt.Fprintln(os.Stderr, "xato:", err)
	os.Exit(1)
}
