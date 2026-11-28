// begona — «boshqa sayt» (GB-6.07): bozor'ga hujum qiladigan ikki sahifa. Bitta server, ikki nom bilan ochiladi:
//
//	http://localhost:9090   — bozor (localhost:8765) bilan bir SITE, lekin boshqa ORIGIN (port boshqa)
//	http://127.0.0.1:9090   — boshqa site ham (host boshqa)
//
//	/         forma: tugma bosilsa brauzer POST /orders/41/cancel'ni bozor'ga yuboradi (CSRF)
//	/oqish    skript: GET /me'ni cookie bilan soʻraydi va javobni OʻQISHGA urinadi (CORS)
//
// Faqat laboratoriya uchun: bozor API'si ishlab turganda brauzerda oching (README).
package main

import (
	"fmt"
	"net/http"
	"os"
)

// formaPage — CSRF: oddiy HTML forma. Skript yoʻq; foydalanuvchi tugmani bosadi, brauzer esa formani
// bozor'ga yuboradi va bozor'ning cookie'larini (SameSite ruxsat bersa) oʻzi qoʻshadi.
const formaPage = `<!doctype html>
<meta charset="utf-8">
<title>Sovgʻa</title>
<form method="post" action="http://localhost:8765/orders/41/cancel">
  <button>Sovgʻani olish</button>
</form>
`

// oqishPage — CORS: fetch bozor'ga cookie bilan (credentials: "include") soʻrov yuboradi. Soʻrov serverga
// yetadi; javobni sahifaga berish-bermaslikni brauzer bozor'ning CORS sarlavhalariga qarab hal qiladi.
const oqishPage = `<!doctype html>
<meta charset="utf-8">
<title>Oʻqish</title>
<pre id="natija">kutilmoqda…</pre>
<script>
fetch("http://localhost:8765/me", {credentials: "include"})
  .then(r => r.text())
  .then(t => { natija.textContent = "oʻqildi: " + t })
  .catch(e => { natija.textContent = "oʻqib boʻlmadi: " + e })
</script>
`

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, formaPage)
	})
	mux.HandleFunc("GET /oqish", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, oqishPage)
	})
	fmt.Println("begona sayt: http://localhost:9090 va http://127.0.0.1:9090")
	if err := http.ListenAndServe(":9090", mux); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
