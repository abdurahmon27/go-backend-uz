package main

// Login, ikki usulda (GB-6.05). GB-6.04'dagi session'lar oʻz joyida, yonida JWT:
//
//	session (GB-6.04)  tasodifiy token cookie'da; serverda uning hash'i → foydalanuvchi. Har soʻrovda qidiruv.
//	                   Logout = yozuvni oʻchirish: shu token bilan keyingi soʻrov darhol 401.
//	JWT (RFC 7519)     imzolangan token Authorization: Bearer sarlavhasida; serverda hech narsa saqlanmaydi.
//	                   Tekshiruv = imzo va claim'lar, qidiruvsiz. Logout serverda hech narsani oʻchirmaydi:
//	                   token exp'gacha yaroqli. Bekor qilish uchun — denylist (VAZIFA, TASK.md).
//
//	POST /login        parol toʻgʻri → yangi session (cookie) + JWT (javob body'sida)
//	requireUser        authenticate: Bearer bor → JWT tekshiruvi; aks holda cookie → hash → session; boʻlmasa 401
//	POST /logout       cookie'dagi session yozuvini oʻchiradi (GB-6.04)
//	POST /logout-all   foydalanuvchining hamma session'larini oʻchiradi (GB-6.04 vazifasi)

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/abdurahmon27/bozor/internal/password"
	"github.com/abdurahmon27/bozor/internal/store"
	"github.com/abdurahmon27/bozor/internal/token"
)

// sessionCookie — cookie nomi. sessionTTL — session'ning eng uzun umri: login'dan 24 soat keyin u
// baribir tugaydi (absolute timeout). 24 soat — bizning tanlov; OWASP Session Management Cheat Sheet'ning
// misoli ish kuni ilovasi uchun 4–8 soat. Harakatsizlik (idle) muddati bu laboratoriyada yoʻq
// (OWASP: 15–30 daqiqa kam xavfli, 2–5 daqiqa yuqori xavfli ilovaga; gb-6-04 README).
// jwtTTL — JWT'ning umri: qisqa, chunki uni muddatidan oldin bekor qilib boʻlmaydi (OWASP JWT Cheat Sheet).
const (
	sessionCookie = "session"
	sessionTTL    = 24 * time.Hour
	jwtTTL        = 15 * time.Minute
)

// newToken — session tokeni: 32 tasodifiy bayt (256 bit), cookie'ga sigʻadigan matn qilib (43 belgi).
func newToken() string {
	b := make([]byte, 32)
	rand.Read(b) // Go 1.24+: xato qaytarmaydi
	return base64.RawURLEncoding.EncodeToString(b)
}

// hashToken — bazaga yoziladigan va bazadan qidiriladigan qiymat: tokenning SHA-256 hash'i.
func hashToken(t string) []byte {
	sum := sha256.Sum256([]byte(t))
	return sum[:]
}

// loginInput — POST /login body'si.
type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginOutput — POST /login javobi: foydalanuvchi va JWT. Session tokeni esa cookie'da.
type loginOutput struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// login — POST /login: parol toʻgʻri boʻlsa yangi session ochadi (cookie) va JWT beradi (body).
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
	sess := newToken()
	expires := time.Now().Add(sessionTTL)
	err = a.st.CreateSession(r.Context(), hashToken(sess), u.ID, expires)
	if err != nil {
		writeError(w, err)
		return
	}
	jwt, err := token.Issue(u.ID, jwtTTL, a.jwtKey)
	if err != nil {
		writeError(w, err)
		return
	}
	setSessionCookie(w, sess, expires)
	writeJSON(w, http.StatusOK, loginOutput{ID: u.ID, Email: u.Email, Token: jwt})
}

// setSessionCookie — tokenni cookie qilib yuboradi. Atributlar GB-6.07'da batafsil:
// HttpOnly — sahifadagi skript cookie'ni oʻqiy olmaydi; Secure — brauzer uni faqat HTTPS orqali
// yuboradi (brauzerlar va curl http://localhost'ni istisno qiladi, MDN: Set-Cookie); SameSite=Lax — boshqa
// saytdan kelgan POST'ga qoʻshilmaydi. Cache-Control: no-store — session ID bor javob keshda qolmasin (OWASP).
func setSessionCookie(w http.ResponseWriter, value string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    value,
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

// errBadToken — JWT yaroqsiz: imzo, alg, iss, aud, exp yoki denylist. Mijozga hammasi bir xil 401,
// asl sabab (kutubxona xatosi) ikki nuqtadan keyin, faqat log'da (GB-6.02).
var errBadToken = errors.New("JWT yaroqsiz")

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

// bearer — "Authorization: Bearer <token>" sarlavhasidan tokenni ajratadi. Sxema nomida katta-kichik
// harf farqi yoʻq (RFC 9110 §11.1), undan keyin bir yoki bir necha boʻsh joy (RFC 6750 §2.1).
func bearer(r *http.Request) (string, bool) {
	scheme, rest, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}
	t := strings.TrimLeft(rest, " ")
	return t, t != ""
}

// authenticate — soʻrov kimniki? Ikki yoʻl:
//   - Authorization: Bearer bor → JWT: imzo va claim'lar tekshiriladi, store'ga qidiruv yoʻq;
//   - aks holda session cookie'si → hash → store'dagi yozuv (GB-6.04).
//
// Bearer bor-u yaroqsiz boʻlsa, cookie'ga oʻtilmaydi: mijoz qaysi usulni tanlagani aniq, javob 401.
func (a *api) authenticate(r *http.Request) (store.User, error) {
	if s, ok := bearer(r); ok {
		c, err := token.Parse(s, a.jwtKey)
		if err != nil {
			return store.User{}, fmt.Errorf("%w: %w", errBadToken, err)
		}
		// VAZIFA (TASK.md, 3-qadam): c.JTI denylist'da boʻlsa (a.st.TokenRevoked) — errBadToken.
		return store.User{ID: c.UserID}, nil
	}
	// aks holda: session cookie'si (GB-6.04)
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return store.User{}, err // http.ErrNoCookie → 401
	}
	return a.st.UserBySession(r.Context(), hashToken(cookie.Value))
}

// requireUser — himoyalangan route'lar uchun middleware (GB-6.01 shakli): authenticate bilan soʻrov
// egasini topadi va uni context'ga qoʻyadi. Topilmasa — 401, handler ishlamaydi.
func (a *api) requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := a.authenticate(r)
		if err != nil {
			if errors.Is(err, errBadToken) {
				// RFC 6750 §3: Bearer token rad etilganda 401 bilan birga shu sarlavha.
				w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
			}
			writeError(w, err) // ErrNoCookie, ErrNoSession, errBadToken → 401
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
// JWT: serverda oʻchiradigan yozuv yoʻq — Bearer bilan kelgan token exp'gacha yaroqli boʻlib qoladi.
func (a *api) logout(w http.ResponseWriter, r *http.Request) {
	// VAZIFA (TASK.md, 2-qadam): soʻrov Bearer bilan kelgan va token yaroqli boʻlsa, uning jti'sini
	// a.st.RevokeToken bilan denylist'ga qoʻshing (exp vaqti bilan).
	if err := a.endCurrentSession(r); err != nil {
		writeError(w, err)
		return
	}
	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

// logoutAll — POST /logout-all (GB-6.04 vazifasi): "hamma qurilmadan chiqish". Foydalanuvchining
// hamma session qatorlari oʻchadi: telefon va noutbukdagi tokenlar ham keyingi soʻrovda 401 oladi.
// GB-6.05: JWT'larga taʼsir qilmaydi — ularning yozuvi yoʻq. Ularni ham oʻldirish uchun foydalanuvchi
// boʻyicha denylist ("shu vaqtdan oldin berilganlar yaroqsiz") kerak boʻlardi: yana server holati (README).
func (a *api) logoutAll(w http.ResponseWriter, r *http.Request) {
	_, err := a.st.DeleteUserSessions(r.Context(), currentUser(r).ID)
	if err != nil {
		writeError(w, err)
		return
	}
	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}
