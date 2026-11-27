package main

// VAZIFA tekshiruvi (TASK.md): requireJSON middleware'i va uning zanjirdagi joyi; 🟡 muxErrors. Bu faylni oʻzgartirmang.

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// post — body va berilgan Content-Type bilan POST (ct boʻsh boʻlsa — sarlavhasiz).
func post(h http.Handler, path, ct, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", path, strings.NewReader(body))
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// 1-qadam: requireJSON oʻzi (zanjirsiz).
func TestVazifaRequireJSON(t *testing.T) {
	// ok — zanjirning keyingi qismi oʻrnida: unga yetib kelgan soʻrov 204 oladi.
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	h := requireJSON(ok)
	order := `{"product_id":1,"quantity":1}`

	for _, ct := range []string{"", "text/plain", "application/x-www-form-urlencoded"} {
		rec := post(h, "/orders", ct, order)
		if rec.Code != http.StatusUnsupportedMediaType {
			t.Fatalf("requireJSON: POST, Content-Type %q → status %d (soʻrov next'ga oʻtib ketdi), kutilgan 415 (http.StatusUnsupportedMediaType). "+
				"requireJSON hali yozilmagan boʻlishi mumkin: TASK.md, 1-qadam", ct, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), `"code":"unsupported_media_type"`) {
			t.Errorf("requireJSON: 415 javobi bizning shaklda boʻlsin: writeError(…, apiError{Code: \"unsupported_media_type\", …}); hozir %s", rec.Body)
		}
	}
	for _, ct := range []string{"application/json", "application/json; charset=utf-8"} {
		if rec := post(h, "/orders", ct, order); rec.Code != http.StatusNoContent {
			t.Errorf("requireJSON: Content-Type %q — JSON, soʻrov next'ga oʻtishi kerak edi, hozir status %d %s "+
				"(maslahat: mime.ParseMediaType)", ct, rec.Code, rec.Body)
		}
	}
	if rec := send(h, "GET", "/products", ""); rec.Code != http.StatusNoContent {
		t.Errorf("requireJSON: GET soʻrovi body'siz, tekshirilmasligi kerak — status %d", rec.Code)
	}
}

// 1-qadam, zanjir: requireJSON newHandler ichida, requestID va logRequests ichkarisida.
func TestVazifaZanjir(t *testing.T) {
	buf := captureLog(t)
	h := newHandler(testStore(t))
	rec := post(h, "/products", "text/plain", `{"name":"anor","price":1800000,"stock":7}`)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("newHandler: POST /products, Content-Type text/plain → %d, kutilgan 415: requireJSON'ni ulang — newHandler zanjiriga yoki newMux'dagi POST route'larga (TASK.md, 1-qadam)", rec.Code)
	}
	id := rec.Header().Get("X-Request-Id")
	if id == "" {
		t.Errorf("415 javobida X-Request-Id yoʻq: requireJSON requestID'ning ichida boʻlsin")
	}
	if !strings.Contains(buf.String(), `POST "/products" 415`) {
		t.Errorf("415 javobi log'ga tushmadi (%q): requireJSON logRequests'ning ichida boʻlsin", buf)
	}
	if rec := send(h, "POST", "/products", `{"name":"anor","price":1800000,"stock":7}`); rec.Code != http.StatusCreated {
		t.Errorf("newHandler: JSON bilan POST /products → %d %s, kutilgan 201", rec.Code, rec.Body)
	}
}

// 2-qadam (🟡 ixtiyoriy): mux'ning 404 va 405 javoblari ham bizning JSON shaklida.
func TestVazifaMuxErrors(t *testing.T) {
	captureLog(t)
	h := newHandler(testStore(t))
	rec := send(h, "GET", "/narsalar", "")
	if strings.HasPrefix(rec.Header().Get("Content-Type"), "text/plain") {
		t.Skip("🟡 ixtiyoriy (TASK.md, 2-qadam): mux'ning 404'i hali oddiy matn — muxErrors yozilmagan")
	}
	if rec.Code != 404 || !strings.Contains(rec.Body.String(), `"code":"not_found"`) || strings.Contains(rec.Body.String(), "page not found") {
		t.Errorf("GET /narsalar: %d %q, kutilgan 404 va faqat bizning JSON xato (code not_found)", rec.Code, rec.Body)
	}
	rec = send(h, "PUT", "/products/1", "")
	if rec.Code != 405 || !strings.Contains(rec.Body.String(), `"code":"method_not_allowed"`) {
		t.Errorf("PUT /products/1: %d %q, kutilgan 405 va JSON xato (code method_not_allowed)", rec.Code, rec.Body)
	}
	if rec.Header().Get("Allow") == "" {
		t.Errorf("405 javobida Allow sarlavhasi saqlanishi kerak (RFC 9110)")
	}
	if rec := send(h, "GET", "/products/1", ""); rec.Code != 200 {
		t.Errorf("muxErrors oddiy route'ni buzdi: GET /products/1 → %d", rec.Code)
	}
}
