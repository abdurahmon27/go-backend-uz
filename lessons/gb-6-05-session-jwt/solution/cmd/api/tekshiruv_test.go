// VAZIFA tekshiruvi (TASK.md): logout JWT'ni ham bekor qilsin — jti denylist.
package main

import (
	"net/http"
	"testing"
)

func TestLogoutJWTniBekorQiladi(t *testing.T) {
	h, _, _ := testAPI(t)
	_, stolen := login(t, h, dilnozaLogin) // shu tokenning nusxasi "oʻgʻirlangan"
	if rec := send(h, "GET", "/me", "", auth{bearer: stolen}); rec.Code != http.StatusOK {
		t.Fatalf("logout'dan oldin JWT bilan /me: %d, kutilgan 200", rec.Code)
	}
	if rec := send(h, "POST", "/logout", "", auth{bearer: stolen}); rec.Code != http.StatusNoContent {
		t.Fatalf("Bearer bilan POST /logout: status %d, kutilgan 204: %s", rec.Code, rec.Body)
	}
	if rec := send(h, "GET", "/me", "", auth{bearer: stolen}); rec.Code != http.StatusUnauthorized {
		t.Fatalf("logout'dan keyin oʻsha JWT bilan /me: status %d, kutilgan 401 — jti denylist'ga tushmadi yoki authenticate uni tekshirmayapti (TASK.md, 2–3-qadam)", rec.Code)
	}
	// Faqat shu token oʻldi: yangi login'ning tokeni ishlaydi.
	_, fresh := login(t, h, dilnozaLogin)
	if rec := send(h, "GET", "/me", "", auth{bearer: fresh}); rec.Code != http.StatusOK {
		t.Errorf("yangi login'dan keyin yangi JWT bilan /me: %d, kutilgan 200 — denylist faqat bitta jti'ni bekor qilsin", rec.Code)
	}
}
