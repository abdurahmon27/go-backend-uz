package main

// VAZIFA tekshiruvi (TASK.md, 2- va 3-qadam): DELETE /products/{id} va uning oʻz testi. Bu faylni oʻzgartirmang.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVazifaDeleteRoute(t *testing.T) {
	h := newMux(testStore(t))

	rec := send(h, "DELETE", "/products/2", "")
	if rec.Code == 405 {
		t.Fatalf("DELETE /products/2: 405 — route yoʻq: newMux'ga \"DELETE /products/{id}\" pattern'ini a.deleteProduct bilan qoʻshing (TASK.md, 2-qadam)")
	}
	if rec.Code != 204 {
		t.Fatalf("DELETE /products/2: status %d, kutilgan 204 No Content: %s", rec.Code, rec.Body)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("DELETE /products/2: 204 javobida body boʻlmaydi, hozir %q", rec.Body)
	}
	if rec := send(h, "GET", "/products/2", ""); rec.Code != 404 {
		t.Errorf("DELETE'dan keyin GET /products/2: status %d, kutilgan 404 (mahsulot oʻchgan boʻlishi kerak)", rec.Code)
	}
	if rec := send(h, "GET", "/products", ""); strings.Contains(rec.Body.String(), `"name":"non"`) || !strings.Contains(rec.Body.String(), `"name":"choy"`) {
		t.Errorf("DELETE'dan keyin GET /products: %s — non boʻlmasligi, qolganlari boʻlishi kerak", rec.Body)
	}

	rec = send(h, "DELETE", "/products/2", "")
	if rec.Code != 404 || !strings.Contains(rec.Body.String(), `"code":"not_found"`) {
		t.Errorf("DELETE /products/2 ikkinchi marta: %d %s, kutilgan 404 va JSON xato (writeStoreError)", rec.Code, rec.Body)
	}
	rec = send(h, "DELETE", "/products/abc", "")
	if rec.Code != 400 || !strings.Contains(rec.Body.String(), `"code":"bad_id"`) {
		t.Errorf("DELETE /products/abc: %d %s, kutilgan 400 va bad_id (pathID)", rec.Code, rec.Body)
	}
}

// 3-qadam: cmd/api'dagi test fayllaridan birida DELETE uchun oʻzingiz yozgan test boʻlsin.
func TestVazifaOzTest(t *testing.T) {
	files, err := filepath.Glob("*_test.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if f == "tekshiruv_test.go" {
			continue
		}
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		s := string(src)
		i := strings.Index(s, "func TestDeleteProduct(")
		if i < 0 {
			continue
		}
		body := s[i:]
		if j := strings.Index(body[1:], "\nfunc "); j >= 0 {
			body = body[:j+1]
		}
		if !strings.Contains(body, `"DELETE"`) && !strings.Contains(body, "http.MethodDelete") {
			t.Errorf("%s: TestDeleteProduct DELETE soʻrovini yubormayapti: send(h, \"DELETE\", …) yoki http.MethodDelete", f)
		}
		if !strings.Contains(body, "204") && !strings.Contains(body, "StatusNoContent") {
			t.Errorf("%s: TestDeleteProduct 204 (http.StatusNoContent) statusini tekshirmayapti", f)
		}
		return
	}
	t.Fatal("cmd/api'da func TestDeleteProduct(t *testing.T) yoʻq: DELETE /products/{id} uchun oʻz testingizni yozing (TASK.md, 3-qadam)")
}
