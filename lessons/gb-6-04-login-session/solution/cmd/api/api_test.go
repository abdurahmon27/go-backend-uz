// API testlari: haqiqiy handler'lar, haqiqiy middleware zanjiri, haqiqiy Postgres (mock yoʻq).
// Server ishga tushmaydi: httptest soʻrovni toʻgʻridan-toʻgʻri newHandler'ga beradi (GB-4.03).
// Har test oʻz sxemasida (pgtest.New): seed.sql'dagi ikki foydalanuvchi bilan.
package main

import (
	"context"
	"crypto/sha256"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/abdurahmon27/bozor/internal/pgtest"
	"github.com/abdurahmon27/bozor/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dilnozaLogin = `{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}`
	jasurLogin   = `{"email":"jasur@bozor.test","password":"choy-non-qand-27"}`
)

// TestMain: server log'i (soʻrov qatorlari, panic stack'i) test chiqishini toʻldirmasin.
func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

// send soʻrovni h'ga beradi; cookie boʻsh boʻlmasa, uni "session" cookie'si qilib qoʻshadi.
func send(h http.Handler, method, path, body, cookie string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if cookie != "" {
		req.AddCookie(&http.Cookie{Name: sessionCookie, Value: cookie})
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// testAPI — yangi sxema (seed.sql bilan) va shu bazaga ulangan butun zanjir.
func testAPI(t *testing.T) (http.Handler, *pgxpool.Pool) {
	t.Helper()
	pool := pgtest.New(t)
	return newHandler(store.New(pool)), pool
}

// login body bilan POST /login yuboradi va javobdagi session cookie'sini qaytaradi.
func login(t *testing.T, h http.Handler, body string) *http.Cookie {
	t.Helper()
	rec := send(h, "POST", "/login", body, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /login: status %d, kutilgan 200: %s", rec.Code, rec.Body)
	}
	for _, c := range rec.Result().Cookies() {
		if c.Name == sessionCookie {
			return c
		}
	}
	t.Fatalf("POST /login javobida %q cookie yoʻq: %v", sessionCookie, rec.Header())
	return nil
}

func TestLoginCookie(t *testing.T) {
	h, pool := testAPI(t)
	c := login(t, h, dilnozaLogin)
	if !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteLaxMode || c.Path != "/" {
		t.Errorf("cookie atributlari: %v; kutilgan HttpOnly, Secure, SameSite=Lax, Path=/", c)
	}
	if len(c.Value) != 43 {
		t.Errorf("token uzunligi %d, kutilgan 43 (32 bayt, base64url)", len(c.Value))
	}
	// Bazada tokenning oʻzi emas, uning SHA-256 hash'i.
	var stored []byte
	if err := pool.QueryRow(context.Background(), "SELECT token_hash FROM sessions").Scan(&stored); err != nil {
		t.Fatalf("sessions qatori oʻqilmadi: %v", err)
	}
	sum := sha256.Sum256([]byte(c.Value))
	if string(stored) != string(sum[:]) {
		t.Errorf("token_hash = %x, kutilgan sha256(token) = %x", stored, sum)
	}
	// Ikkinchi login — boshqa token.
	if c2 := login(t, h, dilnozaLogin); c2.Value == c.Value {
		t.Error("ikki login bir xil token oldi: har login yangi tasodifiy token olishi kerak")
	}
}

func TestLoginFails(t *testing.T) {
	h, pool := testAPI(t)
	for _, body := range []string{
		`{"email":"dilnoza@bozor.test","password":"notoʻgʻri"}`,
		`{"email":"yoq@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}`,
	} {
		rec := send(h, "POST", "/login", body, "")
		if rec.Code != http.StatusUnauthorized || !strings.Contains(rec.Body.String(), "email yoki parol notoʻgʻri") {
			t.Errorf("%s → %d %s; kutilgan 401 «email yoki parol notoʻgʻri»", body, rec.Code, rec.Body)
		}
		if len(rec.Result().Cookies()) != 0 {
			t.Errorf("%s: muvaffaqiyatsiz login cookie bermasligi kerak", body)
		}
	}
	if rec := send(h, "POST", "/login", `{"email":`, ""); rec.Code != http.StatusBadRequest {
		t.Errorf("buzilgan JSON → %d, kutilgan 400", rec.Code)
	}
	var n int
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM sessions").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("sessions'da %d qator, kutilgan 0", n)
	}
}

func TestMe(t *testing.T) {
	h, _ := testAPI(t)
	c := login(t, h, dilnozaLogin)
	tests := []struct {
		name, cookie string
		status       int
		want         string
	}{
		{"cookie_yoq", "", 401, "avval login qiling"},
		{"begona_token", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 401, "avval login qiling"},
		{"tirik_session", c.Value, 200, `{"id":1,"email":"dilnoza@bozor.test"}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := send(h, "GET", "/me", "", tc.cookie)
			if rec.Code != tc.status || !strings.Contains(rec.Body.String(), tc.want) {
				t.Errorf("GET /me: %d %s; kutilgan %d %s", rec.Code, rec.Body, tc.status, tc.want)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	h, _ := testAPI(t)
	c := login(t, h, dilnozaLogin)
	rec := send(h, "POST", "/logout", "", c.Value)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("POST /logout: %d %s, kutilgan 204", rec.Code, rec.Body)
	}
	if got := rec.Header().Get("Set-Cookie"); !strings.Contains(got, "Max-Age=0") {
		t.Errorf("logout javobidagi Set-Cookie %q: cookie oʻchirilishi kerak (Max-Age=0)", got)
	}
	// Eski token nusxasi (brauzer cookie'ni oʻchirmagan deylik) endi ishlamaydi.
	if rec := send(h, "GET", "/me", "", c.Value); rec.Code != http.StatusUnauthorized {
		t.Errorf("logout'dan keyin GET /me: %d, kutilgan 401", rec.Code)
	}
	// Logout idempotent: session yo'q bo'lsa ham 204 va cookie tozalanadi (eskirgan cookie brauzerda qolmasin).
	if rec := send(h, "POST", "/logout", "", c.Value); rec.Code != http.StatusNoContent || !strings.Contains(rec.Header().Get("Set-Cookie"), "Max-Age=0") {
		t.Errorf("ikkinchi logout: %d %q, kutilgan 204 va Max-Age=0", rec.Code, rec.Header().Get("Set-Cookie"))
	}
}

func TestReloginDropsOldSession(t *testing.T) {
	h, pool := testAPI(t)
	old := login(t, h, dilnozaLogin)
	// Eski cookie bilan qayta login (OWASP: yangi ID, eskisi yo'q qilinsin).
	req := httptest.NewRequest("POST", "/login", strings.NewReader(dilnozaLogin))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(old)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("qayta login: %d %s", rec.Code, rec.Body)
	}
	if rec := send(h, "GET", "/me", "", old.Value); rec.Code != http.StatusUnauthorized {
		t.Errorf("qayta login'dan keyin eski token: GET /me %d, kutilgan 401", rec.Code)
	}
	var n int
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM sessions").Scan(&n); err != nil || n != 1 {
		t.Errorf("sessions'da %d qator (%v), kutilgan 1", n, err)
	}
}

func TestHeaders(t *testing.T) {
	h, _ := testAPI(t)
	rec := send(h, "POST", "/login", dilnozaLogin, "")
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("login javobida Cache-Control %q, kutilgan no-store", got)
	}
	c := rec.Result().Cookies()[0]
	if got := send(h, "GET", "/me", "", c.Value).Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("/me javobida Cache-Control %q, kutilgan no-store", got)
	}
	if got := send(h, "GET", "/me", "", "").Header().Get("WWW-Authenticate"); got == "" {
		t.Error("401 javobida WWW-Authenticate yoʻq (RFC 9110 §11.6.1)")
	}
}

func TestRegisterThenLogin(t *testing.T) {
	h, _ := testAPI(t)
	body := `{"email":"malika@bozor.test","password":"oʻrik-uzum-qovun-ariq-9"}`
	if rec := send(h, "POST", "/users", body, ""); rec.Code != http.StatusCreated || strings.Contains(rec.Body.String(), "argon2id") {
		t.Fatalf("POST /users: %d %s; kutilgan 201, hash'siz", rec.Code, rec.Body)
	}
	if rec := send(h, "POST", "/users", body, ""); rec.Code != http.StatusConflict {
		t.Errorf("takror email: %d, kutilgan 409", rec.Code)
	}
	if rec := send(h, "POST", "/users", `{"email":"x","password":"qisqa"}`, ""); rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("notoʻgʻri email/parol: %d, kutilgan 422", rec.Code)
	}
	login(t, h, body)
}

func TestNotFoundStays404(t *testing.T) {
	h, _ := testAPI(t)
	if rec := send(h, "POST", "/nope", "x", ""); rec.Code != http.StatusNotFound {
		t.Errorf("POST /nope: %d, kutilgan 404 (requireJSON faqat route'larda)", rec.Code)
	}
}

func TestExpiredSession(t *testing.T) {
	h, pool := testAPI(t)
	c := login(t, h, dilnozaLogin)
	_, err := pool.Exec(context.Background(), "UPDATE sessions SET expires_at = now() - interval '1 minute'")
	if err != nil {
		t.Fatal(err)
	}
	if rec := send(h, "GET", "/me", "", c.Value); rec.Code != http.StatusUnauthorized {
		t.Errorf("muddati oʻtgan session bilan GET /me: %d %s, kutilgan 401", rec.Code, rec.Body)
	}
}

func TestCurrentUserWithoutMiddleware(t *testing.T) {
	// Route requireUser'siz ulangan boʻlsa: panic → recoverPanic → 500 (jim User{} emas).
	a := &api{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /me", a.me)
	h := recoverPanic(mux)
	if rec := send(h, "GET", "/me", "", ""); rec.Code != http.StatusInternalServerError {
		t.Errorf("requireUser'siz /me: %d, kutilgan 500", rec.Code)
	}
}

func TestRequireJSON(t *testing.T) {
	h, _ := testAPI(t)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(dilnozaLogin))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Content-Type'siz POST /login: %d, kutilgan 415", rec.Code)
	}
}
