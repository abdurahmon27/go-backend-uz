package main

// Bir xil API xatolari (GB-6.02): har xato javobi RFC 9457 "problem details" shaklida.
// Qoida: err.Error() — dasturchi uchun matn, u faqat log'ga boradi. Mijoz faqat biz tanlab
// yozgan matnni oladi (detail, errors) va log'dagi qatorni topish uchun request_id'ni.

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/abdurahmon27/bozor/internal/store"
)

// problem — RFC 9457 obyekti. type, title, status, detail — standart aʼzolar;
// errors va request_id — bozor'ning kengaytmalari (RFC 9457 §3.2).
type problem struct {
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Status    int               `json:"status"`
	Detail    string            `json:"detail,omitempty"`
	Errors    map[string]string `json:"errors,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

// newProblem — maxsus turi yoʻq xato: type "about:blank", title esa statusning standart nomi
// (RFC 9457 §4.2.1: about:blank boʻlsa title — HTTP status nomi).
func newProblem(status int, detail string) problem {
	return problem{
		Type:   "about:blank",
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
	}
}

// writeProblem — har xato javobi shu funksiyadan chiqadi: sarlavha → status → body (GB-4.03 tartibi),
// media turi application/problem+json. request_id — GB-6.01'dagi requestID qoʻygan X-Request-Id: mijoz uni
// aytsa, log'dagi qatorni topamiz. (Middleware'siz, testda, sarlavha boʻsh: omitempty uni tashlaydi.)
func writeProblem(w http.ResponseWriter, p problem) {
	p.RequestID = w.Header().Get("X-Request-Id")
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p) // Encode xatosi: log'ga yozish GB-8.02'da (writeJSON kabi)
}

// API qatlamining oʻz sentinel xatolari (GB-2.08): decodeJSON va pathID asl xatoni shular bilan oʻraydi.
var errBadJSON = errors.New("body JSON sifatida oʻqilmadi")

var errBadID = errors.New("id butun son emas")

// fieldErrors — validate natijasi: maydon nomi → nima notoʻgʻri. Error() metodi bor, demak u error:
// handler uni boshqa xatolar kabi writeError'ga beradi, problemFor esa 422 qiladi.
type fieldErrors map[string]string

func (f fieldErrors) Error() string {
	return fmt.Sprintf("validatsiya: %d ta maydon notoʻgʻri", len(f))
}

// productRules — 422'da mijozga koʻrsatiladigan katalog qoidalari. Store xatosining matni emas:
// u ichki qatlam matni, ertaga unga boshqa tafsilot qoʻshilishi mumkin.
const productRules = "nom 1–20 belgi, narx 1 tiyindan 10 mln soʻmgacha, qoldiq 0–1 000 000"

// problemFor xatoni mijozga ketadigan problem'ga aylantiradi: status errors.Is/As bilan tanlanadi.
// Hech bir case err.Error()ni ishlatmaydi: detail faqat biz tanlagan maydonlardan yigʻiladi.
// switch birinchi mos case'ni oladi: aniqroq case (413) umumiyroq case'dan (errBadJSON) oldin turadi.
func problemFor(err error) problem {
	var se *store.StockError
	var te *json.UnmarshalTypeError
	var tooBig *http.MaxBytesError
	var fe fieldErrors
	switch {
	case errors.Is(err, store.ErrNotFound):
		return newProblem(http.StatusNotFound, "bunday ID'li mahsulot yoʻq")
	case errors.As(err, &se):
		return problem{
			Type:   "/problems/insufficient-stock",
			Title:  "Omborda yetarli emas",
			Status: http.StatusConflict,
			Detail: fmt.Sprintf("%s: omborda %d dona, soʻralgan %d", se.Name, se.Stock, se.Want),
		}
	case errors.As(err, &fe): // vazifa, 1-qadam
		return problem{
			Type:   "/problems/validation",
			Title:  "Qiymatlar notoʻgʻri",
			Status: http.StatusUnprocessableEntity,
			Errors: fe,
		}
	case errors.Is(err, store.ErrBadProduct), errors.Is(err, store.ErrBadQuantity): // vazifa, 2-qadam
		return problem{
			Type:   "/problems/validation",
			Title:  "Qiymatlar notoʻgʻri",
			Status: http.StatusUnprocessableEntity,
			Detail: productRules,
		}
	case errors.As(err, &tooBig): // vazifa, 3-qadam: errBadJSON'dan oldin
		return newProblem(http.StatusRequestEntityTooLarge, fmt.Sprintf("body %d baytdan oshmasin", tooBig.Limit))
	case errors.Is(err, errBadID):
		return newProblem(http.StatusBadRequest, "id butun son boʻlishi kerak")
	case errors.As(err, &te) && te.Field != "":
		return newProblem(http.StatusBadRequest, te.Field+": qiymat turi notoʻgʻri")
	case errors.Is(err, errBadJSON):
		return newProblem(http.StatusBadRequest, "body bitta JSON obyekt boʻlsin, faqat maʼlum kalitlar bilan")
	default:
		return newProblem(http.StatusInternalServerError, "")
	}
}

// writeError — handler'dagi har Go xatosi shu yerdan chiqadi: log'ga toʻliq xato (err.Error()),
// mijozga problemFor tanlagan xavfsiz problem. Ikkalasida bir xil request ID.
// Metod, yoʻl va status logRequests qatorida (GB-6.01); bu qator — xatoning oʻzi. %q: xato matnida
// mijoz yuborgan narsa (yoʻl, kalit) boʻlishi mumkin, yangi qator log'da \n boʻlib qoladi (GB-6.01).
func writeError(w http.ResponseWriter, err error) {
	p := problemFor(err)
	log.Printf("%s xato %d: %q", w.Header().Get("X-Request-Id"), p.Status, err.Error())
	writeProblem(w, p)
}
