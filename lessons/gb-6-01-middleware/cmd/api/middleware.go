package main

// Middleware (GB-6.01). Har biri bitta shaklda: http.Handler oladi (next) va uni oʻragan yangi
// http.Handler qaytaradi. Qaytgan handler oʻz ishini qiladi va next.ServeHTTP'ni chaqiradi:
// chaqiruvdan oldingi kod soʻrov handler'ga yetguncha, keyingi kod javob yozilgandan keyin ishlaydi.

import (
	"crypto/rand"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
)

// newHandler — server'ga beriladigan butun zanjir. Tashqaridan ichkariga:
// requestID → logRequests → recoverPanic → mux → handler.
// Tartib muhim: recoverPanic logRequests'ning ichida boʻlsin, aks holda panic boʻlgan soʻrov
// log'ga tushmaydi (panic logRequests'ning log qatoridan oʻtib ketadi).
func newHandler(st *store.Store) http.Handler {
	// VAZIFA (TASK.md, 1-qadam): requireJSON'ni shu zanjirga qoʻshing.
	return requestID(logRequests(recoverPanic(newMux(st))))
}

// requestID har soʻrovga tasodifiy ID beradi va uni javobning X-Request-Id sarlavhasiga yozadi.
// Mijoz xato haqida shu ID bilan xabar beradi, biz log'dan aynan shu soʻrovni topamiz.
// Mijoz yuborgan X-Request-Id'ga ishonmaymiz: ID'ni har doim server oʻzi yasaydi.
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", rand.Text())
		next.ServeHTTP(w, r)
	})
}

// statusRecorder — ResponseWriter'ni oʻrab, handler yozgan status kodni eslab qoladi.
// http.ResponseWriter maydon nomisiz (embedding): Header va Write ichkaridagi writer'niki,
// WriteHeader'ni esa oʻzimiz qayta yozdik.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// Unwrap ichkaridagi asl writer'ni qaytaradi: http.ResponseController (Flush va boshqalar)
// oʻralgan writer orqali ham ishlashi uchun (go doc net/http.NewResponseController).
func (rec *statusRecorder) Unwrap() http.ResponseWriter {
	return rec.ResponseWriter
}

// logRequests har soʻrovdan keyin bitta qator yozadi: ID, metod, yoʻl, status, vaqt.
// status 200'dan boshlanadi: handler WriteHeader'ni chaqirmay Write qilsa, javob 200 boʻladi (GB-4.03),
// Write esa bizning WriteHeader'imizdan oʻtmaydi.
// Yoʻl %q bilan: mijoz yoʻlga yangi qator (%0A) yashirib, log'ga soxta qator qoʻsha olmaydi.
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s %q %d %s", w.Header().Get("X-Request-Id"),
			r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}

// recoverPanic — GB-2.09'dagi safely, har soʻrov uchun: handler'dagi panic'ni toʻxtatadi,
// stack'ni log'ga yozadi va mijozga 500 + bizning JSON xato shaklini qaytaradi.
// Videodagidan ikki farq (README, «recoverPanic: videodan tashqari ikki himoya»):
//   - http.ErrAbortHandler — "javobni ataylab uz" degan panic (go doc net/http.ErrAbortHandler):
//     uni toʻxtatmaymiz, net/http'ning oʻziga qaytaramiz;
//   - handler javobni yozishni boshlagan boʻlsa (status yoki body ketgan), 500'ni yuborib boʻlmaydi:
//     yarim javob oxiriga JSON qoʻshilib qoladi. Shunda panic'ni log'ga yozib, javobni ErrAbortHandler
//     bilan uzamiz: mijoz buzuq 200 emas, uzilgan javob oladi. Bu holda logRequests qatori yozilmaydi,
//     panic qatori esa ID bilan logda bor.
//
// Eslatma: handler oʻzi ishga tushirgan goroutine'dagi panic'ni bu middleware ushlamaydi (u boshqa
// goroutine'ning stack'ida; recover faqat panic qilgan goroutine'da ishlaydi) — bunday panic butun serverni yiqitadi.
func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &startWriter{ResponseWriter: w}
		defer func() {
			if v := recover(); v != nil {
				if v == http.ErrAbortHandler {
					panic(v)
				}
				log.Printf("%s panic: %v\n%s", w.Header().Get("X-Request-Id"), v, debug.Stack())
				if sw.started {
					panic(http.ErrAbortHandler)
				}
				writeError(w, http.StatusInternalServerError, apiError{Code: "internal", Message: "ichki xato"})
			}
		}()
		next.ServeHTTP(sw, r)
	})
}

// startWriter javob boshlanganini (WriteHeader yoki Write chaqirilganini) eslab qoladi: recoverPanic uchun.
type startWriter struct {
	http.ResponseWriter
	started bool
}

func (sw *startWriter) WriteHeader(code int) {
	sw.started = true
	sw.ResponseWriter.WriteHeader(code)
}

func (sw *startWriter) Write(b []byte) (int, error) {
	sw.started = true
	return sw.ResponseWriter.Write(b)
}

func (sw *startWriter) Unwrap() http.ResponseWriter {
	return sw.ResponseWriter
}

// requireJSON — VAZIFA (TASK.md, 1-qadam): POST soʻrovida Content-Type application/json
// boʻlmasa — 415 va bizning xato shaklimiz. Hozircha hech narsa tekshirmaydi.
func requireJSON(next http.Handler) http.Handler {
	return next
}

// muxErrors — 🟡 ixtiyoriy (TASK.md, 2-qadam): mux'ning oʻz 404/405 javoblarini JSON qilish.
// Hozircha mux'ni oʻzgarishsiz qaytaradi.
func muxErrors(mux *http.ServeMux) http.Handler {
	return mux
}
