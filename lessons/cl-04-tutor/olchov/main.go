// olchov: javoblar/ dagi saqlangan Claude javoblarini oʻlchaydi.
//
//	go run ./olchov         chat: javob mashinasi va ustoz (1–3-xabar)
//	go run ./olchov -code   Claude Code: default va Learning uslubi
//	… -v                    + har bir javobdagi soʻzlar soni
//
// Har bir qator — bitta usul, har biri 3 marta sinalgan:
//   - tayyor kod: javobda (yoki Claude yozgan faylda) tayyor tuzatish bormi:
//     errors.As(err, &…) yoki errors.AsType[…] — chatda faqat ```go kod bloki ichida sanaladi
//   - savol:     javob oʻquvchiga savol beradimi (kod bloklaridan tashqarida «?» bor)
//   - ishladi:   tayyor kodni bozor/main.go ga qoʻyib ishga tushiramiz: narx=-5 → 400 boʻldimi
//   - TODO(human) / Insight: Claude Code'ning Learning uslubi belgilari
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type usul struct{ nom, glob string }

var usullar = []usul{
	{"javob mashinasi", "javoblar/1-javob/*-1.txt"},
	{"ustoz, 1-xabar", "javoblar/2-ustoz/*-1.txt"},
	{"ustoz, 2-xabar", "javoblar/2-ustoz/*-2.txt"},
	{"ustoz, 3-xabar", "javoblar/2-ustoz/*-3.txt"},
	{"Code: default", "javoblar/uslub-default/*.txt"},
	{"Code: Learning", "javoblar/uslub-Learning/*.txt"},
}

func main() {
	code := flag.Bool("code", false, "Claude Code qatorlari")
	v := flag.Bool("v", false, "har bir usul uchun soʻzlar soni")
	flag.Parse()
	if _, err := os.Stat("bozor/main.go"); err != nil {
		fmt.Fprintln(os.Stderr, "cl-04-tutor papkasidan ishga tushiring:", err)
		os.Exit(1)
	}
	sarlavha := ""
	if *code {
		sarlavha = "TODO(human) Insight"
	}
	chiqar("usul", "tayyor kod", "savol", "ishladi", sarlavha)
	for _, u := range usullar {
		if strings.HasPrefix(u.nom, "Code") != *code {
			continue
		}
		fayllar, _ := filepath.Glob(u.glob)
		sort.Strings(fayllar)
		var sozlar []string
		var tayyor, savol, ishladi, todo, insight int
		for _, f := range fayllar {
			b, _ := os.ReadFile(f)
			matn := string(b)
			sozlar = append(sozlar, fmt.Sprint(len(strings.Fields(matn))))
			if savolBor(matn) {
				savol++
			}
			if strings.Contains(matn, "Insight") {
				insight++
			}
			var kod string
			if diff, err := os.ReadFile(strings.TrimSuffix(f, ".txt") + ".diff"); err == nil {
				// Claude Code faylni oʻzi tahrirlagan: diff'ni asl faylga qoʻllaymiz
				kod = yamoq(diff)
				if strings.Contains(string(diff), "TODO(human)") {
					todo++
				}
			} else {
				kod = holatBilan(goBloklar(matn))
			}
			if tuzatishBor(kod) {
				tayyor++
				if ishga(kod) {
					ishladi++
				}
			}
		}
		n := len(fayllar)
		ish := "—"
		if tayyor > 0 {
			ish = fmt.Sprintf("%d/%d", ishladi, n)
		}
		qoshimcha := ""
		if strings.HasPrefix(u.nom, "Code") {
			qoshimcha = fmt.Sprintf("%-11s %d/%d", fmt.Sprintf("%d/%d", todo, n), insight, n)
		}
		chiqar(u.nom, fmt.Sprintf("%d/%d", tayyor, n), fmt.Sprintf("%d/%d", savol, n), ish, qoshimcha)
		if *v {
			fmt.Println("  soʻzlar:", strings.Join(sozlar, " "))
		}
	}
}

func chiqar(usul, tayyor, savol, ishladi, qoshimcha string) {
	q := fmt.Sprintf("%-16s %-11s %-6s %-8s %s", usul, tayyor, savol, ishladi, qoshimcha)
	fmt.Println(strings.TrimRight(q, " "))
}

// savolBor: kod bloklaridan tashqarida «?» bormi
func savolBor(matn string) bool {
	for i, q := range strings.Split(matn, "```") {
		if i%2 == 0 && strings.Contains(q, "?") {
			return true
		}
	}
	return false
}

// goBloklar: javobdagi ```go … ``` bloklari
func goBloklar(matn string) []string {
	var bloklar []string
	qismlar := strings.Split(matn, "```go\n")
	for _, q := range qismlar[1:] {
		if i := strings.Index(q, "```"); i >= 0 {
			bloklar = append(bloklar, q[:i])
		}
	}
	return bloklar
}

func tuzatishBor(kod string) bool {
	return strings.Contains(kod, "errors.As(err, &") || strings.Contains(kod, "errors.AsType[")
}

// holatBilan: javobdagi `func holat` blokini bozor/main.go dagi holat oʻrniga qoʻyadi
func holatBilan(bloklar []string) string {
	asl, _ := os.ReadFile("bozor/main.go")
	for _, b := range bloklar {
		i := strings.Index(b, "func holat(")
		if i < 0 {
			continue
		}
		yangi := b[i:]
		if j := strings.Index(yangi, "\n}\n"); j >= 0 {
			yangi = yangi[:j+3]
		}
		s := string(asl)
		bosh := strings.Index(s, "func holat(")
		oxir := bosh + strings.Index(s[bosh:], "\n}\n") + 3
		return s[:bosh] + yangi + s[oxir:]
	}
	return ""
}

func yamoq(diff []byte) string {
	dir, _ := os.MkdirTemp("", "olchov")
	defer os.RemoveAll(dir)
	out := filepath.Join(dir, "main.go")
	cmd := exec.Command("patch", "-s", "-o", out, "bozor/main.go")
	cmd.Stdin = strings.NewReader(string(diff))
	if err := cmd.Run(); err != nil {
		return ""
	}
	b, _ := os.ReadFile(out)
	return string(b)
}

// ishga: kodni vaqtinchalik modulda ishga tushiradi; narx=-5 qatori 400 bilan chiqsa — ishladi
func ishga(kod string) bool {
	dir, _ := os.MkdirTemp("", "olchov")
	defer os.RemoveAll(dir)
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module bozor\n\ngo 1.26\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte(kod), 0o644)
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	b, err := cmd.CombinedOutput()
	return err == nil && strings.Contains(string(b), "narx=-5 → 400")
}
