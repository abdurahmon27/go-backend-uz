package main

// _video/v2-status: statusRecorder bilan logRequests (video: rec-kod, log2-kod).

import (
	"log"
	"net/http"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
)

func newHandler(st *store.Store) http.Handler {
	return logRequests(newMux(st))
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %q %d %s", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}
