package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// daftar har bir mehmonni faylga yozadi va nechanchi ekanini aytadi.
func daftar(dir string) http.Handler {
	var mu sync.Mutex
	fayl := filepath.Join(dir, "mehmonlar.txt")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		f, err := os.OpenFile(fayl, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		fmt.Fprintln(f, "mehmon")
		b, _ := os.ReadFile(fayl)
		fmt.Fprintf(w, "Salom! Siz %d-mehmonsiz\n", bytes.Count(b, []byte("\n")))
	})
}
