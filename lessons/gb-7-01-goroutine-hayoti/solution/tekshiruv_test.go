package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// Bu fayl vazifani tekshiradi (TASK.md). Uni oʻzgartirmang.
//
// Soxta fetch: har yetkazuvchi oʻz vaqtida javob beradi (tarmoq oʻrniga time.Sleep) va
// tugaganini atomic hisoblagichga yozadi. Keyingi yetkazuvchilar tezroq javob beradi, shuning uchun
// natijani append bilan yigʻsangiz, tartib buziladi.

var testSuppliers = []string{"chorsu", "oloy", "siyob", "urgut", "qoʻqon"}

var testDelays = map[string]time.Duration{
	"chorsu": 150 * time.Millisecond,
	"oloy":   120 * time.Millisecond,
	"siyob":  90 * time.Millisecond,
	"urgut":  60 * time.Millisecond,
	"qoʻqon": 30 * time.Millisecond,
}

var testCounts = map[string]int{"chorsu": 40, "oloy": 120, "siyob": 5, "urgut": 17, "qoʻqon": 64}

// fakeFetch qaytaradi: fetch funksiyasi va tugagan chaqiruvlar soni hisoblagichi.
func fakeFetch(scale float64) (func(string) int, *atomic.Int32) {
	var done atomic.Int32
	f := func(name string) int {
		time.Sleep(time.Duration(float64(testDelays[name]) * scale))
		done.Add(1)
		return testCounts[name]
	}
	return f, &done
}

func want() []int {
	w := make([]int, len(testSuppliers))
	for i, s := range testSuppliers {
		w[i] = testCounts[s]
	}
	return w
}

type fetchAllFunc func([]string, func(string) int) []int

func checkFetchAll(t *testing.T, name string, fn fetchAllFunc) {
	t.Helper()
	n := int32(len(testSuppliers))

	// 1) Egasi kutadimi? fetchAll qaytganda hamma fetch tugagan boʻlishi kerak.
	f, done := fakeFetch(1)
	start := time.Now()
	got := fn(testSuppliers, f)
	elapsed := time.Since(start)
	finished := done.Load()
	if finished != n {
		t.Fatalf("%s qaytdi, lekin %d/%d ta fetch hali tugamagan edi: goroutine'larni boshlagan funksiya "+
			"ularni kutishi kerak (wg.Wait()), aks holda natija hali yozilmagan boʻladi", name, n-finished, n)
	}

	// 2) Tartib: counts[i] = fetch(suppliers[i]).
	if !slices.Equal(got, want()) {
		t.Fatalf("%s natijasi %v, kutilgan %v (suppliers tartibida). Natijani append bilan yigʻmang: "+
			"har goroutine oʻz indeksiga yozsin, counts[i] = fetch(name)", name, got, want())
	}

	// 3) Bir vaqtda: ketma-ket 150+120+90+60+30 = 450ms, bir vaqtda ~150ms (eng sekini).
	if elapsed > 300*time.Millisecond {
		t.Fatalf("%s %d ta yetkazuvchini %v'da oʻqidi: ular ketma-ket ishlayapti (450ms atrofida). "+
			"Bir vaqtda boʻlsa, eng sekini, ya'ni ~150ms atrofida tugaydi: har fetch'ni wg.Go bilan alohida "+
			"goroutine'da boshlang", name, n, elapsed.Round(time.Millisecond))
	}

	// 4) Taxmin emas, tugashini kutish: hamma yetkazuvchi 10 barobar tez boʻlsa, fetchAll ham tez qaytsin.
	f, done = fakeFetch(0.1)
	start = time.Now()
	fn(testSuppliers, f)
	elapsed = time.Since(start)
	if done.Load() != n {
		t.Fatalf("%s tez yetkazuvchilarda ham hammasini kutishi kerak", name)
	}
	if elapsed > 80*time.Millisecond {
		t.Fatalf("yetkazuvchilar ~15ms'da javob berdi, %s esa %v kutdi: time.Sleep bilan vaqtni taxmin "+
			"qilmang, wg.Wait() aynan oxirgi goroutine tugaguncha kutadi", name, elapsed.Round(time.Millisecond))
	}
}

func TestVazifaFetchAll(t *testing.T) {
	checkFetchAll(t, "fetchAll", fetchAll)
}

func TestVazifaFetchAllBosh(t *testing.T) {
	got := fetchAll(nil, func(string) int { panic("boʻsh roʻyxatda fetch chaqirilmasligi kerak") })
	if len(got) != 0 {
		t.Fatalf("fetchAll(nil) = %v, kutilgan boʻsh natija", got)
	}
}

// 🟡 ixtiyoriy: fetchAllAddDone yozilmagan boʻlsa (nil qaytarsa), SKIP.
func TestSariqAddDone(t *testing.T) {
	f, _ := fakeFetch(0.01)
	if fetchAllAddDone(testSuppliers, f) == nil {
		t.Skip("🟡 fetchAllAddDone hali yozilmagan (ixtiyoriy, TASK.md 2-qadam)")
	}
	body := funcSource(t, "fetchAllAddDone")
	for _, s := range []string{".Add(1)", ".Done()", "go func"} {
		if !strings.Contains(body, s) {
			t.Fatalf("fetchAllAddDone eski shaklda boʻlsin: ichida %q topilmadi (wg.Add(1), go func() { defer wg.Done() … }())", s)
		}
	}
	if strings.Contains(body, ".Go(") {
		t.Fatalf("fetchAllAddDone'da wg.Go ishlatmang: bu mashq Add/Done shaklini oʻqish va yozish uchun")
	}
	checkFetchAll(t, "fetchAllAddDone", fetchAllAddDone)
}

// funcSource — main.go'dagi funksiya tanasining matni (tekshiruv uchun; go/parser bilan topiladi).
func funcSource(t *testing.T, name string) string {
	t.Helper()
	src, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range file.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Name.Name == name {
			return string(src[fset.Position(fd.Body.Pos()).Offset:fset.Position(fd.Body.End()).Offset])
		}
	}
	t.Fatalf("main.go'da %s topilmadi", name)
	return ""
}
