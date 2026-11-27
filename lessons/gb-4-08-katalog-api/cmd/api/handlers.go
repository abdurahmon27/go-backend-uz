package main

// Route'lar va handler'lar. Har handler bir xil yoʻldan yuradi: soʻrovni oʻqiydi (yoʻl yoki body),
// store metodini chaqiradi, natijani writeJSON bilan, xatoni writeStoreError bilan yozadi.

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/abdurahmon27/bozor/internal/store"
)

// api — handler'lar umumiy katalogni shu struct orqali koʻradi. Global oʻzgaruvchi yoʻq:
// har test oʻz katalogi bilan oʻz newMux'ini yasaydi.
type api struct {
	st *store.Store
}

// newMux — hamma route'lar shu yerda (GB-4.05): main ham, testlar ham bir xil mux'ni ishlatadi.
func newMux(st *store.Store) *http.ServeMux {
	a := &api{st: st}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /products", a.listProducts)
	mux.HandleFunc("GET /products/{id}", a.getProduct)
	mux.HandleFunc("POST /products", a.createProduct)
	mux.HandleFunc("POST /orders", a.createOrder)
	// VAZIFA (TASK.md, 2-qadam): DELETE /products/{id} → a.deleteProduct
	return mux
}

// healthHandler — server tirikligini tekshirish: har soʻrovga «ok» (GB-4.03).
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

// listProducts — GET /products: butun katalog, ID tartibida.
func (a *api) listProducts(w http.ResponseWriter, r *http.Request) {
	products := slices.Collect(a.st.All())
	if products == nil {
		products = []store.Product{} // boʻsh katalog JSON'da null emas, [] boʻlsin (GB-3.06)
	}
	writeJSON(w, http.StatusOK, products)
}

// getProduct — GET /products/{id}.
func (a *api) getProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	p, err := a.st.Get(id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// productInput — POST /products body'si: faqat mijoz berishi mumkin boʻlgan maydonlar (ID ni server beradi).
type productInput struct {
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

// validate qiymatlarni tekshiradi (GB-4.06): kalit — maydon nomi, qiymat — nima notoʻgʻri.
// Boʻsh jadval — hammasi joyida. Toʻliq qoidalar (nom uzunligi, UTF-8) store.Add ichida: u ham
// xato qaytarsa, writeStoreError uni 422 qiladi.
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

// createProduct — POST /products: decodeJSON → validate (422) → store.Add → 201 va yangi mahsulot.
func (a *api) createProduct(w http.ResponseWriter, r *http.Request) {
	var in productInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if f := validate(in); len(f) > 0 {
		writeError(w, http.StatusUnprocessableEntity, apiError{Code: "validation_failed", Message: "qiymatlar notoʻgʻri", Fields: f})
		return
	}
	p, err := a.st.Add(in.Name, in.Price, in.Stock)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/products/%d", p.ID)) // yangi resursning manzili (RFC 9110 §15.3.2)
	writeJSON(w, http.StatusCreated, p)
}

// orderResponse — POST /orders javobi: bitta mahsulot, necha dona, jami narx tiyinda (GB-4.06'dagi Order).
// Nomi Order emas: bozor'ning domen turi Order (buyurtma, GB-WRITER §6) 5-modulda store'ga keladi.
type orderResponse struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
	Total     int64 `json:"total"` // tiyin
}

// orderInput — POST /orders body'si.
type orderInput struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// validateOrder: product_id 0 dan katta, quantity 1 dan 100 gacha (GB-4.06).
func validateOrder(in orderInput) map[string]string {
	f := map[string]string{}
	if in.ProductID <= 0 {
		f["product_id"] = "0 dan katta boʻlsin"
	}
	if in.Quantity < 1 || in.Quantity > 100 {
		f["quantity"] = "1 dan 100 gacha boʻlsin"
	}
	return f
}

// createOrder — POST /orders: decodeJSON → validateOrder (422) → store.Sell (404, 409) → 201 va orderResponse.
func (a *api) createOrder(w http.ResponseWriter, r *http.Request) {
	var in orderInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if f := validateOrder(in); len(f) > 0 {
		writeError(w, http.StatusUnprocessableEntity, apiError{Code: "validation_failed", Message: "qiymatlar notoʻgʻri", Fields: f})
		return
	}
	p, err := a.st.Sell(in.ProductID, in.Quantity)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	// Narx ≤ maxPrice (store), dona ≤ 100 (validateOrder): koʻpaytma int64'ga sigʻadi.
	writeJSON(w, http.StatusCreated, orderResponse{ProductID: p.ID, Quantity: in.Quantity, Total: p.Price * int64(in.Quantity)})
}

// VAZIFA (TASK.md, 2-qadam): deleteProduct — DELETE /products/{id}: 204 No Content, 400, 404.

// pathID yoʻldagi {id}'ni oʻqiydi (GB-4.05). Son boʻlmasa 400 yozadi va false qaytaradi:
// chaqiruvchi darhol return qiladi (decodeJSON kabi).
func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, apiError{Code: "bad_id", Message: "id butun son boʻlishi kerak"})
		return 0, false
	}
	return id, true
}

// writeStoreError store xatosini HTTP javobiga aylantiradi: bitta joyda, hamma handler uchun
// (GB-3.06'dagi CLI message funksiyasining HTTP koʻrinishi). Store xabarlari mijozga koʻrsatsa
// boʻladigan matn; kutilmagan xatoning matni esa ichki tafsilot, mijozga chiqmaydi: 500.
func writeStoreError(w http.ResponseWriter, err error) {
	var se *store.StockError
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, apiError{Code: "not_found", Message: err.Error()})
	case errors.As(err, &se):
		writeError(w, http.StatusConflict, apiError{Code: "insufficient_stock", Message: err.Error()})
	case errors.Is(err, store.ErrBadQuantity), errors.Is(err, store.ErrBadProduct):
		writeError(w, http.StatusUnprocessableEntity, apiError{Code: "validation_failed", Message: err.Error()})
	default:
		log.Printf("ichki xato: %v", err) // tafsilot server logiga (slog: GB-8.02), mijozga emas
		writeError(w, http.StatusInternalServerError, apiError{Code: "internal", Message: "ichki xato"})
	}
}
