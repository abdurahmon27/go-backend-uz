// taxmin — sizib chiqqan users jadvaliga oflayn hujumni shu kompyuterda oʻlchaydi (GB-6.03).
//
// Har qator («email hash») uchun nomzodlar.txt'dagi parollarni birma-bir sinaydi: sha256 satri uchun
// sha256 hisoblab solishtiradi, argon2id satri uchun bozor'ning oʻz password.Verify'ini chaqiradi —
// hujumchi ham aynan shu ishni qiladi, faqat bizning serverimizdan tashqarida, oʻz kompyuterida.
// Parollar va emaillar oʻylab topilgan (sizdi/ papkasini _video/sizdir yasaydi), haqiqiy sizgan roʻyxat emas.
//
//	go run ./cmd/taxmin sizdi/sha256.txt
//	go run ./cmd/taxmin sizdi/argon2id.txt
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/abdurahmon27/bozor/internal/password"
)

const wordlist = "nomzodlar.txt"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "ishlatish: go run ./cmd/taxmin sizdi/<fayl>.txt")
		os.Exit(2)
	}
	rows, err := lines(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	words, err := lines(wordlist)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	fmt.Printf("%s: %d qator · %s: %d parol\n", os.Args[1], len(rows), wordlist, len(words))

	var guesses int
	var ms0, ms1 runtime.MemStats
	runtime.ReadMemStats(&ms0)
	start := time.Now()
	for _, row := range rows {
		email, hash, _ := strings.Cut(row, " ")
		check := matcher(strings.TrimSpace(hash))
		t := time.Now()
		found := ""
		for i, w := range words {
			guesses++
			if check(w) {
				found = fmt.Sprintf("topildi: %s (%d-nomzod)", w, i+1)
				break
			}
		}
		if found == "" {
			found = "topilmadi"
		}
		fmt.Printf("%-19s %s · %s\n", email, found, short(time.Since(t)))
	}
	total := time.Since(start)
	runtime.ReadMemStats(&ms1)
	perGuess := float64(ms1.TotalAlloc-ms0.TotalAlloc) / float64(guesses)
	fmt.Printf("%d taxmin · %s · sekundiga %s · har taxmin %s xotira\n",
		guesses, short(total), group(float64(guesses)/total.Seconds()), size(perGuess))
}

// matcher sizgan bitta hash uchun «nomzod mos keladimi» funksiyasini qaytaradi.
// argon2id satri — bozor'ning password.Verify'i; qolgani — sha256'ning hex yozuvi (bir marta oʻqiladi).
func matcher(hash string) func(word string) bool {
	if strings.HasPrefix(hash, "$argon2id$") {
		return func(word string) bool {
			ok, err := password.Verify(word, hash)
			return err == nil && ok
		}
	}
	b, err := hex.DecodeString(hash)
	if err != nil || len(b) != sha256.Size {
		fmt.Fprintln(os.Stderr, "xato: hash na argon2id, na sha256:", hash)
		os.Exit(1)
	}
	want := [sha256.Size]byte(b)
	return func(word string) bool {
		return sha256.Sum256([]byte(word)) == want
	}
}

// lines fayldagi boʻsh boʻlmagan qatorlarni qaytaradi.
func lines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if s := strings.TrimSpace(sc.Text()); s != "" {
			out = append(out, s)
		}
	}
	return out, sc.Err()
}

// short — vaqtni oʻqishga qulay qilib yaxlitlaydi.
func short(d time.Duration) string {
	switch {
	case d < time.Millisecond:
		return d.Round(time.Microsecond).String()
	case d < time.Second:
		return d.Round(100 * time.Microsecond).String()
	default:
		return d.Round(100 * time.Millisecond).String()
	}
}

// group — butun sonni uch xonadan boʻshliq bilan yozadi: 12 900 000.
func group(x float64) string {
	s := fmt.Sprintf("%.0f", x)
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(' ')
		}
		b.WriteRune(c)
	}
	return b.String()
}

// size — baytlarni B / KiB / MiB bilan yozadi.
func size(b float64) string {
	switch {
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MiB", b/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1f KiB", b/(1<<10))
	default:
		return fmt.Sprintf("%.0f B", b)
	}
}
