package main

// Vazifa tekshiruvi (TASK.md). Bu faylni oʻzgartirmang.

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// 1–2-qadam: login urinishlari email boʻyicha ham cheklanadi — hujumchi IP'larni almashtirsa ham.
func TestLoginEmailLimit(t *testing.T) {
	h := newTestHandler(t)
	emails := []string{"dilnoza@bozor.test", "Dilnoza@bozor.test", "DILNOZA@BOZOR.TEST"}
	for i := 1; i <= 10; i++ {
		body := fmt.Sprintf(`{"email":%q,"password":"taxmin-%d"}`, emails[i%3], i)
		ip := fmt.Sprintf("10.0.0.%d", i) // har urinish boshqa IP'dan: IP limiti ishlamaydi
		if w := do(h, "POST", "/login", body, ip, nil); w.Code != http.StatusUnauthorized {
			t.Fatalf("%d-urinish (%s, IP %s): status %d, kutilgan 401", i, emails[i%3], ip, w.Code)
		}
	}
	w := do(h, "POST", "/login", `{"email":"dilnoza@bozor.test","password":"taxmin-11"}`, "10.0.0.11", nil)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("11-urinish (dilnoza@bozor.test, yangi IP 10.0.0.11): status %d, kutilgan 429 "+
			"(TASK.md: email boʻyicha limit — 10 ta urinish, harf kattaligi farq qilmaydi)", w.Code)
	}
	if s, err := strconv.Atoi(w.Header().Get("Retry-After")); err != nil || s < 1 {
		t.Errorf("429 javobida Retry-After %q, kutilgan musbat soniya (tooManyRequests'dan foydalaning)",
			w.Header().Get("Retry-After"))
	}
	if w := do(h, "POST", "/login", `{"email":"jasur@bozor.test","password":"taxmin-1"}`, "10.0.0.12", nil); w.Code != http.StatusUnauthorized {
		t.Errorf("jasur@bozor.test: status %d, kutilgan 401 — limit faqat Dilnozaning email'iga tegishli", w.Code)
	}
}

// 🟡 Ixtiyoriy (3-qadam): sweep — eski chelaklarni tozalash.
func TestSweep(t *testing.T) {
	l := newLimiters(12*time.Second, 5)
	s, ok := any(l).(interface{ sweep(time.Time) int })
	if !ok {
		t.Skip("🟡 ixtiyoriy: (*limiters).sweep(now time.Time) int hali yozilmagan (TASK.md, 3-qadam)")
	}
	l.get("192.0.2.1")  // toʻla chelak: hech kim ishlatmagan
	l.wait("192.0.2.2") // bitta token olindi
	l.wait("192.0.2.3") // bitta token olindi
	if n := s.sweep(time.Now()); n != 1 {
		t.Errorf("sweep(hozir) = %d, kutilgan 1: faqat toʻla chelak (192.0.2.1) oʻchadi", n)
	}
	if n := s.sweep(time.Now().Add(time.Minute)); n != 2 {
		t.Errorf("sweep(bir daqiqadan keyin) = %d, kutilgan 2: chelaklar qayta toʻlgan", n)
	}
	if len(l.m) != 0 {
		t.Errorf("sweep'dan keyin map'da %d ta kalit, kutilgan 0", len(l.m))
	}
}
