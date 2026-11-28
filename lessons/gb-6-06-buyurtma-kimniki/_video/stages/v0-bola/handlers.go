// _video bosqichi v0-bola: XATO versiya. Faqat videodagi namoyish uchun.
package main

// Route'lar (GB-4.05). Ochiq: /health, GET /products, GET /products/{id}, POST /users (GB-6.03), POST /login, POST /logout.
// Qolganlari requireUser bilan oʻralgan (auth.go): session tirik boʻlmasa, handler'gacha yetib bormaydi — 401.
// Body'li route'lar requireJSON bilan (GB-6.01, 415), u requireUser'ning ichida: avval kimligi, keyin body.
// Katalogni oʻzgartirish (POST, DELETE /products) va /admin/… — ikki qavat: requireUser (kim?) ichida
// requireAdmin (admin'mi?), authz.go.

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/abdurahmon27/bozor/internal/store"
)

// api — handler'lar bazani shu struct orqali koʻradi (GB-5.10).
type api struct {
	st *store.Store
}

// newMux — hamma route'lar shu yerda: main ham, testlar ham bir xil mux'ni ishlatadi.
func newMux(st *store.Store) *http.ServeMux {
	a := &api{st: st}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /products", a.listProducts)
	mux.HandleFunc("GET /products/{id}", a.getProduct)
	mux.Handle("POST /products", a.requireUser(requireAdmin(requireJSON(http.HandlerFunc(a.createProduct)))))
	mux.Handle("DELETE /products/{id}", a.requireUser(requireAdmin(http.HandlerFunc(a.deleteProduct))))
	mux.Handle("POST /users", requireJSON(http.HandlerFunc(a.register))) // GB-6.03
	mux.Handle("POST /login", requireJSON(http.HandlerFunc(a.login)))
	mux.Handle("GET /me", a.requireUser(http.HandlerFunc(a.me)))
	mux.HandleFunc("POST /logout", a.logout) // requireUser'siz: eskirgan cookie ham tozalansin (GB-6.04)
	mux.Handle("POST /logout-all", a.requireUser(http.HandlerFunc(a.logoutAll)))
	mux.Handle("POST /orders", a.requireUser(requireJSON(http.HandlerFunc(a.createOrder))))
	mux.Handle("GET /orders/{id}", a.requireUser(http.HandlerFunc(a.getOrder)))
	mux.Handle("PATCH /orders/{id}", a.requireUser(requireJSON(http.HandlerFunc(a.updateNote))))
	mux.Handle("GET /admin/users", a.requireUser(requireAdmin(http.HandlerFunc(a.listUsers))))
	return mux
}

// healthHandler — server tirikligini tekshirish: har soʻrovga «ok» (GB-4.03).
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

// me — GET /me: session egasi. requireUser uni context'ga qoʻygan.
func (a *api) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, currentUser(r))
}

// listProducts — GET /products: katalog hammaga ochiq.
func (a *api) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := a.st.List(r.Context())
	if err != nil {
		writeError(w, err)
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
		writeError(w, err)
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
// problemFor uni 422 qiladi.
func validateProduct(in productInput) fieldErrors {
	f := fieldErrors{}
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

// createProduct — POST /products (faqat admin, GB-6.06): decodeJSON → validateProduct (422) → store.Add → 201 va yangi mahsulot.
func (a *api) createProduct(w http.ResponseWriter, r *http.Request) {
	var in productInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if f := validateProduct(in); len(f) > 0 {
		writeError(w, f)
		return
	}
	p, err := a.st.Add(r.Context(), in.Name, in.Price, in.Stock)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/products/%d", p.ID))
	writeJSON(w, http.StatusCreated, p)
}

// deleteProduct — DELETE /products/{id} (faqat admin, GB-6.06): 204 (body yoʻq), 400, 404, buyurtmalarda boʻlsa 409.
func (a *api) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := a.st.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// itemInput — buyurtmadagi bitta qator, mijoz yuborgani (narxni mijoz bermaydi: uni baza beradi).
type itemInput struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// orderInput — POST /orders body'si. Egasi (user_id) bu yerda YOʻQ: uni session beradi. Mijoz body'ga
// "user_id" qoʻshsa, decodeJSON nomaʼlum kalit deb 400 qaytaradi (DisallowUnknownFields, GB-4.06).
type orderInput struct {
	Items []itemInput `json:"items"`
	Note  string      `json:"note"`
}

// Buyurtma chegaralari: bitta buyurtmada 1..maxItems qator, har qatorda 1..100 dona, izoh 200 belgigacha.
const (
	maxItems     = 20
	maxNoteRunes = 200
)

// validateNote izohni tekshiradi: UTF-8, eng koʻpi maxNoteRunes belgi.
func validateNote(note string, f fieldErrors) {
	if !utf8.ValidString(note) || utf8.RuneCountInString(note) > maxNoteRunes {
		f["note"] = fmt.Sprintf("%d belgigacha boʻlsin", maxNoteRunes)
	}
}

// validateOrder buyurtma qiymatlarini tekshiradi (GB-4.06): kalit — maydon nomi («items[1].quantity»).
func validateOrder(in orderInput) fieldErrors {
	f := fieldErrors{}
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
	validateNote(in.Note, f)
	return f
}

// createOrder — POST /orders: buyurtma session egasi nomidan yoziladi (currentUser(r).ID), 201.
func (a *api) createOrder(w http.ResponseWriter, r *http.Request) {
	var in orderInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if f := validateOrder(in); len(f) > 0 {
		writeError(w, f)
		return
	}
	items := make([]store.OrderItem, len(in.Items))
	for i, it := range in.Items {
		items[i] = store.OrderItem{ProductID: it.ProductID, Quantity: it.Quantity}
	}
	o, err := a.st.PlaceOrder(r.Context(), currentUser(r).ID, in.Note, items)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/orders/%d", o.ID))
	writeJSON(w, http.StatusCreated, o)
}

// getOrder — GET /orders/{id}: faqat oʻz buyurtmangiz. Begona yoki yoʻq buyurtma → 404 (authz.go).
func (a *api) getOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	o, err := a.st.Order(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

// noteInput — PATCH /orders/{id} body'si: yangi yetkazish izohi. Pointer: {} (kalit yoʻq) va
// {"note":""} (izohni oʻchirish) farqlansin — kalitsiz body izohni jim oʻchirib yubormasin.
type noteInput struct {
	Note *string `json:"note"`
}

// updateNote — PATCH /orders/{id}: izohni oʻzgartiradi, 204. Yoʻq buyurtma → 404.
func (a *api) updateNote(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in noteInput
	if !decodeJSON(w, r, &in) {
		return
	}
	f := fieldErrors{}
	if in.Note == nil {
		f["note"] = "majburiy (boʻsh satr — izohni oʻchirish)"
	} else {
		validateNote(*in.Note, f)
	}
	if len(f) > 0 {
		writeError(w, f)
		return
	}
	// VAZIFA (TASK.md, 2-qadam): store'ga buyurtma egasini ham bering.
	if err := a.st.UpdateNote(r.Context(), id, *in.Note); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// listUsers — GET /admin/users: hamma foydalanuvchilar. Faqat admin (route requireAdmin bilan, authz.go).
func (a *api) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.st.Users(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// pathID yoʻldagi {id}'ni oʻqiydi (GB-4.05). Son boʻlmasa 400 yozadi (errBadID) va false qaytaradi.
func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, fmt.Errorf("%w: %w", errBadID, err))
		return 0, false
	}
	return id, true
}
