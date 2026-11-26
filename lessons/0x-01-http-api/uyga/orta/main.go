// 🟡 Oʻrta: GET /version qoʻshing: {"go": "go1.26.8", "revision": "<commit>"}.
// version() funksiyasi tayyor (binary ichidagi build maʼlumotini oʻqiydi — go version -m koʻrsatgan maʼlumotning oʻzi).
// Sizning ishingiz: uni GET /version route'iga ulash, keyin `go run .` va `go build` bilan solishtirish
// (qaysi birida revision "unknown"?). Javob: ../javob/orta/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

func version(w http.ResponseWriter, r *http.Request) {
	out := map[string]string{"go": "", "revision": "unknown"}
	if bi, ok := debug.ReadBuildInfo(); ok {
		out["go"] = bi.GoVersion
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" {
				out["revision"] = s.Value
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func main() {
	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	// TODO: GET /version → version
	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
