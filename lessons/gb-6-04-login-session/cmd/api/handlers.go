package main

// Route'lar (GB-4.05). Ochiq: /health, POST /users (GB-6.03), POST /login, POST /logout. Himoyalangan
// route'lar requireUser bilan oʻralgan (auth.go): cookie'dagi session tirik boʻlmasa, handler'gacha
// yetib bormaydi — 401. POST body'li route'lar requireJSON bilan (GB-6.01, 415).

import (
	"fmt"
	"net/http"

	"github.com/abdurahmon27/bozor/internal/store"
)

// api — handler'lar bazani shu struct orqali koʻradi (GB-5.10).
type api struct {
	st *store.Store
}

// newMux — hamma route'lar shu yerda: main ham, testlar ham bir xil mux'ni ishlatadi.
func newMux(st *store.Store) *http.ServeMux {
	a := &api{st: st}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.Handle("POST /users", requireJSON(http.HandlerFunc(a.register))) // GB-6.03
	mux.Handle("POST /login", requireJSON(http.HandlerFunc(a.login)))
	mux.Handle("GET /me", a.requireUser(http.HandlerFunc(a.me)))
	mux.HandleFunc("POST /logout", a.logout) // requireUser'siz: eskirgan cookie ham tozalansin
	// VAZIFA (TASK.md, 2-qadam): POST /logout-all → a.logoutAll, requireUser bilan oʻralgan
	return mux
}

// healthHandler — server tirikligini tekshirish: har soʻrovga «ok» (GB-4.03).
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

// me — GET /me: session egasi. requireUser uni context'ga qoʻygan.
func (a *api) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, currentUser(r))
}
