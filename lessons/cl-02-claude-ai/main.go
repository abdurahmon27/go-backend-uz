// tekshir: matnni chatga (Claude.ai, Project fayllari) yuborishdan oldin unda kalit,
// parol yoki shaxsiy maʼlumotga oʻxshash narsa bor-yoʻqligini qidiradi.
//
//	go run . fayl.log           topilmalar roʻyxati (qiymatning faqat boshi va oxiri)
//	go run . papka/             papkadagi matn fayllari + .env fayllarining nomi
//	go run . -yashir fayl.log   matnning oʻzi, topilmalar [YASHIRIN:tur] bilan
//	cat fayl.log | go run .     stdin
//
//	go run . -qatiy papka/      topilma boʻlsa chiqish kodi 1 (skript va git hook uchun)
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type qoida struct {
	tur string              // topilma turi
	sir bool                // true: kalit/parol, false: shaxsiy maʼlumot
	re  *regexp.Regexp      // qidiruv namunasi
	ok  func(s string) bool // qoʻshimcha tekshiruv (nil = shart yoʻq)
}

// Tartib muhim: bir joyga ikki qoida tushsa, roʻyxatda oldin turgani yutadi.
var qoidalar = []qoida{
	{"anthropic-kalit", true, regexp.MustCompile(`sk-ant-[A-Za-z0-9_-]{20,}`), nil},
	{"github-token", true, regexp.MustCompile(`\b(?:gh[pousr]_|github_pat_)[A-Za-z0-9_]{20,}`), nil},
	{"aws-kalit", true, regexp.MustCompile(`\b(?:AKIA|ASIA)[A-Z0-9]{16}\b`), nil},
	{"maxfiy-kalit", true, regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`), nil},
	{"url-parol", true, regexp.MustCompile(`[a-z][a-z0-9+.-]*://[^/\s:@]+:[^/\s@]+@`), nil},
	{"parol", true, regexp.MustCompile(`(?i)\b[a-z_]*(?:password|passwd|pass|pwd|parol|secret|api_?key|token)["']?\s*[:=]\s*["']?[^\s"',}]{6,}`), nil},
	{"email", false, regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`), nil},
	{"telefon", false, regexp.MustCompile(`\+998[ -]?\(?\d{2}\)?[ -]?\d{3}[ -]?\d{2}[ -]?\d{2}\b`), nil},
	{"karta", false, regexp.MustCompile(`\b\d{4}[ -]?\d{4}[ -]?\d{4}[ -]?\d{4}\b`), luhn},
}

// Nomaʼlum formatdagi kalit: uzun, tasodifiy koʻrinadigan soʻz.
var soz = regexp.MustCompile(`[A-Za-z0-9]{24,}`)

type topilma struct {
	joy        string // fayl:qator
	bosh, oxir int    // qator ichidagi bayt oraligʻi
	tur        string
	sir        bool
	qiymat     string
}

// qatorTekshir bitta qatordagi topilmalarni qaytaradi (bir-birini qoplamaydi).
func qatorTekshir(joy, qator string) []topilma {
	var t []topilma
	band := func(b, o int) bool {
		for _, x := range t {
			if b < x.oxir && x.bosh < o {
				return true
			}
		}
		return false
	}
	for _, q := range qoidalar {
		for _, m := range q.re.FindAllStringIndex(qator, -1) {
			s := qator[m[0]:m[1]]
			if band(m[0], m[1]) || (q.ok != nil && !q.ok(s)) {
				continue
			}
			t = append(t, topilma{joy, m[0], m[1], q.tur, q.sir, s})
		}
	}
	for _, m := range soz.FindAllStringIndex(qator, -1) {
		s := qator[m[0]:m[1]]
		if !band(m[0], m[1]) && tasodifiy(s) {
			t = append(t, topilma{joy, m[0], m[1], "nomaʼlum-kalit", true, s})
		}
	}
	sort.Slice(t, func(i, j int) bool { return t[i].bosh < t[j].bosh })
	return t
}

// tasodifiy: kamida 24 belgi; katta harf, kichik harf va kamida 2 ta raqam aralash; Shannon
// entropiyasi 4 bit/belgidan koʻp (Go'dagi uzun nomlar, masalan http2ErrCodeEnhanceYourCalm, oʻtmaydi).
// Git hash (40 ta hex) bu sinovga kirmaydi: katta harfi yoʻq, 16 xil belgida esa maksimum 4 bit.
// Narxi: 24 belgili tasodifiy kalitlarning ~13 foizini oʻtkazib yuboradi (go test -run Olchov -v).
func tasodifiy(s string) bool {
	var katta, kichik bool
	raqam := 0
	soni := map[rune]int{}
	for _, r := range s {
		katta = katta || (r >= 'A' && r <= 'Z')
		kichik = kichik || (r >= 'a' && r <= 'z')
		if r >= '0' && r <= '9' {
			raqam++
		}
		soni[r]++
	}
	return katta && kichik && raqam >= 2 && entropiya(soni, len(s)) > 4.0
}

func entropiya(soni map[rune]int, n int) float64 {
	var c []int
	for _, x := range soni {
		c = append(c, x)
	}
	sort.Ints(c) // map tartibi tasodifiy: yigʻindi har safar bir xil boʻlsin
	h := 0.0
	for _, x := range c {
		p := float64(x) / float64(n)
		h -= p * math.Log2(p)
	}
	return h
}

// luhn: karta raqamining oxirgi raqami tekshiruv raqami (ISO/IEC 7812-1).
// 16 xonali har qanday son (buyurtma raqami, vaqt) karta emas. 0000 0000 … ham emas.
func luhn(s string) bool {
	raqamlar := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
	if strings.Count(raqamlar, raqamlar[:1]) == len(raqamlar) {
		return false // hammasi bir xil raqam
	}
	sum, ikki := 0, false
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c < '0' || c > '9' {
			continue
		}
		d := int(c - '0')
		if ikki {
			if d *= 2; d > 9 {
				d -= 9
			}
		}
		sum += d
		ikki = !ikki
	}
	return sum%10 == 0
}

// yashir qiymatning boshi va oxirini qoldiradi: sk-a…Qx. Qisqa qiymat toʻliq yashirinadi.
func yashir(s string) string {
	r := []rune(s)
	if len(r) < 16 {
		return strings.Repeat("*", len(r))
	}
	return string(r[:4]) + "…" + string(r[len(r)-2:])
}

var xeshFayl = map[string]bool{"go.sum": true, "package-lock.json": true, "yarn.lock": true, "pnpm-lock.yaml": true}

// envFayl: .env, .env.local … (lekin namunasi .env.example emas)
func envFayl(nom string) bool {
	if nom == ".env.example" || nom == ".env.sample" {
		return false
	}
	return nom == ".env" || strings.HasPrefix(nom, ".env.")
}

func oqi(joy string, r io.Reader, yoz io.Writer) ([]topilma, error) {
	var hammasi []topilma
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for n := 1; sc.Scan(); n++ {
		qator := sc.Text()
		t := qatorTekshir(fmt.Sprintf("%s:%d", joy, n), qator)
		hammasi = append(hammasi, t...)
		if yoz != nil { // -yashir: qatorni topilmalarsiz chiqaramiz
			for i := len(t) - 1; i >= 0; i-- {
				qator = qator[:t[i].bosh] + "[YASHIRIN:" + t[i].tur + "]" + qator[t[i].oxir:]
			}
			fmt.Fprintln(yoz, qator)
		}
	}
	return hammasi, sc.Err()
}

func papka(yol string) ([]topilma, error) {
	var hammasi []topilma
	err := filepath.WalkDir(yol, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if xeshFayl[d.Name()] {
			return nil // xesh'lar roʻyxati: tasodifiy koʻrinadi, lekin sir emas
		}
		if envFayl(d.Name()) {
			hammasi = append(hammasi, topilma{joy: p, tur: "env-fayl", sir: true, qiymat: d.Name()})
		}
		if info, err := d.Info(); err != nil || info.Size() > 1<<20 {
			return nil // 1 MB dan katta fayl: oʻtkazib yuboramiz
		}
		b, err := os.ReadFile(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, "oʻqib boʻlmadi:", err)
			return nil
		}
		if bytes.IndexByte(b[:min(len(b), 512)], 0) >= 0 {
			return nil // binar fayl
		}
		t, err := oqi(p, bytes.NewReader(b), nil)
		hammasi = append(hammasi, t...)
		return err
	})
	return hammasi, err
}

func main() {
	yashirMode := flag.Bool("yashir", false, "matnni chiqarish, topilmalarni [YASHIRIN:tur] bilan almashtirib")
	qatiy := flag.Bool("qatiy", false, "topilma boʻlsa chiqish kodi 1")
	flag.Parse()

	var t []topilma
	var err error
	switch {
	case flag.NArg() == 0:
		t, err = oqi("stdin", os.Stdin, chiqish(*yashirMode))
	case *yashirMode:
		var f *os.File
		if f, err = os.Open(flag.Arg(0)); err == nil {
			t, err = oqi(flag.Arg(0), f, os.Stdout)
			f.Close()
		}
	default:
		for _, yol := range flag.Args() {
			var x []topilma
			x, err = papka(yol) // fayl boʻlsa ham ishlaydi
			t = append(t, x...)
			if err != nil {
				break
			}
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(2)
	}
	if *yashirMode {
		fmt.Fprintf(os.Stderr, "%d ta qiymat yashirildi\n", len(t))
		return
	}
	sirlar := 0
	for _, x := range t {
		turi := "shaxsiy"
		if x.sir {
			turi, sirlar = "sir", sirlar+1
		}
		q := yashir(x.qiymat)
		if x.tur == "env-fayl" {
			q = "Project'ga yuklamang"
		}
		fmt.Printf("%-22s %-7s %-15s %s\n", x.joy, turi, x.tur, q)
	}
	if len(t) == 0 {
		fmt.Println("topilma yoʻq (bu «xavfsiz» degani emas: faqat shu qoidalar tekshirildi)")
		return
	}
	fmt.Printf("jami %d: %d ta sir, %d ta shaxsiy maʼlumot. Chatga yuborishdan oldin olib tashlang.\n",
		len(t), sirlar, len(t)-sirlar)
	if *qatiy {
		os.Exit(1)
	}
}

// chiqish: -yashir boʻlsa matn stdout'ga yoziladi, aks holda hech qayerga.
func chiqish(ha bool) io.Writer {
	if ha {
		return os.Stdout
	}
	return nil
}
