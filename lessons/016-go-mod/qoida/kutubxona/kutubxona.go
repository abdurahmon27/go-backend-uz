// Package kutubxona: go.mod'ida "go 1.21". Uni go 1.26 dagi ilova ishlatadi.
package kutubxona

import (
	"fmt"
	"net/http"
)

// Router "GET /salom" ni roʻyxatga oladi.
// Qaysi router ishlashini (GODEBUG) ilovaning go qatori hal qiladi.
func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /salom", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "salom!")
	})
	return mux
}

// Domenlar 5-darsdagi sikl: har qadamda d ning manzilini saqlaydi.
// Til versiyasi (-lang) shu modulning go qatoridan olinadi: go 1.21.
func Domenlar() string {
	var navbat []*string
	for _, d := range []string{"a.uz", "b.uz", "c.uz"} {
		navbat = append(navbat, &d)
	}
	s := ""
	for _, p := range navbat {
		s += *p + " "
	}
	return s
}
