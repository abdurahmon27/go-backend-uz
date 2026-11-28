package main

// Roʻyxatdan oʻtish (GB-6.03): POST /users. Parol faqat password.Hash'ga kiradi va undan
// keyin hech qayerda qolmaydi: storega satr, mijozga esa faqat id va email.

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

// register — POST /users: decodeJSON → validateUser (422) → password.Hash → users.Add (409) → 201.
func (a *api) register(w http.ResponseWriter, r *http.Request) {
	var in userInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if f := validateUser(in); len(f) > 0 {
		writeError(w, f)
		return
	}
	hash := password.Hash(in.Password) // parolning oʻzi bundan keyin ishlatilmaydi
	u, err := a.users.Add(in.Email, hash)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, u)
}
