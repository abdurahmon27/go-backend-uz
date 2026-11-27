package main

// Middleware testlari (GB-6.01): log qatori, panic → 500, request ID va zanjir tartibi.

import (
	"bytes"
	"log"
	"net/http"
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
	if rec.Code != 500 || !strings.Contains(rec.Body.String(), `"code":"internal"`) {
		t.Errorf("panic'dan keyin: %d %s, kutilgan 500 va internal", rec.Code, rec.Body)
	}
	if strings.Contains(rec.Body.String(), "sinov") {
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
func TestLogInjection(t *testing.T) {
	buf := captureLog(t)
	send(newHandler(testStore(t)), "GET", "/products/1%0Asoxta-qator", "")
	if strings.Count(buf.String(), "\n") != 1 {
		t.Errorf("bitta soʻrov — bitta log qatori boʻlishi kerak: %q", buf)
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
