package main

// Middleware (GB-6.01). Har biri bitta shaklda: http.Handler oladi (next) va uni oʻragan yangi
// http.Handler qaytaradi. Qaytgan handler oʻz ishini qiladi va next.ServeHTTP'ni chaqiradi:
// chaqiruvdan oldingi kod soʻrov handler'ga yetguncha, keyingi kod javob yozilgandan keyin ishlaydi.
// GB-6.02: middleware'larning oʻz xato javoblari (panic 500, 415, mux 404/405) ham writeProblem'dan chiqadi:
// ularda Go xatosi yoʻq, faqat status (problem.go).
// GB-8.02: log'lar log/slog bilan, JSON qatorlar. requestID har soʻrovga request_id biriktirilgan logger beradi
// (soʻrov context'ida); soʻrov ichidagi har log qatori uni loggerFrom(r.Context()) bilan oladi.

import (
	"context"
	"crypto/rand"
	"log/slog"
	"mime"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
)

// newHandler — server'ga beriladigan butun zanjir (GB-6.01, GB-6.07). Tashqaridan ichkariga:
// requestID → logRequests → recoverPanic → cors → protectCSRF → muxErrors → mux → route'ning oʻz qavatlari
// (limitByIP, requireUser, requireAdmin, requireJSON — newMux'da).
// requestID eng tashqarida: undan ichkaridagi hamma (logRequests ham) request_id'li logger'ni oladi.
// Tartib muhim: recoverPanic logRequests'ning ichida boʻlsin, aks holda panic boʻlgan soʻrov
// log'ga tushmaydi (panic logRequests'ning log qatoridan oʻtib ketadi). cors protectCSRF'dan tashqarida:
// 403 javobida ham CORS sarlavhalari boʻlsin, ruxsatli frontend xato matnini oʻqiy olsin (GB-6.07).
func newHandler(st *store.Store, logger *slog.Logger, bg *sync.WaitGroup) http.Handler {
	return requestID(logger,
		logRequests(
			recoverPanic(
				cors(
					protectCSRF(
						muxErrors(newMux(st, bg)))))))
}

// requestID har soʻrovga tasodifiy ID beradi va uni javobning X-Request-Id sarlavhasiga yozadi.
// Mijoz xato haqida shu ID bilan xabar beradi, biz log'dan aynan shu soʻrovni topamiz.
// Mijoz yuborgan X-Request-Id'ga ishonmaymiz: ID'ni har doim server oʻzi yasaydi.
// GB-8.02: ID logger'ga ham biriktiriladi (With) va shu logger soʻrov context'iga qoʻyiladi (GB-6.04'dagi
// userKey kabi): soʻrov ichida yozilgan har qatorda "request_id" bor, uni hech kim qoʻshishni unutmaydi.
func requestID(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := rand.Text()
		w.Header().Set("X-Request-Id", id)
		ctx := context.WithValue(r.Context(), loggerKey{}, logger.With("request_id", id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loggerKey — context'dagi logger'ning kaliti. Oʻz turimiz: boshqa paketning kaliti bilan toʻqnashmaydi
// (go doc context.WithValue).
type loggerKey struct{}

// loggerFrom — requestID context'ga qoʻygan logger. Soʻrovdan tashqarida (yoki zanjirdan oʻtmagan testda)
// slog.Default() — request_id'siz, lekin log yoʻqolmaydi.
func loggerFrom(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
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

// logRequests har soʻrovdan keyin bitta qator yozadi: metod, yoʻl, status, vaqt (request_id — logger'da).
// status 200'dan boshlanadi: handler WriteHeader'ni chaqirmay Write qilsa, javob 200 boʻladi (GB-4.03),
// Write esa bizning WriteHeader'imizdan oʻtmaydi.
// GB-6.01'da yoʻl %q bilan yozilardi (log injection). JSON handler satrni oʻzi qochiradi: yoʻldagi yangi qator
// (%0A) "\n" boʻlib, bitta JSON qator ichida qoladi.
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		loggerFrom(r.Context()).Info("soʻrov",
			"method", r.Method, "path", r.URL.Path,
			"status", rec.status, "dur", time.Since(start))
	})
}

// recoverPanic — GB-2.09'dagi safely, har soʻrov uchun: handler'dagi panic'ni toʻxtatadi,
// stack'ni log'ga yozadi va mijozga 500 problem qaytaradi (GB-6.02; stack mijozga chiqmaydi, RFC 9457 §5).
// GB-6.01 videosidan ikki farq (GB-6.01 README, «recoverPanic: videodan tashqari ikki himoya»):
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
				loggerFrom(r.Context()).Error("panic", "panic", v, "stack", string(debug.Stack()))
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

// requireJSON (GB-6.01 vazifasi): POST va PATCH (GB-6.06) soʻrovining body'si JSON boʻlishi shart.
// Javobda u newMux'da faqat body oladigan route'larni oʻraydi: zanjir boshida tursa, POST /health
// yoki yoʻq yoʻlga ham 405/404 oʻrniga 415 berardi.
// Content-Type yoʻq yoki boshqa tur boʻlsa — 415 Unsupported Media Type, problem shaklida (GB-6.02).
// mime.ParseMediaType "application/json; charset=utf-8"dan faqat turini ajratadi.
func requireJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPatch {
			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil || mediaType != "application/json" {
				writeProblem(w, newProblem(http.StatusUnsupportedMediaType, "Content-Type: application/json boʻlishi kerak"))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// muxErrors (GB-6.01'ning 🟡 vazifasi): mux'ning oʻz 404 va 405 javoblari
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
