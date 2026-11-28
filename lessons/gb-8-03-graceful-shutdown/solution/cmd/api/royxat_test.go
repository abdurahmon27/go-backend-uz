// GB-6.08 vazifasi (GET /orders va GET /admin/orders) — endi oddiy test: bu route'lar va ularning jadvaldagi
// qatorlari keyingi darslarda ham ishlashi kerak. seed.sql: 41 — Dilnozaniki (jami 36 000 soʻm), 42 — Jasurniki (60 000 soʻm).
package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestMyOrders(t *testing.T) {
	h, _ := testAPI(t)
	malika := `{"email":"malika@bozor.test","password":"oʻrik-uzum-qovun-ariq-9"}`
	if rec := send(h, "POST", "/users", malika, ""); rec.Code != http.StatusCreated {
		t.Fatalf("POST /users (Malika): %d %s", rec.Code, rec.Body)
	}
	for _, c := range []struct{ who, login, want string }{
		{"Dilnoza", dilnozaLogin, `[{"id":41,"user_id":1,"status":"yangi","total":3600000}]`},
		{"Jasur", jasurLogin, `[{"id":42,"user_id":2,"status":"yangi","total":6000000}]`},
		{"Malika (buyurtmasiz)", malika, `[]`},
	} {
		rec := send(h, "GET", "/orders", "", login(t, h, c.login).Value)
		got := strings.TrimSpace(rec.Body.String())
		if rec.Code != http.StatusOK || got != c.want {
			t.Errorf("%s GET /orders: %d %s\n  kutilgan 200 %s (TASK.md, 1–2-qadam: faqat oʻz buyurtmalari, boʻlmasa [])",
				c.who, rec.Code, got, c.want)
		}
	}
	if rec := send(h, "GET", "/orders", "", ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("cookie'siz GET /orders: %d, kutilgan 401 (TASK.md, 2-qadam: requireUser)", rec.Code)
	}
}

func TestAdminOrders(t *testing.T) {
	h, _ := testAPI(t)
	want := `[{"id":41,"user_id":1,"status":"yangi","total":3600000},{"id":42,"user_id":2,"status":"yangi","total":6000000}]`
	rec := send(h, "GET", "/admin/orders", "", login(t, h, adminLogin).Value)
	if got := strings.TrimSpace(rec.Body.String()); rec.Code != http.StatusOK || got != want {
		t.Errorf("admin GET /admin/orders: %d %s\n  kutilgan 200 %s (TASK.md, 1-qadam AllOrders, 2-qadam admin mux)", rec.Code, got, want)
	}
	if rec := send(h, "GET", "/admin/orders", "", login(t, h, dilnozaLogin).Value); rec.Code != http.StatusForbidden {
		t.Errorf("Dilnoza GET /admin/orders: %d, kutilgan 403 (TASK.md, 2-qadam: route admin mux'da boʻlsin, asosiy mux'da emas)", rec.Code)
	}
}

// Jadvalda yangi route'larning qatorlari bormi (TASK.md, 3-qadam)? TestAccess ularni ishga tushiradi.
func TestAccessRulesCoverTask(t *testing.T) {
	need := []struct {
		who, method, path string
		want              int
	}{
		{"anonim", "GET", "/orders", 401},
		{"anonim", "GET", "/admin/orders", 401},
		{"xaridor", "GET", "/admin/orders", 403}, // Dilnoza yoki Jasur
	}
	for _, n := range need {
		found := false
		for _, r := range accessRules {
			buyer := r.who == "dilnoza" || r.who == "jasur"
			if (r.who == n.who || n.who == "xaridor" && buyer) && r.method == n.method && r.path == n.path && r.want == n.want {
				found = true
			}
		}
		if !found {
			t.Errorf("accessRules'da qator yoʻq: %s %s %s → %d (TASK.md, 3-qadam: himoya_test.go jadvali)", n.who, n.method, n.path, n.want)
		}
	}
}
