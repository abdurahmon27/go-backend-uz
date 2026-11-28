package main

// Route'lar (GB-4.05). Ochiq: /health, POST /users (GB-6.03), POST /login, POST /logout. Himoyalangan
// route'lar requireUser bilan oʻralgan (auth.go): session cookie'si ham, Bearer JWT ham yaroqsiz boʻlsa,
// handler'gacha yetib bormaydi — 401. POST body'li route'lar requireJSON bilan (GB-6.01, 415).
// GB-6.04'dan farqi: api'da JWT kaliti (jwtKey) bor.

import (
	"fmt"
	"net/http"

	"github.com/abdurahmon27/bozor/internal/store"
)

// api — handler'lar store'ni (GB-5.10) va JWT kalitini (GB-6.05) shu struct orqali koʻradi.
type api struct {
	st     *store.Store
	jwtKey []byte
}

// newMux — hamma route'lar shu yerda: main ham, testlar ham bir xil mux'ni ishlatadi.
func newMux(st *store.Store, jwtKey []byte) *http.ServeMux {
	a := &api{st: st, jwtKey: jwtKey}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.Handle("POST /users", requireJSON(http.HandlerFunc(a.register))) // GB-6.03
	mux.Handle("POST /login", requireJSON(http.HandlerFunc(a.login)))
	mux.Handle("GET /me", a.requireUser(http.HandlerFunc(a.me)))
	mux.HandleFunc("POST /logout", a.logout) // requireUser'siz: eskirgan cookie ham tozalansin
	mux.Handle("POST /logout-all", a.requireUser(http.HandlerFunc(a.logoutAll)))
	return mux
}

// healthHandler — server tirikligini tekshirish: har soʻrovga «ok» (GB-4.03).
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

// me — GET /me: kim ekanimiz. requireUser uni context'ga qoʻygan.
// Session yoʻli foydalanuvchini store'dan oladi (id + email); JWT yoʻli faqat tokendagi sub'ni biladi (id).
func (a *api) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, currentUser(r))
}
