// Chuqurroq: scratch image'da sertifikatlar yoʻq, tashqi HTTPS soʻrov yiqiladi.
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("https://go.dev")
	if err != nil {
		fmt.Println("xato:", err)
		os.Exit(1)
	}
	resp.Body.Close()
	fmt.Println("status:", resp.Status)
}
