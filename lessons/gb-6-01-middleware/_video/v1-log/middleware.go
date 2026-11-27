package main

// _video/v1-log: birinchi logRequests — status hali yoʻq (video: shakl-kod).

import (
	"log"
	"net/http"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
)

func newHandler(st *store.Store) http.Handler {
	return logRequests(newMux(st))
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %q %s", r.Method, r.URL.Path, time.Since(start))
	})
}
