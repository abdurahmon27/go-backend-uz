// API testlari: haqiqiy handler'lar, haqiqiy middleware zanjiri, xotiradagi store (Seed: ikki foydalanuvchi).
// Server ishga tushmaydi: httptest soʻrovni toʻgʻridan-toʻgʻri newHandler'ga beradi (GB-4.03).
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/abdurahmon27/bozor/internal/store"
	"github.com/abdurahmon27/bozor/internal/token"
)

const (
	dilnozaLogin = `{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}`
	jasurLogin   = `{"email":"jasur@bozor.test","password":"choy-non-qand-27"}`
)

// TestMain: server log'i test chiqishini toʻldirmasin.
func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

// auth — soʻrovga nima qoʻshiladi: session cookie'si yoki Bearer JWT (ikkalasi ham boʻsh boʻlishi mumkin).
type auth struct{ cookie, bearer string }

func send(h http.Handler, method, path, body string, a auth) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if a.cookie != "" {
		req.AddCookie(&http.Cookie{Name: sessionCookie, Value: a.cookie})
	}
	if a.bearer != "" {
		req.Header.Set("Authorization", "Bearer "+a.bearer)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// testAPI — yangi store (Seed bilan), yangi kalit va shu store'dagi butun zanjir.
func testAPI(t *testing.T) (http.Handler, *store.Store, []byte) {
	t.Helper()
	st := store.New()
	if err := st.Seed(context.Background()); err != nil {
		t.Fatal(err)
	}
	key := token.NewKey()
	return newHandler(st, key), st, key
}

// login body bilan POST /login yuboradi va session tokenini (cookie) hamda JWT'ni (body) qaytaradi.
func login(t *testing.T, h http.Handler, body string) (sess, jwt string) {
	t.Helper()
	rec := send(h, "POST", "/login", body, auth{})
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /login: status %d, kutilgan 200: %s", rec.Code, rec.Body)
	}
	var out loginOutput
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil || strings.Count(out.Token, ".") != 2 {
		t.Fatalf("POST /login javobida JWT yoʻq: token=%q err=%v", out.Token, err)
	}
	for _, c := range rec.Result().Cookies() {
		if c.Name == sessionCookie {
			return c.Value, out.Token
		}
	}
	t.Fatalf("POST /login javobida %q cookie yoʻq", sessionCookie)
	return "", ""
}

func TestMeIkkiYol(t *testing.T) {
	h, _, _ := testAPI(t)
	sess, jwt := login(t, h, dilnozaLogin)
	if rec := send(h, "GET", "/me", "", auth{cookie: sess}); rec.Code != 200 || strings.TrimSpace(rec.Body.String()) != `{"id":1,"email":"dilnoza@bozor.test"}` {
		t.Errorf("session bilan /me: %d %s", rec.Code, rec.Body)
	}
	// JWT yoʻli faqat tokendagini biladi: sub → id. Email tokenda yoʻq.
	if rec := send(h, "GET", "/me", "", auth{bearer: jwt}); rec.Code != 200 || strings.TrimSpace(rec.Body.String()) != `{"id":1}` {
		t.Errorf("JWT bilan /me: %d %s", rec.Code, rec.Body)
	}
	if rec := send(h, "GET", "/me", "", auth{}); rec.Code != http.StatusUnauthorized {
		t.Errorf("hech narsasiz /me: %d, kutilgan 401", rec.Code)
	}
}

func TestYaroqsizJWT(t *testing.T) {
	h, _, _ := testAPI(t)
	_, jwt := login(t, h, dilnozaLogin)
	parts := strings.Split(jwt, ".")
	bad := map[string]string{
		"imzo buzilgan":  parts[0] + "." + parts[1] + ".AAAA" + parts[2][4:],
		"boshqa kalit":   mustIssue(t, 1, token.NewKey()),
		"oddiy matn":     "salom",
		"session tokeni": "8rgo2ZzSo05nCwhmL1r5PyL7n1M7WPm_C0e2_Q4UkyI",
	}
	for name, tok := range bad {
		rec := send(h, "GET", "/me", "", auth{bearer: tok})
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("%s: /me status %d, kutilgan 401", name, rec.Code)
		}
		if got := rec.Header().Get("WWW-Authenticate"); got != `Bearer error="invalid_token"` {
			t.Errorf("%s: WWW-Authenticate = %q", name, got)
		}
	}
}

func mustIssue(t *testing.T, id int64, key []byte) string {
	t.Helper()
	s, err := token.Issue(id, jwtTTL, key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// Session'ni logout darhol oʻldiradi: oʻsha cookie bilan keyingi soʻrov 401 (GB-6.04).
func TestSessionLogout(t *testing.T) {
	h, _, _ := testAPI(t)
	sess, _ := login(t, h, dilnozaLogin)
	if rec := send(h, "POST", "/logout", "", auth{cookie: sess}); rec.Code != http.StatusNoContent {
		t.Fatalf("POST /logout: %d, kutilgan 204", rec.Code)
	}
	if rec := send(h, "GET", "/me", "", auth{cookie: sess}); rec.Code != http.StatusUnauthorized {
		t.Errorf("logout'dan keyin eski cookie bilan /me: %d, kutilgan 401", rec.Code)
	}
}

// Server qayta ishga tushdi = yangi kalit: eski JWT'lar hammasi oʻladi (bitta foydalanuvchiniki emas — hammaniki).
func TestYangiKalit(t *testing.T) {
	h, st, _ := testAPI(t)
	_, jwt := login(t, h, jasurLogin)
	h2 := newHandler(st, token.NewKey())
	if rec := send(h2, "GET", "/me", "", auth{bearer: jwt}); rec.Code != http.StatusUnauthorized {
		t.Errorf("yangi kalit bilan eski JWT: %d, kutilgan 401", rec.Code)
	}
}

// GB-6.03/6.04'dan qolganlar ham ishlaydi: roʻyxatdan oʻtish va "hamma qurilmadan chiqish".
func TestUsersVaLogoutAll(t *testing.T) {
	h, _, _ := testAPI(t)
	if rec := send(h, "POST", "/users", `{"email":"malika@bozor.test","password":"anor-olma-uzum-27"}`, auth{}); rec.Code != http.StatusCreated {
		t.Fatalf("POST /users: %d %s, kutilgan 201", rec.Code, rec.Body)
	}
	telefon, jwt := login(t, h, `{"email":"malika@bozor.test","password":"anor-olma-uzum-27"}`)
	noutbuk, _ := login(t, h, `{"email":"malika@bozor.test","password":"anor-olma-uzum-27"}`)
	if rec := send(h, "POST", "/logout-all", "", auth{cookie: telefon}); rec.Code != http.StatusNoContent {
		t.Fatalf("POST /logout-all: %d, kutilgan 204", rec.Code)
	}
	if rec := send(h, "GET", "/me", "", auth{cookie: noutbuk}); rec.Code != http.StatusUnauthorized {
		t.Errorf("logout-all'dan keyin boshqa session bilan /me: %d, kutilgan 401", rec.Code)
	}
	// logout-all faqat session yozuvlarini oʻchiradi: JWT'da oʻchiradigan yozuv yoʻq.
	if rec := send(h, "GET", "/me", "", auth{bearer: jwt}); rec.Code != http.StatusOK {
		t.Errorf("logout-all'dan keyin JWT bilan /me: %d, kutilgan 200 (exp'gacha yaroqli)", rec.Code)
	}
}

// logout requireUser'siz (GB-6.04): yaroqsiz Bearer bilan ham 204, server yiqilmaydi.
func TestLogoutYaroqsizBearer(t *testing.T) {
	h, _, _ := testAPI(t)
	if rec := send(h, "POST", "/logout", "", auth{bearer: "salom"}); rec.Code != http.StatusNoContent {
		t.Errorf("yaroqsiz Bearer bilan POST /logout: %d %s, kutilgan 204", rec.Code, rec.Body)
	}
}

// Sxema nomida katta-kichik harf farqi yoʻq (RFC 9110 §11.1), token oldida bir necha boʻsh joy boʻlishi mumkin (RFC 6750 §2.1).
func TestBearerSxemasi(t *testing.T) {
	h, _, _ := testAPI(t)
	_, jwt := login(t, h, dilnozaLogin)
	for _, v := range []string{"Bearer " + jwt, "bearer " + jwt, "BEARER  " + jwt} {
		req := httptest.NewRequest("GET", "/me", nil)
		req.Header.Set("Authorization", v)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("Authorization %q: /me %d, kutilgan 200", v[:8], rec.Code)
		}
	}
}

// logout har credential'ni alohida bekor qiladi: faqat Bearer bilan logout cookie'dagi session'ga tegmaydi.
func TestLogoutHarBiriAlohida(t *testing.T) {
	h, _, _ := testAPI(t)
	sess, jwt := login(t, h, dilnozaLogin)
	send(h, "POST", "/logout", "", auth{cookie: sess})
	if rec := send(h, "GET", "/me", "", auth{bearer: jwt}); rec.Code != http.StatusOK {
		t.Errorf("faqat cookie bilan logout'dan keyin JWT: %d, kutilgan 200 (JWT'da yozuv yoʻq)", rec.Code)
	}
}
