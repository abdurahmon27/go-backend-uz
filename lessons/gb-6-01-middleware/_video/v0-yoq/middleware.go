package main

// _video/v0-yoq: middleware yoʻq — GB-4.08 dagidek, server'ga toʻgʻridan-toʻgʻri mux (video: hook, panic-run).

import (
	"net/http"

	"github.com/abdurahmon27/bozor/internal/store"
)

func newHandler(st *store.Store) http.Handler {
	return newMux(st)
}
