// jadval: `go test -bench` natijasini (stdin) ekranga sig'adigan jadval va grafikka aylantiradi.
// Kesh o'lchamlari shu kompyuterning /sys/devices/system/cpu/cpuN/cache papkasidan o'qiladi.
//
//	go run ./jadval < natija.txt              grafik: n = 2^10 … 2^26, ns/element
//	go run ./jadval -daraja < natija.txt      har kesh darajasi ichidagi o'lchamlar: eng kami … eng ko'pi
//	go run ./jadval -ikki < natija.txt        2×2: tartib (ketma-ket/aralash) × bog'liqlik (mustaqil/zanjir)
//	go run ./jadval -n 16777216 < natija.txt  bitta o'lcham: vaqt va nisbat
//	go run ./jadval -tarqoq < natija.txt      o'lchovlar orasidagi farq
//	go run ./jadval -kesh                     shu kompyuter keshlari
//	go run ./jadval -xy Zanjir < natija.txt   "n median" qatorlari (video chart sahnasi uchun)
//
// Linux bo'lmasa (sysfs yo'q) kesh o'lchamlarini qo'lda bering: -l1 48K -l2 1280K -l3 12M.
//
//	go run ./jadval -median < qadam.txt     istalgan benchmark: median va [min … max]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// usullar: grafik va jadvallardagi asosiy uchtasi. TZanjir (tartibli zanjir) faqat -ikki va -tarqoq da.
var usullar = []string{"Ketma", "Aralash", "Zanjir"}

var nomlar = map[string]string{"Ketma": "ketma", "Aralash": "aralash", "Zanjir": "zanjir", "TZanjir": "t.zanjir"}

type kalit struct {
	usul string
	n    int
}

type olchov struct{ nsEl, nsOp []float64 }

type kesh struct {
	nom  string
	bayt int
}

func main() {
	cpu := flag.Int("cpu", 2, "kesh o'lchamlari qaysi yadrodan o'qiladi (taskset -c bilan bir xil)")
	daraja := flag.Bool("daraja", false, "har kesh darajasi ichidagi o'lchamlar: eng kami … eng ko'pi")
	bittaN := flag.Int("n", 0, "faqat shu n uchun vaqt va nisbat")
	tarqoq := flag.Bool("tarqoq", false, "o'lchovlar orasidagi tarqoqlik")
	medianFlag := flag.Bool("median", false, "istalgan benchmark natijasi: median va oraliq")
	keshFlag := flag.Bool("kesh", false, "faqat shu kompyuter keshlari: P-yadro (-cpu) va E-yadro (-ecpu)")
	ecpu := flag.Int("ecpu", 8, "tejamkor (E) yadro raqami")
	xy := flag.String("xy", "", "shu usul uchun \"n median_ns_el\" qatorlari")
	ikki := flag.Bool("ikki", false, "2×2: tartib × bog'liqlik, eng katta o'lchamda")
	l1 := flag.String("l1", "", "L1d hajmi qo'lda (masalan 48K), sysfs bo'lmasa")
	l2 := flag.String("l2", "", "L2 hajmi qo'lda")
	l3 := flag.String("l3", "", "L3 hajmi qo'lda")
	flag.Parse()
	if *medianFlag {
		medianJadval()
		return
	}
	if *keshFlag {
		keshChiqar(*cpu, *ecpu)
		return
	}

	data, nlar := oqi()
	keshlar := keshlarniOqi(*cpu)
	if *l1 != "" {
		keshlar = nil
		for i, h := range []string{*l1, *l2, *l3} {
			if h != "" {
				keshlar = append(keshlar, kesh{[]string{"L1d", "L2", "L3"}[i], hajmOqi(h)})
			}
		}
	}
	switch {
	case *xy != "":
		for _, n := range nlar {
			if o := data[kalit{*xy, n}]; o != nil {
				fmt.Println(n, strconv.FormatFloat(median(o.nsEl), 'g', 4, 64))
			}
		}
	case *ikki:
		ikkiJadval(data, nlar)
	case *daraja:
		darajaJadval(data, nlar, keshlar)
	case *bittaN > 0:
		bittaJadval(data, *bittaN)
	case *tarqoq:
		tarqoqJadval(data, nlar)
	default:
		grafik(data, nlar, keshlar)
	}
}

var qator = regexp.MustCompile(`^Benchmark(\w+)/n=(\d+)(?:-\d+)?\s+\d+\s+([\d.]+) ns/op\s+([\d.]+) ns/el`)

func oqi() (map[kalit]*olchov, []int) {
	data := map[kalit]*olchov{}
	var nlar []int
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		m := qator.FindStringSubmatch(sc.Text())
		if m == nil {
			continue
		}
		n, _ := strconv.Atoi(m[2])
		op, _ := strconv.ParseFloat(m[3], 64)
		el, _ := strconv.ParseFloat(m[4], 64)
		k := kalit{m[1], n}
		if data[k] == nil {
			data[k] = &olchov{}
		}
		data[k].nsOp = append(data[k].nsOp, op)
		data[k].nsEl = append(data[k].nsEl, el)
		if !slices.Contains(nlar, n) {
			nlar = append(nlar, n)
		}
	}
	slices.Sort(nlar)
	return data, nlar
}

func keshlarniOqi(cpu int) []kesh {
	var k []kesh
	papka := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cache", cpu)
	dirs, _ := filepath.Glob(filepath.Join(papka, "index*"))
	for _, d := range dirs {
		tur := oqiFayl(d, "type")
		if tur == "Instruction" {
			continue
		}
		kb, _ := strconv.Atoi(strings.TrimSuffix(oqiFayl(d, "size"), "K")) // "48K", "1280K", "12288K"
		nom := "L" + oqiFayl(d, "level")
		if tur == "Data" {
			nom += "d"
		}
		if kb > 0 {
			k = append(k, kesh{nom, kb << 10})
		}
	}
	return k
}

func keshChiqar(cpus ...int) {
	for _, l := range strings.Split(oqiFayl("/proc", "cpuinfo"), "\n") {
		if strings.HasPrefix(l, "model name") {
			fmt.Println(strings.TrimSpace(strings.SplitN(l, ":", 2)[1]))
			break
		}
	}
	for i, c := range cpus {
		var q []string
		for _, k := range keshlarniOqi(c) {
			q = append(q, k.nom+" "+hajm(k.bayt))
		}
		qator := oqiFayl(fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cache/index0", c), "coherency_line_size")
		fmt.Printf("cpu%d (%s): %s · qator %s bayt\n", c, []string{"P", "E"}[i], strings.Join(q, " · "), qator)
	}
}

// hajmOqi: "48K" -> 49152, "12M" -> 12582912.
func hajmOqi(h string) int {
	birlik := map[byte]int{'K': 1 << 10, 'M': 1 << 20}[h[len(h)-1]]
	v, _ := strconv.Atoi(strings.TrimRight(h, "KM"))
	if birlik == 0 {
		return v
	}
	return v * birlik
}

func oqiFayl(d, f string) string {
	b, _ := os.ReadFile(filepath.Join(d, f))
	return strings.TrimSpace(string(b))
}

func median(v []float64) float64 {
	s := slices.Sorted(slices.Values(v))
	if len(s)%2 == 1 {
		return s[len(s)/2]
	}
	return (s[len(s)/2-1] + s[len(s)/2]) / 2
}

// hajm: 8192 -> "8K", 1310720 -> "1.25M"
func hajm(b int) string {
	switch {
	case b >= 1<<20 && b%(1<<20) == 0:
		return strconv.Itoa(b>>20) + "M"
	default:
		return strconv.FormatFloat(float64(b)/(1<<10), 'f', -1, 64) + "K"
	}
}

func son(v float64) string {
	switch {
	case v >= 10:
		return strconv.FormatFloat(v, 'f', 0, 64)
	case v >= 1:
		return strconv.FormatFloat(v, 'f', 1, 64)
	default:
		return strconv.FormatFloat(v, 'f', 2, 64)
	}
}

func daraja2(n int) string { return "2^" + strconv.Itoa(int(math.Log2(float64(n)))) }

func nom(u string) string { return nomlar[u] }

// grafik: har usul bitta qator, har o'lcham ikki belgi, balandlik log shkalada.
// "|" — slice (8n bayt) navbatdagi keshdan kattalashgan joy.
func grafik(data map[kalit]*olchov, nlar []int, keshlar []kesh) {
	bars := []rune("▁▂▃▄▅▆▇█")
	lo, hi := math.Inf(1), math.Inf(-1)
	for k, o := range data {
		if !slices.Contains(usullar, k.usul) {
			continue
		}
		m := median(o.nsEl)
		lo, hi = min(lo, m), max(hi, m)
	}
	var leg []string
	for _, k := range keshlar {
		leg = append(leg, k.nom+" "+hajm(k.bayt))
	}
	fmt.Printf("ns/element, n = %s … %s, | = %s\n", daraja2(nlar[0]), daraja2(nlar[len(nlar)-1]), strings.Join(leg, " · "))
	for _, u := range usullar {
		var b strings.Builder
		ki := 0
		for _, n := range nlar {
			for ki < len(keshlar) && 8*n > keshlar[ki].bayt {
				b.WriteRune('|')
				ki++
			}
			m := median(data[kalit{u, n}].nsEl)
			h := int(math.Round((math.Log(m) - math.Log(lo)) / (math.Log(hi) - math.Log(lo)) * float64(len(bars)-1)))
			b.WriteRune(bars[h])
			b.WriteRune(bars[h])
		}
		first, last := median(data[kalit{u, nlar[0]}].nsEl), median(data[kalit{u, nlar[len(nlar)-1]}].nsEl)
		fmt.Printf("%-8s %s  %s → %s\n", nom(u), b.String(), son(first), son(last))
	}
}

// darajaJadval: har kesh darajasi uchun slice (8·n bayt) shu keshga sig'adigan, lekin oldingisiga
// sig'maydigan barcha o'lchamlar; har usul uchun ularning medianlari: eng kami … eng ko'pi.
// Oxirgi qator: hech bir keshga sig'maydigan o'lchamlar (operativ xotira, DRAM).
func darajaJadval(data map[kalit]*olchov, nlar []int, keshlar []kesh) {
	fmt.Printf("%-5s %-10s %10s %10s %10s\n", "qayer", "slice", "ketma", "aralash", "zanjir")
	quyi := 0
	chegaralar := append(slices.Clone(keshlar), kesh{"DRAM", math.MaxInt})
	for _, k := range chegaralar {
		var ichida []int
		for _, n := range nlar {
			if 8*n > quyi && 8*n <= k.bayt {
				ichida = append(ichida, n)
			}
		}
		quyi = k.bayt
		if len(ichida) == 0 {
			continue
		}
		oraliq := hajm(8 * ichida[0])
		if len(ichida) > 1 {
			oraliq += "…" + hajm(8*ichida[len(ichida)-1])
		}
		fmt.Printf("%-5s %-10s", strings.TrimSuffix(k.nom, "d"), oraliq)
		for _, u := range usullar {
			var m []float64
			for _, n := range ichida {
				m = append(m, median(data[kalit{u, n}].nsEl))
			}
			q := son(slices.Min(m))
			if son(slices.Max(m)) != q {
				q += "…" + son(slices.Max(m))
			}
			fmt.Printf(" %10s", q)
		}
		fmt.Println()
	}
}

// ikkiJadval: eng katta o'lchamda 2×2. Tartib prefetcher'ga, bog'liqlik parallel kutishga ta'sir qiladi.
func ikkiJadval(data map[kalit]*olchov, nlar []int) {
	n := nlar[len(nlar)-1]
	g := func(u string) string {
		if o := data[kalit{u, n}]; o != nil {
			return son(median(o.nsEl))
		}
		return "—"
	}
	fmt.Printf("n = %s (%s), ns/element    mustaqil   zanjir\n", daraja2(n), hajm(8*n))
	fmt.Printf("%-27s %8s %8s\n", "ketma-ket tartib", g("Ketma"), g("TZanjir"))
	fmt.Printf("%-27s %8s %8s\n", "aralash tartib", g("Aralash"), g("Zanjir"))
}

func vaqt(ns float64) string {
	switch {
	case ns >= 1e9:
		return son(ns/1e9) + " s"
	case ns >= 1e6:
		return son(ns/1e6) + " ms"
	case ns >= 1e3:
		return son(ns/1e3) + " µs"
	default:
		return son(ns) + " ns"
	}
}

func bittaJadval(data map[kalit]*olchov, n int) {
	fmt.Printf("n = %s, slice %s, %d o'lchov mediani\n", daraja2(n), hajm(8*n), len(data[kalit{"Ketma", n}].nsOp))
	baza := median(data[kalit{"Ketma", n}].nsOp)
	for _, u := range usullar {
		o := data[kalit{u, n}]
		op := median(o.nsOp)
		fmt.Printf("%-8s %8s  %5s ns/el  %4.0f×\n", nom(u), vaqt(op), son(median(o.nsEl)), op/baza)
	}
}

// kvartil: chiziqli interpolyatsiya bilan p-kvantil (0..1).
func kvartil(v []float64, p float64) float64 {
	s := slices.Sorted(slices.Values(v))
	x := p * float64(len(s)-1)
	i := int(x)
	if i+1 >= len(s) {
		return s[len(s)-1]
	}
	return s[i] + (x-float64(i))*(s[i+1]-s[i])
}

// tarqoqJadval: har o'lchamda IQR/median (o'rtadagi 50 % natija) va (max-min)/median;
// ikkalasining ham barcha o'lchamlar bo'yicha mediani va eng kattasi.
func tarqoqJadval(data map[kalit]*olchov, nlar []int) {
	fmt.Printf("tarqoqlik, %d o'lchov: IQR/median · (max-min)/median\n", len(data[kalit{"Ketma", nlar[0]}].nsEl))
	for _, u := range append(slices.Clone(usullar), "TZanjir") {
		if data[kalit{u, nlar[0]}] == nil {
			continue
		}
		var iqr, tola []float64
		for _, n := range nlar {
			o := data[kalit{u, n}].nsEl
			m := median(o)
			iqr = append(iqr, (kvartil(o, 0.75)-kvartil(o, 0.25))/m*100)
			tola = append(tola, (slices.Max(o)-slices.Min(o))/m*100)
		}
		fmt.Printf("%-8s IQR %4.1f%% (eng ko'pi %4.1f%%) · max-min %4.1f%% (%5.1f%%)\n",
			nom(u), median(iqr), slices.Max(iqr), median(tola), slices.Max(tola))
	}
}

var umumiy = regexp.MustCompile(`^Benchmark(\S+?)(?:-\d+)?\s+\d+\s+([\d.]+) ns/op(?:\s+([\d.]+) ns/el)?`)

// medianJadval: har benchmark nomi uchun ns/op mediani, [eng kami … eng ko'pi] va (bo'lsa) ns/el.
func medianJadval() {
	var nomlar []string
	op, el := map[string][]float64{}, map[string][]float64{}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		m := umumiy.FindStringSubmatch(sc.Text())
		if m == nil {
			continue
		}
		if op[m[1]] == nil {
			nomlar = append(nomlar, m[1])
		}
		v, _ := strconv.ParseFloat(m[2], 64)
		op[m[1]] = append(op[m[1]], v)
		if m[3] != "" {
			e, _ := strconv.ParseFloat(m[3], 64)
			el[m[1]] = append(el[m[1]], e)
		}
	}
	nre := regexp.MustCompile(`n=(\d+)`)
	for _, n := range nomlar {
		o := op[n]
		qisqa := nre.ReplaceAllStringFunc(n, func(m string) string {
			v, _ := strconv.Atoi(m[2:])
			return "n=" + daraja2(v)
		})
		qator := fmt.Sprintf("%-16s %7s  [%s … %s]", qisqa, vaqt(median(o)), vaqt(slices.Min(o)), vaqt(slices.Max(o)))
		if e := el[n]; e != nil {
			qator += fmt.Sprintf("  %s ns/el", son(median(e)))
		}
		fmt.Println(qator)
	}
}
