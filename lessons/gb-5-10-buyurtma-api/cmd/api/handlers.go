package main

// Route'lar va handler'lar (GB-4.08'dagi bilan bir xil yoʻl): soʻrovni oʻqiydi (yoʻl yoki body),
// store metodini r.Context() bilan chaqiradi, natijani writeJSON, xatoni writeStoreError bilan yozadi.

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/abdurahmon27/bozor/internal/store"
)

// api — handler'lar bazani shu struct orqali koʻradi. Global oʻzgaruvchi yoʻq: har test oʻz bazasi
// (oʻz sxemasi) bilan oʻz newMux'ini yasaydi.
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
	mux.HandleFunc("DELETE /products/{id}", a.deleteProduct)
	mux.HandleFunc("POST /orders", a.createOrder)
	// VAZIFA (TASK.md, 2-qadam): GET /orders/{id} → a.getOrder
	return mux
}

// healthHandler — server tirikligini tekshirish: har soʻrovga «ok» (GB-4.03).
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

// listProducts — GET /products: hamma mahsulotlar, ID tartibida.
func (a *api) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := a.st.List(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, products)
}

// getProduct — GET /products/{id}.
func (a *api) getProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	p, err := a.st.Get(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// productInput — POST /products body'si: faqat mijoz berishi mumkin boʻlgan maydonlar (ID'ni baza beradi).
type productInput struct {
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

// validateProduct qiymatlarni tekshiradi (GB-4.06): kalit — maydon nomi, qiymat — nima notoʻgʻri.
// Toʻliq qoidalar (nom uzunligi, UTF-8, yuqori chegaralar) store.Add ichida: u ham xato qaytarsa,
// writeStoreError uni 422 qiladi.
func validateProduct(in productInput) map[string]string {
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

// createProduct — POST /products: decodeJSON → validateProduct (422) → store.Add → 201 va yangi mahsulot.
func (a *api) createProduct(w http.ResponseWriter, r *http.Request) {
	var in productInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if f := validateProduct(in); len(f) > 0 {
		writeError(w, http.StatusUnprocessableEntity, apiError{Code: "validation_failed", Message: "qiymatlar notoʻgʻri", Fields: f})
		return
	}
	p, err := a.st.Add(r.Context(), in.Name, in.Price, in.Stock)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/products/%d", p.ID))
	writeJSON(w, http.StatusCreated, p)
}

// deleteProduct — DELETE /products/{id}: 204 (body yoʻq), 400, 404, buyurtmalarda boʻlsa 409.
func (a *api) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := a.st.Delete(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// itemInput — buyurtmadagi bitta qator, mijoz yuborgani (narxni mijoz bermaydi: uni baza beradi).
type itemInput struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// orderInput — POST /orders body'si: {"customer": "Aziza", "items": [{"product_id": 1, "quantity": 2}]}.
type orderInput struct {
	Customer string      `json:"customer"`
	Items    []itemInput `json:"items"`
}

// Buyurtma chegaralari: bitta buyurtmada 1..maxItems qator, har qatorda 1..100 dona.
const (
	maxItems         = 20
	maxCustomerRunes = 100
)

// validateOrder buyurtma qiymatlarini tekshiradi (GB-4.06): kalit — maydon nomi («items[1].quantity»).
func validateOrder(in orderInput) map[string]string {
	f := map[string]string{}
	if in.Customer == "" || utf8.RuneCountInString(in.Customer) > maxCustomerRunes {
		f["customer"] = fmt.Sprintf("1 dan %d belgigacha boʻlsin", maxCustomerRunes)
	}
	if len(in.Items) == 0 || len(in.Items) > maxItems {
		f["items"] = fmt.Sprintf("1 dan %d qatorgacha boʻlsin", maxItems)
	}
	seen := map[int64]bool{}
	for i, it := range in.Items {
		if it.ProductID <= 0 {
			f[fmt.Sprintf("items[%d].product_id", i)] = "0 dan katta boʻlsin"
		}
		if seen[it.ProductID] {
			f[fmt.Sprintf("items[%d].product_id", i)] = "bir mahsulot bir qatorda boʻlsin"
		}
		seen[it.ProductID] = true
		if it.Quantity < 1 || it.Quantity > 100 {
			f[fmt.Sprintf("items[%d].quantity", i)] = "1 dan 100 gacha boʻlsin"
		}
	}
	return f
}

// createOrder — POST /orders: decodeJSON → validateOrder (422) → store.PlaceOrder (404, 409) → 201 va buyurtma.
func (a *api) createOrder(w http.ResponseWriter, r *http.Request) {
	var in orderInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if f := validateOrder(in); len(f) > 0 {
		writeError(w, http.StatusUnprocessableEntity, apiError{Code: "validation_failed", Message: "qiymatlar notoʻgʻri", Fields: f})
		return
	}
	items := make([]store.OrderItem, len(in.Items))
	for i, it := range in.Items {
		items[i] = store.OrderItem{ProductID: it.ProductID, Quantity: it.Quantity}
	}
	o, err := a.st.PlaceOrder(r.Context(), in.Customer, items)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/orders/%d", o.ID))
	writeJSON(w, http.StatusCreated, o)
}

// pathID yoʻldagi {id}'ni oʻqiydi (GB-4.05). Son boʻlmasa 400 yozadi va false qaytaradi.
func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, apiError{Code: "bad_id", Message: "id butun son boʻlishi kerak"})
		return 0, false
	}
	return id, true
}

// writeStoreError store xatosini HTTP javobiga aylantiradi: bitta joyda, hamma handler uchun (GB-4.08).
// Store xabarlari mijozga koʻrsatsa boʻladigan matn; kutilmagan xatoning (masalan baza uzilgan) matni
// ichki tafsilot: u server logiga yoziladi, mijozga faqat 500 va «ichki xato» chiqadi.
func writeStoreError(w http.ResponseWriter, err error) {
	var se *store.StockError
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, apiError{Code: "not_found", Message: err.Error()})
	// VAZIFA (TASK.md, 3-qadam): store.ErrOrderNotFound ham 404 not_found boʻlsin
	case errors.As(err, &se):
		writeError(w, http.StatusConflict, apiError{Code: "insufficient_stock", Message: err.Error()})
	case errors.Is(err, store.ErrDuplicate), errors.Is(err, store.ErrInUse):
		writeError(w, http.StatusConflict, apiError{Code: "conflict", Message: err.Error()})
	case errors.Is(err, store.ErrBadQuantity), errors.Is(err, store.ErrBadProduct):
		writeError(w, http.StatusUnprocessableEntity, apiError{Code: "validation_failed", Message: err.Error()})
	default:
		log.Printf("ichki xato: %v", err) // tafsilot server logiga (slog: GB-8.02), mijozga emas
		writeError(w, http.StatusInternalServerError, apiError{Code: "internal", Message: "ichki xato"})
	}
}
