package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang. Testlar haqiqiy internetga chiqmaydi: har test oʻz soxta serverini ochadi.
//   httptest.NewServer(…) — shu kompyuterda vaqtincha server ochadi, srv.URL — uning manzili;
//   unga http.HandlerFunc(func(w, r) {…}) beriladi — "soʻrov kelsa shu funksiya javob beradi".
//   Ularni GB-4.03 da tushuntiramiz; hozir faqat shu: test serveri kerakli javobni beradi, siz esa
//   fetchRate'ni oʻsha serverga qaratib tekshirasiz
//   w.WriteHeader(code) — javob statusi; fmt.Fprint(w, body) — javob body'si
//   time.Sleep(d) — d vaqt kutadi (server "javob bermay turishi" uchun)
//   t.Cleanup(srv.Close) — test tugagach serverni yopadi
//   errors.As (GB-2.08) — oʻralgan xato ichidan *StatusError turini topadi
//   strings.Contains(s, sub) — s ichida sub bormi; strings.Repeat("x", n) — "x"ni n marta takrorlaydi

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const okBody = `{"result":"success","base_code":"USD","rates":{"EUR":0.91,"UZS":12000.5}}`

// newServer status va body bilan javob beradigan test serverini ochadi.
func newServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// Videodagi fetchRate: 200 va toʻgʻri JSON → kurs (bu test boshidanoq oʻtadi).
func TestVideo(t *testing.T) {
	srv := newServer(t, http.StatusOK, okBody)
	client := &http.Client{Timeout: time.Second}
	got, err := fetchRate(client, srv.URL, "UZS")
	if err != nil {
		t.Fatalf("fetchRate xato qaytardi: %v", err)
	}
	if got != 12000.5 {
		t.Errorf("fetchRate(…, \"UZS\") = %v; kutilgan 12000.5", got)
	}
}

// Videodagi Timeout: server javob bermay tursa, client kutib qolmaydi (bu test ham boshidanoq oʻtadi).
func TestVideoTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
		fmt.Fprint(w, okBody)
	}))
	t.Cleanup(srv.Close)
	client := &http.Client{Timeout: 50 * time.Millisecond}
	_, err := fetchRate(client, srv.URL, "UZS")
	if err == nil {
		t.Fatal("server kech javob berdi, lekin fetchRate xato qaytarmadi: client'da Timeout boʻlishi kerak")
	}
	if !strings.Contains(err.Error(), "Client.Timeout") {
		t.Errorf("xato «%v»: Client.Timeout haqida boʻlishi kerak edi", err)
	}
}

// 1-qadam (a): javobda soʻralgan valyuta yoʻq.
func TestVazifaKodYoq(t *testing.T) {
	srv := newServer(t, http.StatusOK, okBody)
	client := &http.Client{Timeout: time.Second}
	rate, err := fetchRate(client, srv.URL, "XXX")
	if err == nil {
		t.Fatalf("«XXX» kodi javobda yoʻq: fetchRate xato qaytarishi kerak, hozir %v va nil — "+
			"map'da yoʻq kalit nol beradi, vergulli shaklni ishlating (TASK.md, 1-qadam)", rate)
	}
	if !strings.Contains(err.Error(), "XXX") {
		t.Errorf("xato matni «%v»: kod «XXX» matnda boʻlishi kerak (TASK.md, 1-qadam)", err)
	}
}

// 1-qadam (b): status 200, lekin API body'da "result":"error" deydi.
func TestVazifaResultXato(t *testing.T) {
	srv := newServer(t, http.StatusOK, `{"result":"error","error-type":"unsupported-code"}`)
	client := &http.Client{Timeout: time.Second}
	rate, err := fetchRate(client, srv.URL, "UZS")
	if err == nil {
		t.Fatalf("API \"result\":\"error\" dedi (status 200), fetchRate esa %v va nil qaytardi; "+
			"kerak: xato, matnida API'ning error-type qiymati (TASK.md, 1-qadam)", rate)
	}
	if !strings.Contains(err.Error(), "unsupported-code") {
		t.Errorf("xato matni «%v»: «unsupported-code» boʻlishi kerak (TASK.md, 1-qadam)", err)
	}
}

// 2-qadam: 200 dan boshqa status — *StatusError, chaqiruvchi kodini bilib oladi.
func TestVazifaStatusError(t *testing.T) {
	for _, code := range []int{http.StatusNotFound, http.StatusServiceUnavailable} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			srv := newServer(t, code, "xato")
			client := &http.Client{Timeout: time.Second}
			_, err := fetchRate(client, srv.URL, "UZS")
			if err == nil {
				t.Fatalf("status %d: fetchRate xato qaytarmadi", code)
			}
			var se *StatusError
			if !errors.As(err, &se) {
				t.Fatalf("status %d: xato «%v» *StatusError emas; errors.As topishi uchun "+
					"&StatusError{Code: resp.StatusCode} qaytaring (TASK.md, 2-qadam)", code, err)
			}
			if se.Code != code {
				t.Errorf("status %d: se.Code = %d; kutilgan %d (TASK.md, 2-qadam)", code, se.Code, code)
			}
		})
	}
}

// 3-qadam: tashqi serverning juda katta javobi (chuqurroq qismi: io.LimitReader).
func TestVazifaKattaJavob(t *testing.T) {
	big := `{"result":"success","rates":{"UZS":12000.5},"pad":"` + strings.Repeat("x", 2<<20) + `"}`
	srv := newServer(t, http.StatusOK, big)
	client := &http.Client{Timeout: 5 * time.Second}
	rate, err := fetchRate(client, srv.URL, "UZS")
	if err == nil {
		t.Fatalf("javob 2 MB dan katta, fetchRate esa %v va nil qaytardi; kerak: xato — body'ni "+
			"io.LimitReader(resp.Body, 1<<20) orqali oʻqing (TASK.md, 3-qadam)", rate)
	}
}
