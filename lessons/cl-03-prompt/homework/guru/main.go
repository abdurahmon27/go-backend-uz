// Guru: test javoblarida nechta Test funksiyasi va nechta jadval holati borligini sanaydi.
//
//	go run ./homework/guru     (cl-03-prompt papkasidan)
//
// Holat = jadvaldagi bitta qator: `{"nom", …}` yoki `nom: "…"` bilan boshlangan qator.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	testFunc = regexp.MustCompile(`(?m)^func Test`)
	holat    = regexp.MustCompile(`(?m)^\s*(\{\s*"|\{?\s*nom:\s*")`)
)

func main() {
	papkalar, _ := filepath.Glob("javoblar/*-test*")
	sort.Strings(papkalar)
	for _, p := range papkalar {
		fayllar, _ := filepath.Glob(filepath.Join(p, "*.txt"))
		sort.Strings(fayllar)
		var funklar, holatlar []string
		for _, f := range fayllar {
			b, err := os.ReadFile(f)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			funklar = append(funklar, fmt.Sprint(len(testFunc.FindAll(b, -1))))
			holatlar = append(holatlar, fmt.Sprint(len(holat.FindAll(b, -1))))
		}
		fmt.Printf("%-13s Test funksiya: %s · holat: %s\n", filepath.Base(p),
			strings.Join(funklar, " "), strings.Join(holatlar, " "))
	}
}
