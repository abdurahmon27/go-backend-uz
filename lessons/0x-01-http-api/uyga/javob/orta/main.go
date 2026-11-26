// 🟡 Javob: binary ichidagi build maʼlumoti (go version -m) endpoint orqali.
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
	http.HandleFunc("GET /version", version)
	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
