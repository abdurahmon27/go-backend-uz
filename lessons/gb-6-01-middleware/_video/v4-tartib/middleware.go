package main

// _video/v4-tartib: notoʻgʻri tartib — recoverPanic logRequests'dan tashqarida (video: bashorat3, tartib-run).

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
)

func newHandler(st *store.Store) http.Handler {
	return recoverPanic(logRequests(newMux(st)))
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

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				log.Printf("panic: %v\n%s", v, debug.Stack())
				writeError(w, http.StatusInternalServerError, apiError{Code: "internal", Message: "ichki xato"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
