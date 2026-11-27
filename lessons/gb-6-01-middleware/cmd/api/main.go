// api — bozor katalog API'si, endi middleware bilan (GB-6.01). Store — GB-4.08'dagi xotiradagi
// katalog (bozor.json'dan oʻqiladi): bu darsning mavzusi Postgres emas, middleware. Middleware faqat
// http.Handler'ni koʻradi, shuning uchun middleware.go GB-5.10'dagi Postgres versiyasiga ham
// oʻzgarishsiz qoʻshiladi.
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

// newServer — bozor API serveri (GB-4.04). Handler — newHandler: middleware zanjiri, ichida newMux (GB-6.01).
func newServer(addr string, st *store.Store) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newHandler(st),
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
