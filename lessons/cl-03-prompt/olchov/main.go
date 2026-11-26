// olchov: javoblar/ dagi saqlangan Claude javoblarini oʻlchaydi.
//
//	go run ./olchov          1–5-promptlar: bozor xatosi haqidagi savol
//	go run ./olchov -test    6–7-promptlar: holat() uchun table test
//	… -v                     + har bir javob boʻyicha tafsilot
//
// 1–5 uchun har bir javobda:
//   - soʻz: soʻzlar soni (boʻsh joy bilan ajratilgan)
//   - format: aynan uch qator: «Sabab:», «Tuzatish:», «Tekshirish:»
//   - tayyor: javobda koʻchirib qoʻyish mumkin boʻlgan errors.Is(err, ErrTopilmadi) bormi
//   - 404: «Tuzatish:» qatorini bozor/main.go dagi `if err == ErrTopilmadi {` oʻrniga qoʻyib,
//     dasturni ishga tushiramiz va GET /buyurtma/42 haqiqatan 404 qaytardimi, tekshiramiz
//
// 6–7 uchun javobdagi Go kod bloki main_test.go boʻladi va:
//   - soʻz: javobdagi soʻzlar soni
//   - nomlar: misoldagi (007-test-bench) nomlardan nechtasi olingan: jadval, tc.nom, «kutilgan»
//   - tuzatilgan: tuzatildi/main.go bilan go test oʻtdimi (PASS)
//   - eski: bozor/main.go (== bilan, xato bor) bilan go test xatoni ushladimi (FAIL)
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

const eski = "if err == ErrTopilmadi {"

func main() {
	v := flag.Bool("v", false, "har bir javob boʻyicha tafsilot")
	test := flag.Bool("test", false, "6–7-promptlar: table test")
	flag.Parse()

	bozor, err := os.ReadFile("bozor/main.go")
	if err != nil {
		fmt.Fprintln(os.Stderr, "bu papkadan ishga tushiring:", err)
		os.Exit(1)
	}
	tuzatildi, _ := os.ReadFile("tuzatildi/main.go")
	papkalar, _ := filepath.Glob("javoblar/*")
	sort.Strings(papkalar)

	if *test {
		qator("prompt", "soʻz (3 marta)", "nomlar", "tuzatilgan", "eski xato")
	} else {
		qator("prompt", "soʻz (3 marta)", "format", "tayyor", "404")
	}
	for _, p := range papkalar {
		if strings.Contains(p, "-test") != *test {
			continue
		}
		fayllar, _ := filepath.Glob(filepath.Join(p, "*.txt"))
		sort.Strings(fayllar)
		var sozlar, uslublar []string
		var a, b, c, sinaldi int
		for _, f := range fayllar {
			x, _ := os.ReadFile(f)
			javob := string(x)
			if *test {
				kod := goBlok(javob)
				u := uslub(kod)
				sozlar = append(sozlar, fmt.Sprint(len(strings.Fields(javob))))
				uslublar = append(uslublar, fmt.Sprint(u))
				pass := goTest(string(tuzatildi), kod) == "PASS"
				fail := goTest(string(bozor), kod) == "FAIL"
				if pass {
					a++
				}
				if fail {
					b++
				}
				if *v {
					fmt.Printf("  %s: uslub %d, tuzatilgan PASS %v, eski FAIL %v\n", f, u, pass, fail)
				}
				continue
			}
			n := len(strings.Fields(javob))
			sozlar = append(sozlar, fmt.Sprint(n))
			fok := formatda(javob)
			if fok {
				a++
			}
			tok := strings.Contains(javob, "errors.Is(err, ErrTopilmadi)")
			if tok {
				b++
			}
			natija := "—"
			if t, ok := tuzatish(javob); ok {
				sinaldi++
				natija = ishlat(string(bozor), t)
				if strings.HasPrefix(natija, "404") {
					c++
				}
			}
			if *v {
				fmt.Printf("  %s: %d soʻz, format %v, tayyor %v, %s\n", f, n, fok, tok, natija)
			}
		}
		k := len(fayllar)
		if *test {
			qator(filepath.Base(p), strings.Join(sozlar, " "), strings.Join(uslublar, " "),
				fmt.Sprintf("PASS %d/%d", a, k), fmt.Sprintf("FAIL %d/%d", b, k))
			continue
		}
		s404 := "—"
		if sinaldi > 0 {
			s404 = fmt.Sprintf("%d/%d", c, k)
		}
		qator(filepath.Base(p), strings.Join(sozlar, " "), fmt.Sprintf("%d/%d", a, k), fmt.Sprintf("%d/%d", b, k), s404)
	}
}

// qator: jadval qatori; kenglik baytda emas, harfda (ʻ ikki bayt)
func qator(ustunlar ...string) {
	kenglik := []int{13, 15, 7, 11, 0}
	var b strings.Builder
	for i, u := range ustunlar {
		b.WriteString(u)
		for n := utf8.RuneCountInString(u); n < kenglik[i]; n++ {
			b.WriteByte(' ')
		}
		if i < len(ustunlar)-1 {
			b.WriteByte(' ')
		}
	}
	fmt.Println(b.String())
}

// formatda: boʻsh boʻlmagan qatorlar aynan uchta va kerakli soʻzlar bilan boshlanadi
func formatda(javob string) bool {
	var q []string
	for _, l := range strings.Split(javob, "\n") {
		if strings.TrimSpace(l) != "" {
			q = append(q, strings.TrimSpace(l))
		}
	}
	return len(q) == 3 && strings.HasPrefix(q[0], "Sabab:") &&
		strings.HasPrefix(q[1], "Tuzatish:") && strings.HasPrefix(q[2], "Tekshirish:")
}

// tuzatish: «Tuzatish:» qatoridagi Go kodi (` belgilarisiz)
func tuzatish(javob string) (string, bool) {
	for _, l := range strings.Split(javob, "\n") {
		l = strings.TrimSpace(l)
		if s, ok := strings.CutPrefix(l, "Tuzatish:"); ok {
			s = strings.TrimSpace(strings.Trim(strings.TrimSpace(s), "`"))
			return s, s != ""
		}
	}
	return "", false
}

// ishlat: tuzatilgan kodni vaqtinchalik papkada ishga tushiradi, /buyurtma/42 holatini qaytaradi
func ishlat(kod, yangi string) string {
	if !strings.Contains(kod, eski) {
		return "eski qator topilmadi"
	}
	dir := vaqtinchalik(map[string]string{"main.go": strings.Replace(kod, eski, yangi, 1)})
	defer os.RemoveAll(dir)
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "kompilyatsiya xatosi"
	}
	for _, l := range strings.Split(string(out), "\n") {
		if s, ok := strings.CutPrefix(l, "GET /buyurtma/42 → "); ok {
			return strings.Fields(s)[0]
		}
	}
	return "natija yoʻq"
}

// goBlok: javobdagi birinchi ```go … ``` blok
func goBlok(javob string) string {
	_, s, ok := strings.Cut(javob, "```go\n")
	if !ok {
		return ""
	}
	s, _, _ = strings.Cut(s, "```")
	return s
}

// uslub: kursdagi testning nechta belgisi bor (007-test-bench/sodda/jadval_test.go)
func uslub(kod string) int {
	n := 0
	for _, b := range []string{"jadval", "tc.nom", "kutilgan"} {
		if strings.Contains(kod, b) {
			n++
		}
	}
	return n
}

// goTest: main.go + javobdagi main_test.go → "PASS", "FAIL" yoki "build"
func goTest(main, test string) string {
	dir := vaqtinchalik(map[string]string{"main.go": main, "main_test.go": test})
	defer os.RemoveAll(dir)
	cmd := exec.Command("go", "test", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	switch {
	case err == nil:
		return "PASS"
	case strings.Contains(string(out), "--- FAIL"):
		return "FAIL"
	default:
		return "build"
	}
}

func vaqtinchalik(fayllar map[string]string) string {
	dir, err := os.MkdirTemp("", "olchov")
	if err != nil {
		panic(err)
	}
	fayllar["go.mod"] = "module bozor\n\ngo 1.26\n"
	for nom, s := range fayllar {
		os.WriteFile(filepath.Join(dir, nom), []byte(s), 0o644)
	}
	return dir
}
