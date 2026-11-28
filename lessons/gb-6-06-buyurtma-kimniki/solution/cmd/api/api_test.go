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
	adminLogin   = `{"email":"admin@bozor.test","password":"qovun-anjir-tut-uzum-31"}`
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
		{"tirik_session", c.Value, 200, `{"id":1,"email":"dilnoza@bozor.test","role":"xaridor"}`},
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
	if rec := send(h, "POST", "/users", body, ""); rec.Code != http.StatusCreated || strings.Contains(rec.Body.String(), "argon2id") ||
		!strings.Contains(rec.Body.String(), `"role":"xaridor"`) {
		t.Fatalf("POST /users: %d %s; kutilgan 201, hash'siz, role xaridor (GB-6.06)", rec.Code, rec.Body)
	}
	// GB-6.06: rolni mijoz bera olmaydi — nomaʼlum kalit, 400 (OWASP API3:2023, mass assignment).
	admin := `{"email":"hiyla@bozor.test","password":"oʻrik-uzum-qovun-ariq-9","role":"admin"}`
	if rec := send(h, "POST", "/users", admin, ""); rec.Code != http.StatusBadRequest {
		t.Errorf("body'da role: %d %s, kutilgan 400", rec.Code, rec.Body)
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

// GB-6.04 vazifasi: POST /logout-all — "hamma qurilmadan chiqish".
func TestLogoutAll(t *testing.T) {
	h, _ := testAPI(t)
	phone := login(t, h, dilnozaLogin)  // Dilnoza telefonda
	laptop := login(t, h, dilnozaLogin) // va noutbukda
	jasur := login(t, h, jasurLogin)    // boshqa foydalanuvchi

	rec := send(h, "POST", "/logout-all", "", phone.Value)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("POST /logout-all: status %d, kutilgan 204: %s", rec.Code, rec.Body)
	}
	if got := rec.Header().Get("Set-Cookie"); !strings.Contains(got, "Max-Age=0") {
		t.Errorf("POST /logout-all javobida Set-Cookie %q: cookie oʻchirilsin, clearSessionCookie", got)
	}
	for name, c := range map[string]*http.Cookie{"telefon": phone, "noutbuk": laptop} {
		if rec := send(h, "GET", "/me", "", c.Value); rec.Code != http.StatusUnauthorized {
			t.Errorf("logout-all'dan keyin Dilnozaning %s tokeni: GET /me %d, kutilgan 401", name, rec.Code)
		}
	}
	if rec := send(h, "GET", "/me", "", jasur.Value); rec.Code != http.StatusOK {
		t.Errorf("Jasurning session'i ham oʻchdi: GET /me %d, kutilgan 200", rec.Code)
	}
	if rec := send(h, "POST", "/logout-all", "", ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("cookie'siz POST /logout-all: %d, kutilgan 401", rec.Code)
	}
}

// GB-6.06: buyurtmani faqat egasi koʻradi. Begona (42) va yoʻq (999) buyurtma — bir xil 404.
func TestOrderOwner(t *testing.T) {
	h, _ := testAPI(t)
	dilnoza := login(t, h, dilnozaLogin)
	for id, want := range map[string]int{"41": 200, "42": 404, "999": 404} {
		rec := send(h, "GET", "/orders/"+id, "", dilnoza.Value)
		if rec.Code != want {
			t.Errorf("Dilnoza GET /orders/%s: %d, kutilgan %d", id, rec.Code, want)
		}
	}
}

// Begona buyurtmaning 404 javobi yoʻq buyurtmanikidan farq qilmasin (request_id'dan boshqa hamma narsa).
func TestForeignOrderLooksMissing(t *testing.T) {
	h, _ := testAPI(t)
	dilnoza := login(t, h, dilnozaLogin)
	foreign := send(h, "GET", "/orders/42", "", dilnoza.Value)
	missing := send(h, "GET", "/orders/999", "", dilnoza.Value)
	strip := func(rec *httptest.ResponseRecorder) string {
		body := rec.Body.String()
		return body[:strings.Index(body, `"request_id"`)]
	}
	if foreign.Code != missing.Code || strip(foreign) != strip(missing) {
		t.Errorf("begona: %d %s\nyoʻq: %d %s\nikkalasi bir xil boʻlishi kerak", foreign.Code, foreign.Body, missing.Code, missing.Body)
	}
	if strings.Contains(foreign.Body.String(), "Yunusobod") {
		t.Error("Jasurning manzili Dilnozaga chiqdi")
	}
	jasur := login(t, h, jasurLogin)
	if rec := send(h, "GET", "/orders/42", "", jasur.Value); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Yunusobod") {
		t.Errorf("Jasur oʻz buyurtmasi: %d %s, kutilgan 200", rec.Code, rec.Body)
	}
	if rec := send(h, "GET", "/orders/41", "", jasur.Value); rec.Code != http.StatusNotFound {
		t.Errorf("Jasur GET /orders/41 (Dilnozaniki): %d, kutilgan 404", rec.Code)
	}
	if rec := send(h, "GET", "/orders/42", "", ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("cookie'siz GET /orders/42: %d, kutilgan 401", rec.Code)
	}
}

// Buyurtma egasini session beradi: yangi buyurtma Dilnozaniki, body'dagi "user_id" esa rad etiladi.
func TestCreateOrderOwner(t *testing.T) {
	h, _ := testAPI(t)
	dilnoza := login(t, h, dilnozaLogin)
	rec := send(h, "POST", "/orders", `{"items":[{"product_id":3,"quantity":1}]}`, dilnoza.Value)
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"user_id":1`) {
		t.Fatalf("POST /orders: %d %s; kutilgan 201, user_id 1", rec.Code, rec.Body)
	}
	rec = send(h, "POST", "/orders", `{"user_id":2,"items":[{"product_id":3,"quantity":1}]}`, dilnoza.Value)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("body'da user_id: %d %s, kutilgan 400", rec.Code, rec.Body)
	}
}

// GB-6.06: /admin/users faqat admin uchun: xaridor 403, cookie'siz 401, admin 200.
func TestAdminUsers(t *testing.T) {
	h, _ := testAPI(t)
	tests := []struct {
		name, cookie string
		status       int
	}{
		{"cookie_yoq", "", 401},
		{"xaridor", login(t, h, dilnozaLogin).Value, 403},
		{"admin", login(t, h, adminLogin).Value, 200},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := send(h, "GET", "/admin/users", "", tc.cookie)
			if rec.Code != tc.status {
				t.Errorf("GET /admin/users: %d %s; kutilgan %d", rec.Code, rec.Body, tc.status)
			}
			if tc.status == 200 && !strings.Contains(rec.Body.String(), "jasur@bozor.test") {
				t.Errorf("admin roʻyxatida Jasur yoʻq: %s", rec.Body)
			}
			if strings.Contains(rec.Body.String(), "argon2id") {
				t.Error("javobda parol hash'i bor")
			}
		})
	}
}

// requireAdmin requireUser'dan tashqarida boʻlsa: currentUser panic → recoverPanic → 500 (jim oʻtkazish emas).
func TestRequireAdminOutside(t *testing.T) {
	pool := pgtest.New(t)
	a := &api{st: store.New(pool)}
	mux := http.NewServeMux()
	mux.Handle("GET /admin/users", requireAdmin(a.requireUser(http.HandlerFunc(a.listUsers))))
	h := recoverPanic(mux)
	c := login(t, newHandler(a.st), adminLogin)
	if rec := send(h, "GET", "/admin/users", "", c.Value); rec.Code != http.StatusInternalServerError {
		t.Errorf("tartib teskari: %d, kutilgan 500", rec.Code)
	}
}

// PATCH body'sida "note" kaliti boʻlmasa — 422, izoh jim oʻchib ketmaydi; Content-Type'siz — 415, lekin login'dan keyin.
func TestUpdateNoteInput(t *testing.T) {
	h, _ := testAPI(t)
	dilnoza := login(t, h, dilnozaLogin)
	if rec := send(h, "PATCH", "/orders/41", `{}`, dilnoza.Value); rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("PATCH {}: %d %s, kutilgan 422", rec.Code, rec.Body)
	}
	req := httptest.NewRequest("PATCH", "/orders/41", strings.NewReader(`{"note":"x"}`))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("cookie'siz, Content-Type'siz PATCH: %d, kutilgan 401 (avval kimligi)", rec.Code)
	}
}

// GB-6.06: katalogni faqat admin oʻzgartiradi (POST, DELETE /products); oʻqish hammaga ochiq.
func TestProductsAdminOnly(t *testing.T) {
	h, _ := testAPI(t)
	dilnoza := login(t, h, dilnozaLogin).Value
	admin := login(t, h, adminLogin).Value
	body := `{"name":"anor","price":1800000,"stock":5}`
	for _, c := range []struct {
		method, path, body, cookie string
		want                       int
	}{
		{"POST", "/products", body, "", 401},
		{"POST", "/products", body, dilnoza, 403},
		{"DELETE", "/products/1", "", dilnoza, 403},
		{"POST", "/products", body, admin, 201},
		{"DELETE", "/products/5", "", admin, 204},
		{"DELETE", "/products/1", "", admin, 409}, // olma buyurtmalarda bor
		{"GET", "/products/1", "", "", 200},
	} {
		if rec := send(h, c.method, c.path, c.body, c.cookie); rec.Code != c.want {
			t.Errorf("%s %s: %d %s, kutilgan %d", c.method, c.path, rec.Code, rec.Body, c.want)
		}
	}
}
