// Vazifa tekshiruvi (TASK.md): POST /logout-all — "hamma qurilmadan chiqish". Bu faylni oʻzgartirmang.
package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestLogoutAll(t *testing.T) {
	h, _ := testAPI(t)
	phone := login(t, h, dilnozaLogin)  // Dilnoza telefonda
	laptop := login(t, h, dilnozaLogin) // va noutbukda
	jasur := login(t, h, jasurLogin)    // boshqa foydalanuvchi

	rec := send(h, "POST", "/logout-all", "", phone.Value)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("POST /logout-all: status %d, kutilgan 204 (TASK.md, 2-qadam): %s", rec.Code, rec.Body)
	}
	if got := rec.Header().Get("Set-Cookie"); !strings.Contains(got, "Max-Age=0") {
		t.Errorf("POST /logout-all javobida Set-Cookie %q: cookie oʻchirilsin, clearSessionCookie (TASK.md, 2-qadam)", got)
	}
	for name, c := range map[string]*http.Cookie{"telefon": phone, "noutbuk": laptop} {
		if rec := send(h, "GET", "/me", "", c.Value); rec.Code != http.StatusUnauthorized {
			t.Errorf("logout-all'dan keyin Dilnozaning %s tokeni: GET /me %d, kutilgan 401 (TASK.md, 1-qadam: hamma session'lari oʻchsin)", name, rec.Code)
		}
	}
	if rec := send(h, "GET", "/me", "", jasur.Value); rec.Code != http.StatusOK {
		t.Errorf("Jasurning session'i ham oʻchdi: GET /me %d, kutilgan 200 (TASK.md, 1-qadam: faqat shu foydalanuvchiniki)", rec.Code)
	}
	if rec := send(h, "POST", "/logout-all", "", ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("cookie'siz POST /logout-all: %d, kutilgan 401 (TASK.md, 2-qadam: route requireUser bilan)", rec.Code)
	}
}
