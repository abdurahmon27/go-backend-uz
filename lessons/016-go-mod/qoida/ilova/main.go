package main

import (
	"fmt"
	"net/http/httptest"

	"example.com/kutubxona"
)

// ilova: go 1.26. kutubxona: go 1.21.
func main() {
	w := httptest.NewRecorder()
	kutubxona.Router().ServeHTTP(w, httptest.NewRequest("GET", "/salom", nil))
	fmt.Println("kutubxona router:", w.Code)
	fmt.Println("kutubxona sikl:  ", kutubxona.Domenlar())
	fmt.Println("ilova sikl:      ", domenlar())
}

// domenlar kutubxonadagi bilan bir xil sikl, lekin go 1.26 modulda.
func domenlar() string {
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
