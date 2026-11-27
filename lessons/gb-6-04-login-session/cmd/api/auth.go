package main

// Login va bekor qilinadigan session'lar (GB-6.04).
//
//	POST /login        parol toʻgʻri → yangi session: tokenning hash'i bazaga, tokenning oʻzi cookie'ga
//	requireUser        cookie → hash → sessions qatori → foydalanuvchi context'ga; yoʻq yoki eski → 401
//	POST /logout       cookie'dagi session qatorini oʻchiradi: shu token bilan keyingi soʻrov 401
//	POST /logout-all   foydalanuvchining hamma session'larini oʻchiradi (VAZIFA, TASK.md)

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/abdurahmon27/bozor/internal/password"
	"github.com/abdurahmon27/bozor/internal/store"
)

// sessionCookie — cookie nomi. sessionTTL — session'ning eng uzun umri: login'dan 24 soat keyin u
// baribir tugaydi (absolute timeout). 24 soat — bizning tanlov; OWASP Session Management Cheat Sheet'ning
// misoli ish kuni ilovasi uchun 4–8 soat, harakatsizlik (idle) muddati esa alohida (README).
const (
	sessionCookie = "session"
	sessionTTL    = 24 * time.Hour
)

// newToken — session tokeni: 32 tasodifiy bayt (256 bit), cookie'ga sigʻadigan matn qilib (43 belgi).
func newToken() string {
	b := make([]byte, 32)
	rand.Read(b) // Go 1.24+: xato qaytarmaydi
	return base64.RawURLEncoding.EncodeToString(b)
}

// hashToken — bazaga yoziladigan va bazadan qidiriladigan qiymat: tokenning SHA-256 hash'i.
func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}

// loginInput — POST /login body'si.
type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// login — POST /login: parol toʻgʻri boʻlsa yangi session ochadi va tokenni cookie'da qaytaradi.
// Har login yangi token oladi, so'rovda eski session cookie'si bo'lsa, uning qatori o'chiriladi: OWASP
// "Renew the Session ID After Any Privilege Level Change" — yangi ID, eskisi yo'q qilinsin (session fixation).
func (a *api) login(w http.ResponseWriter, r *http.Request) {
	var in loginInput
	if !decodeJSON(w, r, &in) {
		return
	}
	u, err := a.checkPassword(r.Context(), in.Email, in.Password)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := a.endCurrentSession(r); err != nil {
		writeError(w, err)
		return
	}
	token := newToken()
	expires := time.Now().Add(sessionTTL)
	err = a.st.CreateSession(r.Context(), hashToken(token), u.ID, expires)
	if err != nil {
		writeError(w, err)
		return
	}
	setSessionCookie(w, token, expires)
	writeJSON(w, http.StatusOK, u)
}

// setSessionCookie — tokenni cookie qilib yuboradi. Atributlar GB-6.07'da batafsil:
// HttpOnly — sahifadagi skript cookie'ni oʻqiy olmaydi; Secure — brauzer uni faqat HTTPS orqali
// yuboradi (brauzerlar va curl http://localhost'ni istisno qiladi, MDN: Set-Cookie); SameSite=Lax — boshqa
// saytdan kelgan POST'ga qoʻshilmaydi. Cache-Control: no-store — session ID bor javob keshda qolmasin (OWASP).
func setSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("Cache-Control", "no-store")
}

// clearSessionCookie — brauzerga cookie'ni oʻchirishni aytadi (MaxAge: -1 → "Max-Age=0").
// Bu faqat mijoz tomoni: session'ni bekor qiladigan narsa — bazadagi qatorni oʻchirish.
func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// errBadLogin — email yoki parol notoʻgʻri (problemFor → 401, ikkalasiga bir xil matn).
var errBadLogin = errors.New("email yoki parol notoʻgʻri")

// dummyHash — yoʻq email uchun ham argon2id bir marta ishlaydi: javob vaqti email bazada bor-yoʻqligini
// bildirmasin (OWASP Authentication Cheat Sheet, "Authentication Responses").
var dummyHash = password.Hash("bozor: email topilmadi")

// checkPassword — Verify (GB-6.03) ustidagi login tekshiruvi: email boʻyicha foydalanuvchi, keyin argon2id Verify.
func (a *api) checkPassword(ctx context.Context, email, pw string) (store.User, error) {
	u, err := a.st.UserByEmail(ctx, email)
	if errors.Is(err, store.ErrUserNotFound) {
		password.Verify(pw, dummyHash)
		return store.User{}, fmt.Errorf("%w: %w", errBadLogin, err)
	}
	if err != nil {
		return store.User{}, err
	}
	ok, err := password.Verify(pw, u.PasswordHash)
	if err != nil {
		return store.User{}, fmt.Errorf("foydalanuvchi %d: %w", u.ID, err)
	}
	if !ok {
		return store.User{}, fmt.Errorf("%w: foydalanuvchi %d", errBadLogin, u.ID)
	}
	u.PasswordHash = "" // hash bundan keyin kerak emas: javobga yoki log'ga tasodifan chiqmasin
	return u, nil
}

// userKey — context'dagi kalit turi. Oʻz turimiz: boshqa paketning kaliti bilan toʻqnashmaydi
// (go doc context.WithValue: kalit string yoki boshqa ichki tur boʻlmasin).
type userKey struct{}

// requireUser — himoyalangan route'lar uchun middleware (GB-6.01 shakli): cookie'dagi token hash'i
// boʻyicha session egasini topadi va uni soʻrov context'iga qoʻyadi. Topilmasa — 401, handler ishlamaydi.
func (a *api) requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookie)
		if err != nil {
			writeError(w, err) // http.ErrNoCookie → 401
			return
		}
		u, err := a.st.UserBySession(r.Context(), hashToken(c.Value))
		if err != nil {
			writeError(w, err) // store.ErrNoSession → 401
			return
		}
		w.Header().Set("Cache-Control", "no-store") // shaxsiy javob: keshda qolmasin (OWASP)
		ctx := context.WithValue(r.Context(), userKey{}, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// currentUser — requireUser context'ga qoʻygan foydalanuvchi. Route requireUser bilan oʻralmagan
// boʻlsa — bu dasturchi xatosi: panic, recoverPanic uni 500 qiladi (jim User{} qaytarishdan xavfsizroq).
func currentUser(r *http.Request) store.User {
	u, ok := r.Context().Value(userKey{}).(store.User)
	if !ok {
		panic("currentUser: route requireUser bilan oʻralmagan")
	}
	return u
}

// endCurrentSession — so'rovdagi session cookie'sining qatorini o'chiradi; cookie yo'q bo'lsa hech narsa
// qilmaydi. logout va login (eski session'ni yo'q qilish) ishlatadi.
func (a *api) endCurrentSession(r *http.Request) error {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return nil // cookie yoʻq: oʻchiradigan session ham yoʻq
	}
	return a.st.DeleteSession(r.Context(), hashToken(c.Value))
}

// logout — POST /logout: cookie'dagi session qatorini o'chiradi va brauzerga cookie'ni o'chirishni aytadi.
// requireUser'siz: session allaqachon eskirgan yoki o'chirilgan bo'lsa ham javob 204 va cookie tozalanadi.
func (a *api) logout(w http.ResponseWriter, r *http.Request) {
	if err := a.endCurrentSession(r); err != nil {
		writeError(w, err)
		return
	}
	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}
