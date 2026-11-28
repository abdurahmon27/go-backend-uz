package main

// Middleware testlari (GB-6.01): log qatori, panic → 500, request ID va zanjir tartibi.

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// captureLog standart log'ni test davomida buferga yoʻnaltiradi va oxirida qaytaradi.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	return &buf
}

// boom — handler'dagi bug oʻrnida: har soʻrovda panic qiladi.
var boom = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	panic("sinov: 13-mahsulotda bug")
})

func TestLogRequests(t *testing.T) {
	tests := []struct {
		path string
		want string // log qatorida boʻlishi kerak
	}{
		{"/products/1", `GET "/products/1" 200`},
		{"/products/99", `GET "/products/99" 404`},
		{"/health", `GET "/health" 200`}, // WriteHeader yoʻq: status 200'dan boshlanadi
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			buf := captureLog(t)
			send(newHandler(testStore(t)), "GET", tc.path, "")
			if !strings.Contains(buf.String(), tc.want) {
				t.Errorf("log: %q, kutilgan boʻlak %q", buf, tc.want)
			}
		})
	}
}

func TestRecoverPanic(t *testing.T) {
	buf := captureLog(t)
	rec := send(recoverPanic(boom), "GET", "/products/13", "")
	if rec.Code != 500 || !strings.Contains(rec.Body.String(), `"title":"Internal Server Error"`) {
		t.Errorf("panic'dan keyin: %d %s, kutilgan 500 problem", rec.Code, rec.Body)
	}
	if strings.Contains(rec.Body.String(), "sinov") || strings.Contains(rec.Body.String(), "goroutine") {
		t.Errorf("panic matni mijozga chiqdi: %s", rec.Body)
	}
	if !strings.Contains(buf.String(), "panic: sinov: 13-mahsulotda bug") || !strings.Contains(buf.String(), "goroutine") {
		t.Errorf("log'da panic va stack boʻlishi kerak: %q", buf)
	}
}

// Toʻgʻri tartib: recoverPanic logRequests ichida — panic boʻlgan soʻrov ham 500 bilan log'da.
// Notoʻgʻri tartib: recoverPanic tashqarida — log qatori umuman yozilmaydi.
func TestOrder(t *testing.T) {
	buf := captureLog(t)
	rec := send(requestID(logRequests(recoverPanic(boom))), "GET", "/products/13", "")
	id := rec.Header().Get("X-Request-Id")
	if !strings.Contains(buf.String(), id+` GET "/products/13" 500`) {
		t.Errorf("toʻgʻri tartib: log'da %q bilan 500 qatori yoʻq: %q", id, buf)
	}

	buf.Reset()
	send(requestID(recoverPanic(logRequests(boom))), "GET", "/products/13", "")
	if strings.Contains(buf.String(), `GET "/products/13"`) {
		t.Errorf("notoʻgʻri tartibda log qatori boʻlmasligi kerak edi: %q", buf)
	}
}

func TestRequestID(t *testing.T) {
	captureLog(t)
	h := newHandler(testStore(t))
	a := send(h, "GET", "/health", "").Header().Get("X-Request-Id")
	b := send(h, "GET", "/health", "").Header().Get("X-Request-Id")
	// rand.Text hozir 26 belgi beradi; hujjat faqat «kamida 128 bit» deydi, keyin uzunroq boʻlishi mumkin.
	if len(a) < 26 || a == b {
		t.Errorf("X-Request-Id: %q va %q — kamida 26 belgili va har soʻrovda boshqa boʻlishi kerak", a, b)
	}
}

// http.ErrAbortHandler toʻxtatilmaydi: net/http'ning oʻziga qaytadi.
func TestAbortHandler(t *testing.T) {
	captureLog(t)
	abort := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { panic(http.ErrAbortHandler) })
	defer func() {
		if v := recover(); v != http.ErrAbortHandler {
			t.Errorf("recover: %v, kutilgan http.ErrAbortHandler qayta panic qilinishi", v)
		}
	}()
	send(recoverPanic(abort), "GET", "/", "")
}

// Yoʻldagi yangi qator (%0A) log'da yangi qator boʻlmaydi: %q uni \n deb yozadi.
// GB-6.02: xatoli soʻrov ikki qator yozadi (writeError'ning xato qatori + logRequests qatori),
// ikkalasi ham %q bilan — soxta qator boshlanmaydi.
func TestLogInjection(t *testing.T) {
	buf := captureLog(t)
	send(newHandler(testStore(t)), "GET", "/products/1%0Asoxta-qator", "")
	if strings.Count(buf.String(), "\n") != 2 || strings.Contains(buf.String(), "\nsoxta-qator") {
		t.Errorf("xatoli soʻrov — ikki log qatori, soxta qatorsiz boʻlishi kerak: %q", buf)
	}
}

// Javob boshlangandan keyin panic: 500 yuborib boʻlmaydi — javob ErrAbortHandler bilan uziladi.
func TestPanicAfterWrite(t *testing.T) {
	buf := captureLog(t)
	half := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":1}`))
		panic("sinov: yozib boʻlgach")
	})
	defer func() {
		if v := recover(); v != http.ErrAbortHandler {
			t.Errorf("recover: %v, kutilgan http.ErrAbortHandler (javob uzilishi kerak)", v)
		}
		if !strings.Contains(buf.String(), "panic: sinov: yozib boʻlgach") {
			t.Errorf("panic log'da boʻlishi kerak: %q", buf)
		}
	}()
	send(recoverPanic(half), "GET", "/products", "")
}

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

// GB-6.01 vazifasi: requireJSON oʻzi (zanjirsiz).
func TestRequireJSON(t *testing.T) {
	// ok — zanjirning keyingi qismi oʻrnida: unga yetib kelgan soʻrov 204 oladi.
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	h := requireJSON(ok)
	order := `{"product_id":1,"quantity":1}`

	for _, ct := range []string{"", "text/plain", "application/x-www-form-urlencoded"} {
		rec := post(h, "/orders", ct, order)
		if rec.Code != http.StatusUnsupportedMediaType {
			t.Fatalf("requireJSON: POST, Content-Type %q → status %d, kutilgan 415", ct, rec.Code)
		}
		if rec.Header().Get("Content-Type") != "application/problem+json" || !strings.Contains(rec.Body.String(), `"status":415`) {
			t.Errorf("requireJSON: 415 javobi problem shaklida boʻlsin; hozir %s", rec.Body)
		}
	}
	for _, ct := range []string{"application/json", "application/json; charset=utf-8"} {
		if rec := post(h, "/orders", ct, order); rec.Code != http.StatusNoContent {
			t.Errorf("requireJSON: Content-Type %q — JSON, soʻrov next'ga oʻtishi kerak edi, hozir status %d %s", ct, rec.Code, rec.Body)
		}
	}
	if rec := send(h, "GET", "/products", ""); rec.Code != http.StatusNoContent {
		t.Errorf("requireJSON: GET soʻrovi body'siz, tekshirilmasligi kerak — status %d", rec.Code)
	}
}

// Zanjir: requireJSON (newMux'dagi POST route'larda) requestID va logRequests ichkarisida.
func TestRequireJSONChain(t *testing.T) {
	buf := captureLog(t)
	h := newHandler(testStore(t))
	rec := post(h, "/products", "text/plain", `{"name":"anor","price":1800000,"stock":7}`)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("newHandler: POST /products, Content-Type text/plain → %d, kutilgan 415", rec.Code)
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

// mux'ning 404 va 405 javoblari ham problem shaklida, Allow saqlanadi.
func TestMuxErrors(t *testing.T) {
	captureLog(t)
	h := newHandler(testStore(t))
	rec := send(h, "GET", "/narsalar", "")
	if rec.Code != 404 || !strings.Contains(rec.Body.String(), `"title":"Not Found"`) || strings.Contains(rec.Body.String(), "page not found") {
		t.Errorf("GET /narsalar: %d %q, kutilgan 404 va faqat problem", rec.Code, rec.Body)
	}
	rec = send(h, "PUT", "/products/1", "")
	if rec.Code != 405 || !strings.Contains(rec.Body.String(), `"title":"Method Not Allowed"`) {
		t.Errorf("PUT /products/1: %d %q, kutilgan 405 problem", rec.Code, rec.Body)
	}
	if rec.Header().Get("Allow") == "" {
		t.Errorf("405 javobida Allow sarlavhasi saqlanishi kerak (RFC 9110)")
	}
	if rec := send(h, "GET", "/products/1", ""); rec.Code != 200 {
		t.Errorf("muxErrors oddiy route'ni buzdi: GET /products/1 → %d", rec.Code)
	}
}
