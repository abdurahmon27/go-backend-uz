package main

// Roʻyxatdan oʻtish (GB-6.03): POST /users. GB-6.04: foydalanuvchilar endi Postgres'dagi users jadvalida
// (00004_users.sql); handler GB-6.03'dagi bilan bir xil, faqat a.users.Add o'rniga a.st.AddUser.

import (
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/abdurahmon27/bozor/internal/password"
)

// userInput — POST /users body'si.
type userInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Parol uzunligi (OWASP Authentication Cheat Sheet, NIST SP 800-63B): ikkinchi omil (MFA) yoʻq
// boʻlsa — kamida 15 belgi; uzun iboralar uchun kamida 64 belgigacha ruxsat.
const (
	minPassword = 15
	maxPassword = 64
)

// validateUser — 422 uchun maydon xatolari (GB-6.02'dagi fieldErrors). Belgilar baytda emas,
// rune'da sanaladi (GB-2.10): «oʻ» bitta belgi.
func validateUser(in userInput) fieldErrors {
	f := fieldErrors{}
	if !strings.Contains(in.Email, "@") {
		f["email"] = "email manzil boʻlsin"
	}
	if n := utf8.RuneCountInString(in.Password); n < minPassword || n > maxPassword {
		f["password"] = fmt.Sprintf("%d–%d belgi boʻlsin", minPassword, maxPassword)
	}
	return f
}

// register — POST /users: decodeJSON → validateUser (422) → password.Hash → st.AddUser (409) → 201.
func (a *api) register(w http.ResponseWriter, r *http.Request) {
	var in userInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if f := validateUser(in); len(f) > 0 {
		writeError(w, r, f)
		return
	}
	hash := password.Hash(in.Password) // parolning oʻzi bundan keyin ishlatilmaydi
	u, err := a.st.AddUser(r.Context(), in.Email, hash)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, u)
}
