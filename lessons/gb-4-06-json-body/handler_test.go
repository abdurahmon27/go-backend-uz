package main

// Videodagi bosqichlarni tekshiradi (go test .): starterda ham oʻtadi. Soʻrovlar serversiz, mux orqali yuboriladi (GB-4.03).

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func post(path, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, strings.NewReader(body)))
	return rec
}

// errorBody javobni {"error": {…}} shaklida oʻqiydi.
func errorBody(t *testing.T, rec *httptest.ResponseRecorder) apiError {
	t.Helper()
	var got map[string]apiError
	dec := json.NewDecoder(rec.Body)
	if err := dec.Decode(&got); err != nil {
		t.Fatalf("xato javobi JSON emas (%v): body %q", err, rec.Body.String())
	}
	if dec.More() {
		t.Errorf("body'da ikkinchi JSON ham bor: writeError'dan keyin return unutilgan")
	}
	return got["error"]
}

func TestCreateProduct(t *testing.T) {
	before := len(catalog)
	rec := post("/products", `{"name":"choy","price":2500000,"stock":30}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status %d, kutilgan %d: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	var got Product
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("body JSON emas: %v", err)
	}
	if got.Name != "choy" || got.Price != 2500000 || got.Stock != 30 || got.ID <= 3 {
		t.Errorf("mahsulot %+v: nom, narx, qoldiq soʻrovdagidek, ID esa serverdan (3 dan katta) boʻlishi kerak", got)
	}
	if len(catalog) != before+1 {
		t.Errorf("katalogda %d ta mahsulot, kutilgan %d", len(catalog), before+1)
	}
}

func TestCreateProductXato(t *testing.T) {
	cases := []struct {
		name, body string
		status     int
		code       string
	}{
		{"buzuq JSON", `{"name":`, http.StatusBadRequest, "bad_json"},
		{"nomaʼlum kalit", `{"name":"choy","price":1,"rang":"qora"}`, http.StatusBadRequest, "bad_json"},
		{"ID ni mijoz bermaydi", `{"id":99,"name":"choy","price":1}`, http.StatusBadRequest, "bad_json"},
		{"birinchi qiymatdan keyin ikkinchi JSON", `{"name":"choy","price":1}{"x":1}`, http.StatusBadRequest, "bad_json"},
		{"birinchi qiymatdan keyin 2000 boʻsh joy", `{"name":"choy","price":1}` + strings.Repeat(" ", 2000), http.StatusRequestEntityTooLarge, "too_large"},
		{"juda katta body", `{"name":"` + strings.Repeat("a", 2000) + `"}`, http.StatusRequestEntityTooLarge, "too_large"},
		{"notoʻgʻri qiymatlar", `{"name":"","price":-5}`, http.StatusUnprocessableEntity, "validation_failed"},
	}
	for _, c := range cases {
		before := len(catalog)
		rec := post("/products", c.body)
		if rec.Code != c.status {
			t.Errorf("%s: status %d, kutilgan %d", c.name, rec.Code, c.status)
		}
		if got := rec.Result().Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("%s: Content-Type «%s», kutilgan «application/json»", c.name, got)
		}
		e := errorBody(t, rec)
		if e.Code != c.code {
			t.Errorf("%s: error.code «%s», kutilgan «%s»", c.name, e.Code, c.code)
		}
		if c.code != "validation_failed" && len(e.Fields) != 0 {
			t.Errorf("%s: fields faqat validatsiya xatosida boʻladi, hozir %v", c.name, e.Fields)
		}
		if len(catalog) != before {
			t.Errorf("%s: xatoli soʻrov katalogga mahsulot qoʻshdi (return unutilganmi?)", c.name)
		}
	}
}

func TestValidateBarchaMaydon(t *testing.T) {
	rec := post("/products", `{"name":"","price":-5,"stock":-1}`)
	f := errorBody(t, rec).Fields
	for _, k := range []string{"name", "price", "stock"} {
		if f[k] == "" {
			t.Errorf("fields[%q] boʻsh: validate birinchi xatoda toʻxtamay, hamma notoʻgʻri maydonni qaytarsin", k)
		}
	}
}
