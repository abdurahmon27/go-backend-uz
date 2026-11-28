package main

// Authorization — «nima qila olasiz» (GB-6.06). Authentication (requireUser, auth.go) faqat kim ekaningizni
// aniqlaydi. Ruxsatning ikki qavati:
//
//	obyekt   «bu buyurtma sizniki?» — har store metodida, SQL shartida: WHERE id = $1 AND user_id = $2.
//	         Begona buyurtma → 404, xuddi yoʻq buyurtma kabi (RFC 9110 §15.5.4: mavjudligini yashirish).
//	funksiya «bu amal sizning rolingizga ochiq?» — requireAdmin middleware'i, /admin/… route'larida → 403.
//
// OWASP API Security Top 10 2023: API1 Broken Object Level Authorization, API5 Broken Function Level Authorization.

import (
	"errors"
	"net/http"

	"github.com/abdurahmon27/bozor/internal/store"
)

// errForbidden — foydalanuvchi kim ekani maʼlum, lekin bu amal uning roliga ruxsat etilmagan (problemFor → 403).
var errForbidden = errors.New("bu amal uchun ruxsat yoʻq")

// requireAdmin — faqat admin rolidagi foydalanuvchini oʻtkazadi, qolganlarga 403.
// requireUser'ning ICHIDA turishi shart: foydalanuvchini context'ga requireUser qoʻyadi.
// Tashqarida tursa currentUser panic qiladi → 500 (GB-6.04): jim oʻtkazib yuborishdan xavfsizroq.
func requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if currentUser(r).Role != store.RoleAdmin {
			writeError(w, errForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
