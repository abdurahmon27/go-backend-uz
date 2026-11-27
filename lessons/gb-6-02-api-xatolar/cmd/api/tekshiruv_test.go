package main

// Vazifa tekshiruvi (TASK.md). Bu faylni oʻzgartirmang.

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/abdurahmon27/bozor/internal/store"
)

// vazifaSend soʻrovni toʻliq zanjir (newHandler) orqali yuboradi va javobni problem sifatida oʻqiydi.
func vazifaSend(t *testing.T, method, path, body string) (int, problem, string) {
	t.Helper()
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	rec := send(newHandler(testStore(t)), method, path, body)
	var p problem
	if ct := rec.Header().Get("Content-Type"); ct == "application/problem+json" {
		json.Unmarshal(rec.Body.Bytes(), &p)
	}
	return rec.Code, p, rec.Body.String()
}

// 1-qadam: validate natijasi (fieldErrors) → 422, errors kengaytmasi bilan.
func TestVazifaFieldErrors(t *testing.T) {
	for _, c := range []struct{ path, body, field string }{
		{"/products", `{"name":"","price":0}`, "name"},
		{"/orders", `{"product_id":1,"quantity":0}`, "quantity"},
	} {
		status, p, body := vazifaSend(t, "POST", c.path, c.body)
		if status != 422 {
			t.Fatalf("POST %s %s: status %d, kutilgan 422. fieldErrors hozir problemFor'da default'ga tushyapti: "+
				"errors.As(err, &fe) case'ini qoʻshing (TASK.md, 1-qadam). Body: %s", c.path, c.body, status, body)
		}
		if p.Type != "/problems/validation" || p.Title != "Qiymatlar notoʻgʻri" || p.Errors[c.field] == "" {
			t.Errorf("POST %s: problem %+v: type «/problems/validation», title «Qiymatlar notoʻgʻri» va errors[%q] boʻlsin (TASK.md, 1-qadam)", c.path, p, c.field)
		}
	}
}

// 2-qadam: store qoidasi buzildi (ErrBadProduct, ErrBadQuantity) → 422, detail — productRules, store matni emas.
func TestVazifaStoreRules(t *testing.T) {
	long := strings.Repeat("a", 21)
	status, p, body := vazifaSend(t, "POST", "/products", `{"name":"`+long+`","price":100}`)
	if status != 422 {
		t.Fatalf("POST /products (21 belgili nom): status %d, kutilgan 422. store.ErrBadProduct va store.ErrBadQuantity "+
			"uchun case qoʻshing (TASK.md, 2-qadam). Body: %s", status, body)
	}
	if p.Type != "/problems/validation" || p.Detail != productRules {
		t.Errorf("problem %+v: type «/problems/validation», detail esa productRules boʻlsin (TASK.md, 2-qadam)", p)
	}
	if strings.Contains(body, "eng koʻpi") || strings.Contains(body, "mahsulot notoʻgʻri") {
		t.Errorf("body'da store xatosining matni bor: %s. detail'ga err.Error() emas, productRules yozing (TASK.md, 2-qadam)", body)
	}
	// ErrBadQuantity HTTP orqali kelmaydi (validateOrder uni oldinroq ushlaydi), shuning uchun problemFor'ning oʻzini chaqiramiz.
	if q := problemFor(fmt.Errorf("Sell: %w", store.ErrBadQuantity)); q.Status != 422 || q.Detail != productRules {
		t.Errorf("problemFor(…ErrBadQuantity): status %d, detail %q — kutilgan 422 va productRules: "+
			"ErrBadQuantity'ni ErrBadProduct bilan bitta case'ga qoʻshing (TASK.md, 2-qadam)", q.Status, q.Detail)
	}
}

// 3-qadam: juda katta body → 413, 400 emas. Case tartibi: MaxBytesError errBadJSON'dan oldin.
func TestVazifaTooLarge(t *testing.T) {
	status, p, body := vazifaSend(t, "POST", "/products", `{"name":"`+strings.Repeat("a", 2000)+`"}`)
	if status != 413 {
		t.Fatalf("POST /products (2000 baytli body): status %d, kutilgan 413. *http.MaxBytesError uchun case qoʻshing va "+
			"uni errBadJSON case'idan OLDIN qoʻying: switch birinchi mos case'ni oladi (TASK.md, 3-qadam). Body: %s", status, body)
	}
	if p.Title != "Request Entity Too Large" || !strings.Contains(p.Detail, "1024") {
		t.Errorf("problem %+v: newProblem(http.StatusRequestEntityTooLarge, …) va detail'da limit (1024) boʻlsin (TASK.md, 3-qadam)", p)
	}
}
