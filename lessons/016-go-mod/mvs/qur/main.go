// qur: manba/<nom>@<versiya>/ papkalaridan proxy/ papkasini quradi.
// proxy/ — GOPROXY protokolidagi oddiy fayllar (go help goproxy):
//
//	example.com/log/@v/list         mavjud versiyalar
//	example.com/log/@v/v1.3.0.info  {"Version": ..., "Time": ...}
//	example.com/log/@v/v1.3.0.mod   modulning go.mod fayli
//	example.com/log/@v/v1.3.0.zip   kod: example.com/log@v1.3.0/... ichida
//
// Ishga tushirish (016-go-mod/mvs papkasidan):
//
//	go run ./qur manba proxy
//	go run ./qur zaharli-manba zaharli
package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const domen = "example.com/"

func main() {
	if len(os.Args) != 3 {
		log.Fatal("ishlatish: go run ./qur <manba-papka> <proxy-papka>")
	}
	manba, proxy := os.Args[1], os.Args[2]
	dirs, err := os.ReadDir(manba)
	if err != nil {
		log.Fatal(err)
	}
	versiyalar := map[string][]string{} // modul → versiyalar
	for _, d := range dirs {
		nom, ver, ok := strings.Cut(d.Name(), "@")
		if !ok {
			continue
		}
		mod := domen + nom
		at := filepath.Join(proxy, mod, "@v")
		if err := os.MkdirAll(at, 0o755); err != nil {
			log.Fatal(err)
		}
		src := filepath.Join(manba, d.Name())
		gomod, err := os.ReadFile(filepath.Join(src, "go.mod"))
		if err != nil {
			log.Fatal(err)
		}
		yoz(filepath.Join(at, ver+".mod"), gomod)
		yoz(filepath.Join(at, ver+".info"),
			[]byte(fmt.Sprintf(`{"Version":%q,"Time":"2026-11-26T00:00:00Z"}`+"\n", ver)))
		zipla(filepath.Join(at, ver+".zip"), src, mod+"@"+ver+"/")
		versiyalar[mod] = append(versiyalar[mod], ver)
	}
	for mod, vs := range versiyalar {
		sort.Strings(vs)
		yoz(filepath.Join(proxy, mod, "@v", "list"), []byte(strings.Join(vs, "\n")+"\n"))
		fmt.Println(mod, vs)
	}
}

func yoz(nom string, b []byte) {
	if err := os.WriteFile(nom, b, 0o644); err != nil {
		log.Fatal(err)
	}
}

// zipla src papkasidagi fayllarni prefiks ostida zip'ga yozadi (go mod zip formati).
func zipla(nom, src, prefiks string) {
	f, err := os.Create(nom)
	if err != nil {
		log.Fatal(err)
	}
	zw := zip.NewWriter(f)
	fayllar, err := os.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range fayllar {
		b, err := os.ReadFile(filepath.Join(src, e.Name()))
		if err != nil {
			log.Fatal(err)
		}
		w, err := zw.CreateHeader(&zip.FileHeader{Name: prefiks + e.Name(), Method: zip.Deflate,
			Modified: time.Date(2026, 11, 26, 0, 0, 0, 0, time.UTC)})
		if err != nil {
			log.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			log.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
