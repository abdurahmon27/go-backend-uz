package main

// Route'lar (GB-4.05). Ochiq: /health, POST /login, POST /logout. Qolganlari requireUser bilan oʻralgan
// (auth.go): session tirik boʻlmasa, handler'gacha yetib bormaydi — 401. POST body'li route'lar requireJSON
// bilan (GB-6.01, 415). GB-6.07: POST /login oldida IP boʻyicha rate limit (ratelimit.go).

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
)

// api — handler'lar store'ni va login limitlarini shu struct orqali koʻradi (GB-5.10).
type api struct {
	st      *store.Store
	loginIP *limiters // GB-6.07: bitta IP'dan login urinishlari
	// VAZIFA (TASK.md, 1-qadam): loginEmail — bitta email'ga login urinishlari.
}

// newMux — hamma route'lar shu yerda: main ham, testlar ham bir xil mux'ni ishlatadi.
func newMux(st *store.Store) *http.ServeMux {
	a := &api{
		st:      st,
		loginIP: newLimiters(12*time.Second, 5), // 5 ta urinish, keyin har 12 soniyada bittadan
		// VAZIFA (TASK.md, 1-qadam): loginEmail — 10 ta urinish, keyin har daqiqada bittadan.
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.Handle("POST /login", requireJSON(http.HandlerFunc(a.login))) // VIDEO: darsgacha — rate limit yoʻq
	mux.Handle("GET /me", a.requireUser(http.HandlerFunc(a.me)))
	mux.HandleFunc("POST /logout", a.logout) // requireUser'siz: eskirgan cookie ham tozalansin (GB-6.04)
	mux.Handle("POST /logout-all", a.requireUser(http.HandlerFunc(a.logoutAll)))
	mux.Handle("GET /orders/{id}", a.requireUser(http.HandlerFunc(a.getOrder)))
	mux.Handle("POST /orders/{id}/cancel", a.requireUser(http.HandlerFunc(a.cancelOrder)))
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

// getOrder — GET /orders/{id}: faqat oʻz buyurtmasi (GB-6.06), begona — 404.
func (a *api) getOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	o, err := a.st.Order(r.Context(), id, currentUser(r).ID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

// cancelOrder — POST /orders/{id}/cancel: buyurtmani bekor qiladi. Body yoʻq: shuning uchun requireJSON
// yoʻq va oddiy HTML forma ham shu route'ga yetib keladi — CSRF misoli aynan shu route'da (GB-6.07).
// Egalik tekshiruvi (GB-6.06) CSRF'ni toʻxtatmaydi: soxta soʻrov egasining oʻz brauzeridan, oʻz cookie'si bilan keladi.
func (a *api) cancelOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	o, err := a.st.CancelOrder(r.Context(), id, currentUser(r).ID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

// pathID yoʻldagi {id}'ni oʻqiydi (GB-4.05). Son boʻlmasa 400 yozadi (errBadID) va false qaytaradi.
func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, fmt.Errorf("%w: %w", errBadID, err))
		return 0, false
	}
	return id, true
}
