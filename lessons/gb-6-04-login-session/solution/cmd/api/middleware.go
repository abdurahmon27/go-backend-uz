package main

// Middleware (GB-6.01). Har biri bitta shaklda: http.Handler oladi (next) va uni oʻragan yangi
// http.Handler qaytaradi. Qaytgan handler oʻz ishini qiladi va next.ServeHTTP'ni chaqiradi:
// chaqiruvdan oldingi kod soʻrov handler'ga yetguncha, keyingi kod javob yozilgandan keyin ishlaydi.
// GB-6.02: middleware'larning oʻz xato javoblari (panic 500, 415, mux 404/405) ham writeProblem'dan chiqadi:
// ularda Go xatosi yoʻq, faqat status (problem.go).

import (
	"crypto/rand"
	"log"
	"mime"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
)

// newHandler — server'ga beriladigan butun zanjir. Tashqaridan ichkariga:
// requestID → logRequests → recoverPanic → muxErrors → mux → handler (requireJSON — faqat POST route'larda, newMux).
// Tartib muhim: recoverPanic logRequests'ning ichida boʻlsin, aks holda panic boʻlgan soʻrov
// log'ga tushmaydi (panic logRequests'ning log qatoridan oʻtib ketadi).
func newHandler(st *store.Store) http.Handler {
	return requestID(logRequests(recoverPanic(muxErrors(newMux(st)))))
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
// stack'ni log'ga yozadi va mijozga 500 problem qaytaradi (GB-6.02; stack mijozga chiqmaydi, RFC 9457 §5).
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
				writeProblem(w, newProblem(http.StatusInternalServerError, ""))
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

// requireJSON (vazifa javobi, TASK.md 1-qadam): POST soʻrovining body'si JSON boʻlishi shart.
// Javobda u newMux'da faqat body oladigan route'larni oʻraydi: zanjir boshida tursa, POST /health
// yoki yoʻq yoʻlga ham 405/404 oʻrniga 415 berardi.
// Content-Type yoʻq yoki boshqa tur boʻlsa — 415 Unsupported Media Type, problem shaklida (GB-6.02).
// mime.ParseMediaType "application/json; charset=utf-8"dan faqat turini ajratadi.
func requireJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil || mediaType != "application/json" {
				writeProblem(w, newProblem(http.StatusUnsupportedMediaType, "Content-Type: application/json boʻlishi kerak"))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// muxErrors (🟡 ixtiyoriy vazifa javobi, TASK.md 2-qadam): mux'ning oʻz 404 va 405 javoblari
// oddiy matn (GB-4.08). mux.Handler(r) pattern'ni soʻrovni bajarmasdan topadi: pattern boʻsh boʻlsa,
// javobni mux'ning oʻzi yozadi — shunda writer'ni jsonErrorWriter bilan oʻraymiz.
func muxErrors(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, pattern := mux.Handler(r); pattern != "" {
			mux.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(&jsonErrorWriter{ResponseWriter: w}, r)
	})
}

// jsonErrorWriter 404 va 405'ni problem shaklida yozadi (GB-6.02) va mux'ning matnini tashlab yuboradi.
// 405'dagi Allow sarlavhasi saqlanadi: mux uni WriteHeader'dan oldin qoʻygan.
type jsonErrorWriter struct {
	http.ResponseWriter
	replaced bool
}

func (j *jsonErrorWriter) WriteHeader(code int) {
	switch code {
	case http.StatusNotFound:
		j.replaced = true
		writeProblem(j.ResponseWriter, newProblem(code, "bunday yoʻl yoʻq"))
	case http.StatusMethodNotAllowed:
		j.replaced = true
		writeProblem(j.ResponseWriter, newProblem(code, "bu yoʻl uchun bunday metod yoʻq"))
	default:
		j.ResponseWriter.WriteHeader(code)
	}
}

func (j *jsonErrorWriter) Write(b []byte) (int, error) {
	if j.replaced {
		return len(b), nil
	}
	return j.ResponseWriter.Write(b)
}

func (j *jsonErrorWriter) Unwrap() http.ResponseWriter {
	return j.ResponseWriter
}
