package main

// Bir xil API xatolari (GB-6.02): har xato javobi RFC 9457 "problem details" shaklida.
// err.Error() — dasturchi uchun matn, u faqat log'ga boradi; mijoz faqat biz yozgan matnni oladi.
// GB-6.03: 422 (fieldErrors) va 409 (email band). GB-6.04: login va session xatolari — 401 Unauthorized.
// GB-6.06: rol yetmasa — 403 Forbidden; begona buyurtma — 404, xuddi yoʻq buyurtma kabi.
// GB-6.07: boshqa origin'dan kelgan POST — 403; login urinishlari koʻp — 429 (ratelimit.go, writeProblem'ni oʻzi chaqiradi).

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/abdurahmon27/bozor/internal/store"
)

// problem — RFC 9457 obyekti. errors va request_id — bozor'ning kengaytmalari (RFC 9457 §3.2).
type problem struct {
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Status    int               `json:"status"`
	Detail    string            `json:"detail,omitempty"`
	Errors    map[string]string `json:"errors,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

// newProblem — maxsus turi yoʻq xato: type "about:blank", title — statusning standart nomi (RFC 9457 §4.2.1).
func newProblem(status int, detail string) problem {
	return problem{Type: "about:blank", Title: http.StatusText(status), Status: status, Detail: detail}
}

// writeProblem — har xato javobi shu funksiyadan chiqadi, media turi application/problem+json.
// 401'da RFC 9110 §11.6.1 WWW-Authenticate sarlavhasini talab qiladi: "Cookie" ro'yxatdan o'tgan sxema emas
// (cookie session uchun standart sxema yo'q), lekin sintaksisi to'g'ri va mijozga qanday kirishni aytadi.
func writeProblem(w http.ResponseWriter, p problem) {
	p.RequestID = w.Header().Get("X-Request-Id")
	if p.Status == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", `Cookie realm="bozor"`)
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p)
}

// errBadJSON — decodeJSON asl xatoni shu bilan oʻraydi.
var errBadJSON = errors.New("body JSON sifatida oʻqilmadi")

// errBadID — pathID yoʻldagi {id}'ni son qila olmadi (GB-6.02).
var errBadID = errors.New("id butun son emas")

// fieldErrors — validate natijasi (GB-6.02): maydon nomi → nima notoʻgʻri; problemFor uni 422 qiladi.
type fieldErrors map[string]string

func (f fieldErrors) Error() string {
	return fmt.Sprintf("validatsiya: %d ta maydon notoʻgʻri", len(f))
}

// problemFor xatoni mijozga ketadigan problem'ga aylantiradi (GB-6.02): status errors.Is/As bilan tanlanadi,
// detail faqat biz yozgan matn. 401'ning ikki holati:
//   - login: email yoki parol notoʻgʻri — qaysi biri ekanini aytmaymiz (OWASP Authentication Cheat Sheet);
//   - session: cookie yoʻq, token notanish yoki muddati oʻtgan — mijoz uchun bir xil: qaytadan login.
func problemFor(err error) problem {
	var te *json.UnmarshalTypeError
	var tooBig *http.MaxBytesError
	var fe fieldErrors
	var se *store.StockError
	switch {
	case errors.Is(err, errBadLogin):
		return newProblem(http.StatusUnauthorized, "email yoki parol notoʻgʻri")
	case errors.Is(err, http.ErrNoCookie), errors.Is(err, store.ErrNoSession):
		return newProblem(http.StatusUnauthorized, "avval login qiling")
	case errors.Is(err, errForbidden): // GB-6.06: kimligi maʼlum, roli yetmaydi
		return newProblem(http.StatusForbidden, "bu amal uchun ruxsat yoʻq")
	case errors.Is(err, errCrossOrigin): // GB-6.07: CSRF tekshiruvi
		return newProblem(http.StatusForbidden, "boshqa saytdan kelgan soʻrov")
	case errors.Is(err, store.ErrOrderNotFound): // GB-6.06: yoʻq ham, begona ham — bir xil javob
		return newProblem(http.StatusNotFound, "bunday buyurtma yoʻq")
	case errors.Is(err, store.ErrNotFound):
		return newProblem(http.StatusNotFound, "bunday ID'li mahsulot yoʻq")
	case errors.As(err, &se):
		return problem{
			Type:   "/problems/insufficient-stock",
			Title:  "Omborda yetarli emas",
			Status: http.StatusConflict,
			Detail: fmt.Sprintf("%s: omborda %d dona, soʻralgan %d", se.Name, se.Stock, se.Want),
		}
	case errors.Is(err, store.ErrDuplicate):
		return newProblem(http.StatusConflict, "bu nomli mahsulot bor")
	case errors.Is(err, store.ErrInUse):
		return newProblem(http.StatusConflict, "mahsulot buyurtmalarda bor: oʻchirib boʻlmaydi")
	case errors.Is(err, store.ErrBadProduct):
		return problem{Type: "/problems/validation", Title: "Qiymatlar notoʻgʻri", Status: http.StatusUnprocessableEntity,
			Detail: "nom 1–20 belgi, narx 1 tiyindan 10 mln soʻmgacha, qoldiq 0–1 000 000"}
	case errors.Is(err, store.ErrBadQuantity):
		return problem{Type: "/problems/validation", Title: "Qiymatlar notoʻgʻri", Status: http.StatusUnprocessableEntity,
			Detail: "savatda 1–20 qator, har qatorda 1–100 dona"}
	case errors.Is(err, errBadID):
		return newProblem(http.StatusBadRequest, "id butun son boʻlishi kerak")
	case errors.As(err, &fe):
		return problem{
			Type:   "/problems/validation",
			Title:  "Qiymatlar notoʻgʻri",
			Status: http.StatusUnprocessableEntity,
			Errors: fe,
		}
	case errors.Is(err, store.ErrEmailTaken): // GB-6.03
		return newProblem(http.StatusConflict, "bu email bilan foydalanuvchi bor")
	case errors.As(err, &tooBig):
		return newProblem(http.StatusRequestEntityTooLarge, fmt.Sprintf("body %d baytdan oshmasin", tooBig.Limit))
	case errors.As(err, &te) && te.Field != "":
		return newProblem(http.StatusBadRequest, te.Field+": qiymat turi notoʻgʻri")
	case errors.Is(err, errBadJSON):
		return newProblem(http.StatusBadRequest, "body bitta JSON obyekt boʻlsin, faqat maʼlum kalitlar bilan")
	default:
		return newProblem(http.StatusInternalServerError, "")
	}
}

// writeError — handler'dagi har Go xatosi shu yerdan chiqadi: log'ga toʻliq xato, mijozga xavfsiz problem.
func writeError(w http.ResponseWriter, err error) {
	p := problemFor(err)
	log.Printf("%s xato %d: %q", w.Header().Get("X-Request-Id"), p.Status, err.Error())
	writeProblem(w, p)
}
