// Himoyalangan bozor'ning isboti (GB-6.08): spetsifikatsiya jadval shaklida va uni haqiqiy Postgres'da
// tekshiradigan testlar. Har test oʻz sxemasida (pgtest.New), seed.sql'dagi uch foydalanuvchi bilan:
// Dilnoza (xaridor, 41-buyurtma), Jasur (xaridor, 42-buyurtma), admin.
package main

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

// Jadvaldagi body'lar.
const (
	orderBody   = `{"items":[{"product_id":3,"quantity":1}]}`
	productBody = `{"name":"anor","price":1800000,"stock":5}`
	noteBody    = `{"note":"Jasurning manzili"}`
)

// rule — jadvalning bir qatori: kim, qaysi soʻrov, qaysi status kutiladi.
type rule struct {
	who                string // "anonim" (cookie'siz), "dilnoza", "jasur", "admin"
	method, path, body string
	want               int
}

// accessRules — kim nima qila oladi (README, «Spetsifikatsiya»). Ikki xaridor shart: bitta foydalanuvchi
// bilan yozilgan test begona buyurtmani hech qachon soʻramaydi va BOLA'ni koʻrmaydi.
var accessRules = []rule{
	{"anonim", "GET", "/products", "", 200},
	{"anonim", "GET", "/orders/41", "", 401},
	{"anonim", "POST", "/orders", orderBody, 401},
	{"anonim", "POST", "/admin/products", productBody, 401},
	{"dilnoza", "GET", "/orders/41", "", 200},
	{"dilnoza", "GET", "/orders/42", "", 404},
	{"jasur", "GET", "/orders/41", "", 404},
	{"jasur", "PATCH", "/orders/41", noteBody, 404},
	{"admin", "GET", "/orders/41", "", 404}, // /orders/{id} — «mening buyurtmam»: admin uchun ham
	{"dilnoza", "POST", "/admin/products", productBody, 403},
	{"jasur", "GET", "/admin/users", "", 403},
	{"dilnoza", "GET", "/admin/yoq", "", 403}, // yoʻq admin yoʻli ham 403: ichida nima borligi bilinmaydi
	{"admin", "GET", "/admin/yoq", "", 404},
	{"admin", "GET", "/admin/users", "", 200},
	{"admin", "POST", "/admin/products", productBody, 201},
	{"admin", "DELETE", "/admin/products/1", "", 409}, // olma buyurtmalarda bor
	{"anonim", "GET", "/orders", "", 401},
	{"anonim", "GET", "/admin/orders", "", 401},
	{"jasur", "GET", "/admin/orders", "", 403},
}

// TestAccess jadvalning har qatorini butun zanjir (newHandler) orqali yuboradi.
func TestAccess(t *testing.T) {
	h, _ := testAPI(t)
	cookies := map[string]string{
		"anonim":  "",
		"dilnoza": login(t, h, dilnozaLogin).Value,
		"jasur":   login(t, h, jasurLogin).Value,
		"admin":   login(t, h, adminLogin).Value,
	}
	for _, r := range accessRules {
		rec := send(h, r.method, r.path, r.body, cookies[r.who])
		if rec.Code != r.want {
			t.Errorf("%s %s %s: %d, kutilgan %d", r.who, r.method, r.path, rec.Code, r.want)
		}
	}
}

// Parol bazaga ochiq matn boʻlib tushmaydi: users'da faqat argon2id satri (GB-6.03), javobda esa hash ham yoʻq.
func TestPasswordNotStored(t *testing.T) {
	h, pool := testAPI(t)
	const pw = "anor-behi-yongoq-uzum-5"
	rec := send(h, "POST", "/users", `{"email":"aziza@bozor.test","password":"`+pw+`"}`, "")
	if rec.Code != http.StatusCreated || strings.Contains(rec.Body.String(), "argon2id") {
		t.Fatalf("POST /users: %d %s; kutilgan 201, javobda hash yoʻq", rec.Code, rec.Body)
	}
	var hash string
	err := pool.QueryRow(context.Background(),
		"SELECT password_hash FROM users WHERE email = $1", "aziza@bozor.test").Scan(&hash)
	if err != nil {
		t.Fatalf("users qatori oʻqilmadi: %v", err)
	}
	if strings.Contains(hash, pw) || !strings.HasPrefix(hash, "$argon2id$") {
		t.Errorf("password_hash = %q: parol ochiq turibdi yoki argon2id emas", hash)
	}
	login(t, h, `{"email":"aziza@bozor.test","password":"`+pw+`"}`) // hash'dan parolni tekshirish ishlaydi
}

// Logout session'ni serverda bekor qiladi: bazadagi qator oʻchadi, cookie nusxasi bilan keyingi soʻrov — 401.
func TestLogoutRevokes(t *testing.T) {
	h, pool := testAPI(t)
	c := login(t, h, dilnozaLogin).Value
	if rec := send(h, "GET", "/orders/41", "", c); rec.Code != http.StatusOK {
		t.Fatalf("login'dan keyin GET /orders/41: %d, kutilgan 200", rec.Code)
	}
	if rec := send(h, "POST", "/logout", "", c); rec.Code != http.StatusNoContent {
		t.Fatalf("POST /logout: %d, kutilgan 204", rec.Code)
	}
	if rec := send(h, "GET", "/orders/41", "", c); rec.Code != http.StatusUnauthorized {
		t.Errorf("logout'dan keyin oʻsha cookie bilan GET /orders/41: %d, kutilgan 401", rec.Code)
	}
	var n int
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM sessions").Scan(&n); err != nil || n != 0 {
		t.Errorf("sessions'da %d qator (%v), kutilgan 0", n, err)
	}
}

// Login'ga rate limit (GB-6.07): bitta IP'dan 5 urinish, keyin 429 — 6-urinish toʻgʻri parol bilan boʻlsa ham.
func TestLoginRateLimit(t *testing.T) {
	h, _ := testAPI(t)
	wrong := `{"email":"dilnoza@bozor.test","password":"taxmin-1"}`
	for i := 1; i <= 5; i++ {
		if rec := send(h, "POST", "/login", wrong, ""); rec.Code != http.StatusUnauthorized {
			t.Fatalf("%d-urinish: %d, kutilgan 401", i, rec.Code)
		}
	}
	rec := send(h, "POST", "/login", dilnozaLogin, "")
	if rec.Code != http.StatusTooManyRequests || rec.Header().Get("Retry-After") == "" {
		t.Errorf("6-urinish (toʻgʻri parol): %d, Retry-After %q; kutilgan 429 va Retry-After",
			rec.Code, rec.Header().Get("Retry-After"))
	}
}
