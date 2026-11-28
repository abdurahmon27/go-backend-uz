package main

// Route'lar (GB-4.05) — himoyalangan bozor (GB-6.08). Uch guruh:
//
//	ochiq        GET /health, GET /products…, POST /users (roʻyxatdan oʻtish), POST /login (rate limit), POST /logout
//	login kerak  requireUser (auth.go) bilan: session tirik boʻlmasa handler'gacha yetmaydi — 401
//	admin        /admin/… — alohida mux, BIR MARTA oʻralgan: requireUser (kim?) ichida requireAdmin (admin'mi?) → 403
//
// Body'li route'lar requireJSON bilan (GB-6.01, 415), u requireUser'ning ichida: avval kimligi, keyin body.

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/abdurahmon27/bozor/internal/store"
)

// api — handler'lar bazani va login limitlarini shu struct orqali koʻradi (GB-5.10, GB-6.07).
type api struct {
	st         *store.Store
	loginIP    *limiters // GB-6.07: bitta IP'dan login urinishlari
	loginEmail *limiters // GB-6.07: bitta email'ga login urinishlari (IP'lar almashsa ham)
}

// newMux — hamma route'lar shu yerda: main ham, testlar ham bir xil mux'ni ishlatadi.
func newMux(st *store.Store) *http.ServeMux {
	a := &api{
		st:         st,
		loginIP:    newLimiters(12*time.Second, 5), // 5 ta urinish, keyin har 12 soniyada bittadan
		loginEmail: newLimiters(time.Minute, 10),   // 10 ta urinish, keyin har daqiqada bittadan
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /products", a.listProducts)
	mux.HandleFunc("GET /products/{id}", a.getProduct)
	mux.Handle("POST /users", requireJSON(http.HandlerFunc(a.register))) // GB-6.03
	mux.Handle("POST /login", a.loginIP.limitByIP(requireJSON(http.HandlerFunc(a.login))))
	mux.HandleFunc("POST /logout", a.logout) // requireUser'siz: eskirgan cookie ham tozalansin (GB-6.04)

	mux.Handle("GET /me", a.requireUser(http.HandlerFunc(a.me)))
	mux.Handle("POST /logout-all", a.requireUser(http.HandlerFunc(a.logoutAll)))
	mux.Handle("GET /orders", a.requireUser(http.HandlerFunc(a.listOrders)))
	mux.Handle("POST /orders", a.requireUser(requireJSON(http.HandlerFunc(a.createOrder))))
	mux.Handle("GET /orders/{id}", a.requireUser(http.HandlerFunc(a.getOrder)))
	mux.Handle("PATCH /orders/{id}", a.requireUser(requireJSON(http.HandlerFunc(a.updateNote))))

	// Admin: alohida mux. U bir marta oʻraladi, shuning uchun unga qoʻshilgan har route — keyin qoʻshilgani
	// ham — himoyalangan boʻladi (OWASP API5:2023: bitta, tekshirish oson authorization moduli). Yoʻq yoʻl
	// ham avval 401/403'dan oʻtadi: admin boʻlmagan foydalanuvchi /admin/ ostida nima borligini bilmaydi.
	admin := http.NewServeMux()
	admin.HandleFunc("GET /admin/users", a.listUsers)
	admin.HandleFunc("GET /admin/orders", a.allOrders)
	admin.Handle("POST /admin/products", requireJSON(http.HandlerFunc(a.createProduct)))
	admin.HandleFunc("DELETE /admin/products/{id}", a.deleteProduct)
	mux.Handle("/admin/",
		a.requireUser(requireAdmin(muxErrors(admin))))
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
		writeError(w, r, err)
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
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// productInput — POST /admin/products body'si: faqat mijoz berishi mumkin boʻlgan maydonlar (ID'ni baza beradi).
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

// createProduct — POST /admin/products (faqat admin): decodeJSON → validateProduct (422) → store.Add → 201 va yangi mahsulot.
func (a *api) createProduct(w http.ResponseWriter, r *http.Request) {
	var in productInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if f := validateProduct(in); len(f) > 0 {
		writeError(w, r, f)
		return
	}
	p, err := a.st.Add(r.Context(), in.Name, in.Price, in.Stock)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/products/%d", p.ID))
	writeJSON(w, http.StatusCreated, p)
}

// deleteProduct — DELETE /admin/products/{id} (faqat admin): 204 (body yoʻq), 400, 404, buyurtmalarda boʻlsa 409.
func (a *api) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := a.st.Delete(r.Context(), id); err != nil {
		writeError(w, r, err)
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
		writeError(w, r, f)
		return
	}
	items := make([]store.OrderItem, len(in.Items))
	for i, it := range in.Items {
		items[i] = store.OrderItem{ProductID: it.ProductID, Quantity: it.Quantity}
	}
	o, err := a.st.PlaceOrder(r.Context(), currentUser(r).ID, in.Note, items)
	if err != nil {
		writeError(w, r, err)
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
	o, err := a.st.Order(r.Context(), id, currentUser(r).ID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

// noteInput — PATCH /orders/{id} body'si: yangi yetkazish izohi. Pointer: {} (kalit yoʻq) va
// {"note":""} (izohni boʻshatish) farqlansin — kalitsiz body izohni jim oʻchirib yubormasin.
type noteInput struct {
	Note *string `json:"note"`
}

// updateNote — PATCH /orders/{id}: izohni oʻzgartiradi, 204. Begona yoki yoʻq buyurtma → 404 (vazifa, TASK.md).
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
		f["note"] = "majburiy (boʻsh satr — izohni boʻshatish)"
	} else {
		validateNote(*in.Note, f)
	}
	if len(f) > 0 {
		writeError(w, r, f)
		return
	}
	if err := a.st.UpdateNote(r.Context(), id, currentUser(r).ID, *in.Note); err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// listOrders — GET /orders: faqat session egasining buyurtmalari, ID tartibida; boʻlmasa [] (VAZIFA, TASK.md).
// Egasi currentUser'dan: mijoz uni yoʻlda ham, body'da ham bera olmaydi.
func (a *api) listOrders(w http.ResponseWriter, r *http.Request) {
	list, err := a.st.Orders(r.Context(), currentUser(r).ID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// allOrders — GET /admin/orders: hamma buyurtmalar (VAZIFA). Rol tekshiruvi handler'da emas: admin mux
// uni bir marta oʻragan.
func (a *api) allOrders(w http.ResponseWriter, r *http.Request) {
	list, err := a.st.AllOrders(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// listUsers — GET /admin/users: hamma foydalanuvchilar (admin mux'da).
func (a *api) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.st.Users(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// pathID yoʻldagi {id}'ni oʻqiydi (GB-4.05). Son boʻlmasa 400 yozadi (errBadID) va false qaytaradi.
func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, r, fmt.Errorf("%w: %w", errBadID, err))
		return 0, false
	}
	return id, true
}
