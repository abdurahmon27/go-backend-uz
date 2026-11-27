package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Product — katalogdagi bitta mahsulot (GB-3.02).
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

var catalog = []Product{
	{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
	{ID: 2, Name: "non", Price: 400_000, Stock: 120},
	{ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
}

// productInput — POST /products body'si: faqat mijoz berishi mumkin boʻlgan maydonlar (ID ni server beradi).
type productInput struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}

// apiError — hamma xato javobining ichki shakli: {"error": {"code": …, "message": …, "fields": …}}.
type apiError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// writeJSON: sarlavha → status → body (GB-4.03 tartibi).
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) // Encode xatosi: log'ga yozish GB-8.02'da
}

// writeError — barcha xatolar shu funksiyadan chiqadi.
func writeError(w http.ResponseWriter, status int, e apiError) {
	writeJSON(w, status, map[string]apiError{"error": e})
}

// decodeJSON body'ni dst'ga oʻqiydi: eng koʻpi 1024 bayt, nomaʼlum kalit yoʻq, bitta JSON qiymat.
// Xato boʻlsa javobni oʻzi yozadi va false qaytaradi: chaqiruvchi darhol return qiladi.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeDecodeError(w, err)
		return false
	}
	// Decode faqat birinchi JSON qiymatni oʻqiydi: undan keyin hech narsa boʻlmasligi kerak (limit ham shunga tegadi).
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			err = errors.New("body'da bitta JSON boʻlishi kerak")
		}
		writeDecodeError(w, err)
		return false
	}
	return true
}

// writeDecodeError: juda katta body — 413, qolgan hamma oʻqish xatosi — 400.
func writeDecodeError(w http.ResponseWriter, err error) {
	var tooBig *http.MaxBytesError
	if errors.As(err, &tooBig) {
		writeError(w, http.StatusRequestEntityTooLarge,
			apiError{Code: "too_large", Message: fmt.Sprintf("body %d baytdan oshmasin", tooBig.Limit)})
		return
	}
	writeError(w, http.StatusBadRequest, apiError{Code: "bad_json", Message: err.Error()})
}

// validate qiymatlarni tekshiradi: kalit — maydon nomi, qiymat — nima notoʻgʻri. Boʻsh jadval — hammasi joyida.
func validate(in productInput) map[string]string {
	f := map[string]string{}
	if in.Name == "" {
		f["name"] = "boʻsh boʻlmasin"
	}
	if in.Price <= 0 {
		f["price"] = "0 dan katta boʻlsin"
	}
	if in.Stock < 0 {
		f["stock"] = "manfiy boʻlmasin"
	}
	return f
}

// addProduct katalogga qoʻshadi: ID = eng katta ID + 1 (GB-3.06'dagi Add'ning oddiy nusxasi).
// Bir vaqtda ikki soʻrov kelsa bu xavfsiz emas: GB-4.07'da -race bilan koʻramiz.
func addProduct(in productInput) Product {
	var last int64
	for _, p := range catalog {
		last = max(last, p.ID)
	}
	p := Product{ID: last + 1, Name: in.Name, Price: in.Price, Stock: in.Stock}
	catalog = append(catalog, p)
	return p
}

// createProduct — POST /products.
func createProduct(w http.ResponseWriter, r *http.Request) {
	var in productInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if f := validate(in); len(f) > 0 {
		writeError(w, http.StatusUnprocessableEntity, apiError{Code: "validation_failed", Message: "qiymatlar notoʻgʻri", Fields: f})
		return
	}
	writeJSON(w, http.StatusCreated, addProduct(in))
}

// Order — buyurtma (VAZIFA, TASK.md): bitta mahsulot, necha dona, jami narx tiyinda.
type Order struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
	Total     int64 `json:"total"` // tiyin
}

// orderInput — POST /orders body'si.
type orderInput struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// validateOrder — VAZIFA: hozir hech narsa tekshirmaydi (har doim "hammasi joyida").
//
// Kerak: product_id 0 dan katta boʻlsin, quantity 1 dan 100 gacha. Notoʻgʻri maydon uchun jadvalga
// kalit (maydon nomi, videodagi kabi json nomi) va xabar qoʻying.
func validateOrder(in orderInput) map[string]string {
	return map[string]string{}
}

// createOrder — VAZIFA: hozir hech narsa yozmaydi.
//
// Kerak: decodeJSON → validateOrder (422, validation_failed) → mahsulot katalogda yoʻq boʻlsa 404 (not_found) →
// 201 va Order (Total = narx × dona). Hamma xato writeError orqali.
func createOrder(w http.ResponseWriter, r *http.Request) {
	// VAZIFA: shu yerga yozing
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /products", createProduct)
	mux.HandleFunc("POST /orders", createOrder)
	return mux
}

// send — serversiz: soʻrovni mux'ga beradi va javobni "status body" qilib chiqaradi (GB-4.03).
func send(h http.Handler, method, path, body string) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	fmt.Printf("%d %s", rec.Code, rec.Body.String())
}

func main() {
	mux := newMux()
	send(mux, http.MethodPost, "/products", `{"name":"non","price":400000}{"x":1}`)
	send(mux, http.MethodPost, "/products", `{"name":"non","price":400000}`+strings.Repeat(" ", 2000))
}
