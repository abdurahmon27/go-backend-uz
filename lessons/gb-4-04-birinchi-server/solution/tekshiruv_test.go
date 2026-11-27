package main

// Vazifani tekshiradi (go test .). Oʻzgartirmang.
// newServer("127.0.0.1:0") qaytargan *http.Server'ning Handler'ini httptest.NewServer bilan boʻsh portda haqiqiy server
// qilib ochadi (GB-4.03) va unga haqiqiy HTTP soʻrovlar yuboradi; ReadHeaderTimeout maydonini esa toʻgʻridan-toʻgʻri oʻqiydi.
// 127.0.0.1 — shu kompyuter, :0 — portni operatsion tizim tanlaydi (Addr faqat oʻqiladi, server bu manzilda ochilmaydi).

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// get — test serveriga GET yuboradi va status bilan body'ni qaytaradi (Timeout bilan, GB-4.02).
func get(t *testing.T, base, path string) (int, string) {
	t.Helper()
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(base + path)
	if err != nil {
		t.Fatalf("GET %s: soʻrov bajarilmadi: %v", path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("GET %s: body oʻqilmadi: %v", path, err)
	}
	return resp.StatusCode, string(body)
}

func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := newServer("127.0.0.1:0")
	if srv == nil {
		t.Fatal("newServer nil qaytardi: mux'ni yasab, &http.Server{…} qaytaring (TASK.md, 2-qadam)")
	}
	if srv.Addr != "127.0.0.1:0" {
		t.Errorf("Addr «%s», kutilgan «127.0.0.1:0»: Addr'ga newServer'ning addr parametrini bering", srv.Addr)
	}
	if srv.Handler == nil {
		t.Fatal("Handler nil: server DefaultServeMux'ni ishlatadi, unda bitta ham yoʻl yoʻq (hamma soʻrov 404). Handler: mux bering (TASK.md, 2-qadam)")
	}
	ts := httptest.NewServer(srv.Handler)
	t.Cleanup(ts.Close)
	return ts
}

func TestVazifaYollar(t *testing.T) {
	ts := testServer(t)
	if code, body := get(t, ts.URL, "/health"); code != http.StatusOK || body != "ok\n" {
		t.Errorf("GET /health: %d %q, kutilgan 200 \"ok\\n\": mux'da \"/health\" → healthHandler boʻlsin (TASK.md, 1-qadam)", code, body)
	}
	code, body := get(t, ts.URL, "/products")
	var got []Product
	if code != http.StatusOK || json.Unmarshal([]byte(body), &got) != nil || len(got) != len(catalog) {
		t.Errorf("GET /products: %d %q, kutilgan 200 va %d ta mahsulotli JSON massiv: mux'da \"/products\" → catalogHandler boʻlsin (TASK.md, 1-qadam)", code, body, len(catalog))
	}
}

func TestVazifaOlma(t *testing.T) {
	ts := testServer(t)
	code, body := get(t, ts.URL, "/products/1")
	if code != http.StatusOK {
		t.Fatalf("GET /products/1: status %d, kutilgan 200: productHandler'ni \"/products/1\" patterniga ulang (TASK.md, 1-qadam)", code)
	}
	var got Product
	if err := json.Unmarshal([]byte(body), &got); err != nil || got != olma {
		t.Errorf("GET /products/1: body %q, kutilgan olma JSON'i: bu yoʻl productHandler'ga ulanishi kerak (TASK.md, 1-qadam)", body)
	}
	if code, _ := get(t, ts.URL, "/products/2"); code != http.StatusNotFound {
		t.Errorf("GET /products/2: status %d, kutilgan 404: faqat aynan \"/products/1\" patternini qoʻshing, oxiri slesh bilan emas (TASK.md, 1-qadam)", code)
	}
}

func TestVazifaTimeout(t *testing.T) {
	srv := newServer("127.0.0.1:0")
	if srv == nil {
		t.Fatal("newServer nil qaytardi: mux'ni yasab, &http.Server{…} qaytaring (TASK.md, 2-qadam)")
	}
	switch d := srv.ReadHeaderTimeout; {
	case d <= 0:
		t.Errorf("ReadHeaderTimeout %v: server sarlavhalarni cheksiz kutadi. ReadHeaderTimeout: 5 * time.Second bering (TASK.md, 2-qadam)", d)
	case d < time.Second:
		t.Errorf("ReadHeaderTimeout %v: bu soniya emas. time.Duration nanosoniyalarda: 5 = 5 nanosoniya; 5 * time.Second yozing (TASK.md, 2-qadam)", d)
	case d > 30*time.Second:
		t.Errorf("ReadHeaderTimeout %v juda katta: sekin mijoz ulanishni shuncha ushlab turadi; 5 * time.Second bering (TASK.md, 2-qadam)", d)
	}
}
