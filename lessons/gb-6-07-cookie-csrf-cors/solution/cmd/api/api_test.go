package main

// Darsda koʻrsatilgan qoidalar testlari (GB-6.07): starter'da ham, solution'da ham oʻtadi.
// Soʻrovlar newHandler'ning butun zanjiridan oʻtadi (cors, protectCSRF, rate limit — hammasi).

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/abdurahmon27/bozor/internal/store"
)

const dilnozaLogin = `{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}`

// newTestHandler — Seed'langan xotira store'i bilan butun zanjir.
func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	log.SetOutput(io.Discard) // server log'i test chiqishini toʻldirmasin
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	st := store.New()
	if err := st.Seed(context.Background()); err != nil {
		t.Fatal(err)
	}
	return newHandler(st)
}

// do — bitta soʻrov. ip boʻsh boʻlmasa RemoteAddr shu IP'dan (rate limit kaliti).
func do(h http.Handler, method, path, body, ip string, header map[string]string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://localhost:8765"+path, strings.NewReader(body))
	if body != "" {
		r.Header.Set("Content-Type", "application/json")
	}
	if ip != "" {
		r.RemoteAddr = ip + ":40000"
	}
	for k, v := range header {
		r.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

// loginCookie — Dilnoza login qiladi; session cookie'si qaytadi.
func loginCookie(t *testing.T, h http.Handler) *http.Cookie {
	t.Helper()
	w := do(h, "POST", "/login", dilnozaLogin, "198.51.100.7", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /login: status %d, kutilgan 200: %s", w.Code, w.Body)
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == sessionCookie {
			return c
		}
	}
	t.Fatalf("POST /login javobida %s cookie'si yoʻq", sessionCookie)
	return nil
}

func TestSessionCookieAttributes(t *testing.T) {
	c := loginCookie(t, newTestHandler(t))
	if c.Name != "__Host-session" || !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteLaxMode ||
		c.Path != "/" || c.Domain != "" || c.MaxAge != 86400 {
		t.Errorf("Set-Cookie: %s — kutilgan: __Host-session, Path=/, Domain yoʻq, Max-Age=86400, HttpOnly, Secure, SameSite=Lax", c)
	}
}

func TestCrossOriginProtection(t *testing.T) {
	cases := []struct {
		name   string
		header map[string]string
		want   int
	}{
		{"curl: Sec-Fetch-Site va Origin yoʻq", nil, http.StatusOK},
		{"oʻz sahifasi", map[string]string{"Sec-Fetch-Site": "same-origin"}, http.StatusOK},
		{"localhost:9090 (bir site)", map[string]string{"Sec-Fetch-Site": "same-site", "Origin": "http://localhost:9090"}, http.StatusForbidden},
		{"127.0.0.1:9090 (boshqa site)", map[string]string{"Sec-Fetch-Site": "cross-site", "Origin": "http://127.0.0.1:9090"}, http.StatusForbidden},
		{"eski brauzer: faqat Origin", map[string]string{"Origin": "http://localhost:9090"}, http.StatusForbidden},
		{"ruxsatli frontend", map[string]string{"Sec-Fetch-Site": "same-site", "Origin": "http://localhost:3000"}, http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t) // har holatda buyurtma 41 yana «yangi»
			c := loginCookie(t, h)
			tc.header = mergeCookie(tc.header, c)
			w := do(h, "POST", "/orders/41/cancel", "", "", tc.header)
			if w.Code != tc.want {
				t.Errorf("POST /orders/41/cancel: status %d, kutilgan %d: %s", w.Code, tc.want, w.Body)
			}
		})
	}
}

func mergeCookie(h map[string]string, c *http.Cookie) map[string]string {
	m := map[string]string{"Cookie": c.Name + "=" + c.Value}
	for k, v := range h {
		m[k] = v
	}
	return m
}

func TestCORSAllowList(t *testing.T) {
	h := newTestHandler(t)
	c := loginCookie(t, h)
	ok := do(h, "GET", "/me", "", "", mergeCookie(map[string]string{"Origin": "http://localhost:3000"}, c))
	if got := ok.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("Origin: http://localhost:3000 → Access-Control-Allow-Origin %q, kutilgan shu origin", got)
	}
	if got := ok.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("Access-Control-Allow-Credentials %q, kutilgan true", got)
	}
	bad := do(h, "GET", "/me", "", "", mergeCookie(map[string]string{"Origin": "http://localhost:9090"}, c))
	if got := bad.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Origin: http://localhost:9090 → Access-Control-Allow-Origin %q, kutilgan: sarlavha yoʻq", got)
	}
	if !strings.Contains(bad.Header().Get("Vary"), "Origin") {
		t.Errorf("Vary: %q, kutilgan Origin", bad.Header().Get("Vary"))
	}
	pre := do(h, "OPTIONS", "/login", "", "", map[string]string{"Origin": "http://localhost:3000",
		"Access-Control-Request-Method": "POST", "Access-Control-Request-Headers": "content-type"})
	if pre.Code != http.StatusNoContent || pre.Header().Get("Access-Control-Allow-Headers") != "Content-Type" {
		t.Errorf("preflight: status %d, Allow-Headers %q — kutilgan 204 va Content-Type", pre.Code,
			pre.Header().Get("Access-Control-Allow-Headers"))
	}
}

func TestLoginRateLimitByIP(t *testing.T) {
	h := newTestHandler(t)
	bad := `{"email":"dilnoza@bozor.test","password":"taxmin-1"}`
	for i := 1; i <= 5; i++ {
		if w := do(h, "POST", "/login", bad, "203.0.113.5", nil); w.Code != http.StatusUnauthorized {
			t.Fatalf("%d-urinish: status %d, kutilgan 401", i, w.Code)
		}
	}
	w := do(h, "POST", "/login", bad, "203.0.113.5", nil)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("6-urinish: status %d, kutilgan 429", w.Code)
	}
	if s, err := strconv.Atoi(w.Header().Get("Retry-After")); err != nil || s < 1 || s > 12 {
		t.Errorf("Retry-After %q, kutilgan 1..12 soniya", w.Header().Get("Retry-After"))
	}
	if w := do(h, "POST", "/login", bad, "203.0.113.99", nil); w.Code != http.StatusUnauthorized {
		t.Errorf("boshqa IP: status %d, kutilgan 401 (har IP'ning oʻz chelagi)", w.Code)
	}
}
