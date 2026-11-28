package main

// Brauzer va bozor (GB-6.07). Brauzer cookie'ni soʻrovga oʻzi qoʻshadi — soʻrovni qaysi sahifa
// boshlaganidan qatʼi nazar. Shuning uchun server ikki narsani alohida tekshiradi:
//
//	protectCSRF  holat oʻzgartiradigan soʻrov (POST, PUT, PATCH, DELETE) boshqa origin'dan kelganmi?
//	             Kelgan boʻlsa — 403, handler ishlamaydi (CSRF himoyasi, http.CrossOriginProtection).
//	cors         boshqa origin'dagi sahifa javobni OʻQIY oladimi? Faqat allowedOrigins'dagilar.
//	             Oddiy soʻrovni (GET, forma POST) CORS toʻxtatmaydi: u faqat brauzerga «javobni shu sahifaga ber»
//	             deydi. Oddiy boʻlmagan soʻrovdan (JSON body, PATCH, DELETE) oldin brauzer preflight (OPTIONS)
//	             yuboradi; ruxsat boʻlmasa haqiqiy soʻrov yuborilmaydi.
//
// origin = sxema + host + port (http://localhost:8765). site = sxema + roʻyxatdan oʻtgan domen
// (eTLD+1), port hisobga olinmaydi: http://localhost:9090 bozor bilan bir site, lekin boshqa origin.
// SameSite cookie atributi site'ga qaraydi, Sec-Fetch-Site va Origin sarlavhalari esa origin'ni ham ajratadi.

import (
	"errors"
	"fmt"
	"net/http"
)

// allowedOrigins — bozor API'siga boshqa origin'dan murojaat qila oladigan sahifalar: bozor'ning
// oʻz veb ilovasi (laboratoriyada http://localhost:3000). Bu roʻyxat ikki joyda ishlaydi: CORS javob
// sarlavhalari va CSRF tekshiruvining ishonchli origin'lari. Hech qachon «har qanday Origin» emas.
var allowedOrigins = []string{"http://localhost:3000"}

// errCrossOrigin — soʻrov boshqa origin'dagi sahifadan kelgan (problemFor → 403).
var errCrossOrigin = errors.New("boshqa saytdan kelgan soʻrov")

// protectCSRF — Go 1.25'dan beri standart kutubxonadagi http.CrossOriginProtection bilan:
//   - GET, HEAD, OPTIONS — har doim oʻtadi (ular holatni oʻzgartirmasligi kerak);
//   - Sec-Fetch-Site: same-origin yoki none (manzil qatoriga yozilgan) — oʻtadi; boshqa qiymat — rad;
//   - Sec-Fetch-Site yoʻq, Origin bor — Origin'ning host'i Host sarlavhasiga teng boʻlsa oʻtadi;
//   - ikkalasi ham yoʻq — Go buni brauzer emas deb hisoblaydi (curl, mobil ilova, boshqa server): oʻtadi (fail-open);
//   - AddTrustedOrigin'dagi origin'lar — oʻtadi.
//
// Origin ↔ Host solishtiruvida sxema (http/https) koʻrilmaydi: production'da HTTPS + HSTS (GB-8.06).
//
// Brauzer bo'lmagan mijozda begona cookie yoʻq: CSRF — aynan brauzer cookie'ni oʻzi qoʻshgani uchun xavf.
// Rad sababi (Check'ning xatosi) log'ga, mijozga — bizning matn (GB-6.02).
func protectCSRF(next http.Handler) http.Handler {
	cop := http.NewCrossOriginProtection()
	for _, o := range allowedOrigins {
		if err := cop.AddTrustedOrigin(o); err != nil {
			panic(err) // roʻyxat kodimizda: notoʻgʻri yozilgan boʻlsa server ishga tushmasin
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := cop.Check(r); err != nil {
			writeError(w, fmt.Errorf("%w: %w", errCrossOrigin, err))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// cors — CORS sarlavhalari faqat roʻyxatdagi origin'larga (Fetch standarti, "CORS protocol").
// Origin roʻyxatda boʻlmasa hech narsa qoʻshilmaydi: soʻrov baribir bajariladi, lekin brauzer javobni
// sahifaga bermaydi. Origin'ni tekshirmasdan qaytarish (Allow-Credentials bilan) — har qanday sahifaga
// foydalanuvchi maʼlumotini oʻqishga ruxsat: bu xato README'da koʻrsatilgan.
// Vary: Origin — javob Origin'ga qarab farq qiladi: kesh bir origin'ning javobini boshqasiga bermasin.
func cors(next http.Handler) http.Handler {
	allowed := map[string]bool{}
	for _, o := range allowedOrigins {
		allowed[o] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		origin := r.Header.Get("Origin")
		if !allowed[origin] {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Preflight: brauzer haqiqiy soʻrovdan oldin OPTIONS bilan ruxsat soʻraydi (JSON body, PATCH, DELETE…).
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Max-Age", "600") // brauzer javobni 10 daqiqa eslab qoladi
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
