// jadval: olchov.sh natijasini (stdin) mashhur "Latency numbers" jadvali bilan solishtiradi.
//
//	go run ./jadval < natija.txt            bizning jadval va mashhur jadval yonma-yon
//	go run ./jadval -tarqoq < natija.txt    har qator: 5 turning medianalari, eng kami … eng ko'pi
//	go run ./jadval -xy biz < natija.txt    "kategoriya ns" qatorlari (video chart sahnasi uchun)
//	go run ./jadval -xy mashhur             mashhur jadval qiymatlari, xuddi shu kategoriyalar
//	go run ./jadval -tur 5 -xy biz < natija.txt   faqat 5-tur (video terminal sahnalari shu turdan)
//
// Mashhur jadval: Jeff Dean, LADIS 2009 keynote, "Numbers Everyone Should Know" (24-slayd);
// SSD qatori: Jonas Bonér gist'i, 2012 ("~1GB/sec SSD"). Manbalar: README.md.
// DC = "Round trip within same datacenter", CA-NL = "Send packet CA->Netherlands->CA".
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type qator struct {
	nom     string  // jadvaldagi nom
	kat     string  // chart kategoriyasi ("" = chartda yo'q)
	mashhur float64 // mashhur jadvaldagi qiymat, ns (0 = u yerda yo'q)
	izoh    string  // mashhur qiymat nimani o'lchagan
	kalit   string  // natija.txt dagi nom (benchmark yoki olcha qatori)
	doimiy  float64 // DSA-A.01 da o'lchangan (qayta o'lchanmaydi), ns
}

// DSA-A.01 (../dsa-a01-cache/natija.txt, Zanjir, 10 o'lchov mediani, taskset -c 2):
// n=1024 (8 KiB, L1 ichida) 1.52 ns · n=32768 (256 KiB, L2) 4.23 ns · n=16777216 (128 MiB, RAM) 105.8 ns.
var qatorlar = []qator{
	{nom: "L1 keshdan o'qish", kat: "kesh", mashhur: 0.5, kalit: "dsa", doimiy: 1.52},
	{nom: "L2 keshdan o'qish", mashhur: 7, kalit: "dsa", doimiy: 4.23},
	{nom: "mutex lock/unlock", mashhur: 25, kalit: "BenchmarkMutex"},
	{nom: "RAM'dan o'qish", kat: "xotira", mashhur: 100, kalit: "dsa", doimiy: 105.8},
	{nom: "RAM'dan 1 MiB", mashhur: 250e3, kalit: "BenchmarkMegabayt"},
	{nom: "nusxadan 4 KiB", kalit: "nusxadan 4 KiB"},
	{nom: "SSD'dan 4 KiB", kat: "SSD", mashhur: 150e3, kalit: "SSD'dan 4 KiB"},
	{nom: "4 KiB yozish", kalit: "4 KiB yozish"},
	{nom: "4 KiB yozish + fsync", kat: "fsync", kalit: "4 KiB yozish + fsync"},
	{nom: "tarmoq borib-kelish", kat: "tarmoq", mashhur: 500e3, izoh: "DC", kalit: "Wi-Fi: borib-kelish"},
	{nom: "internet borib-kelish", kat: "internet", mashhur: 150e6, izoh: "CA-NL", kalit: "internet: borib-kelish"},
	{nom: "sahifa so'rovi", kalit: "sahifa so'rovi"},
}

func main() {
	tarqoq := flag.Bool("tarqoq", false, "har qator: turlar medianalari va oraliq")
	xy := flag.String("xy", "", "biz | mashhur: chart uchun \"kategoriya ns\" qatorlari")
	turN := flag.Int("tur", 0, "faqat shu tur (1…5): video ekranidagi terminal natijasi bilan bir xil")
	flag.Parse()

	if *xy == "mashhur" {
		for _, q := range qatorlar {
			if q.kat != "" && q.mashhur > 0 {
				fmt.Println(q.kat, xyson(q.mashhur))
			}
		}
		return
	}
	turlar := oqi()
	if *turN > 0 && *turN <= len(turlar) {
		turlar = turlar[*turN-1 : *turN]
	}
	biz := func(q qator) (med, lo, hi float64, ok bool) {
		if q.kalit == "dsa" {
			return q.doimiy, q.doimiy, q.doimiy, true
		}
		var m []float64
		for _, t := range turlar {
			if v, ok := t[q.kalit]; ok {
				m = append(m, median(v))
			}
		}
		if len(m) == 0 {
			return 0, 0, 0, false
		}
		return median(m), slices.Min(m), slices.Max(m), true
	}

	switch {
	case *xy == "biz":
		for _, q := range qatorlar {
			if v, _, _, ok := biz(q); ok && q.kat != "" {
				fmt.Println(q.kat, xyson(v))
			}
		}
	case *tarqoq:
		fmt.Printf("%-21s %8s %18s %6s\n", "amal", "mediana", "eng kami … ko'pi", "turlar")
		for _, q := range qatorlar {
			if q.kalit == "dsa" {
				continue
			}
			if v, lo, hi, ok := biz(q); ok {
				fmt.Printf("%-21s %8s %18s %6d\n", q.nom, qisqa(v), qisqa(lo)+" … "+qisqa(hi), len(turlar))
			}
		}
	default:
		fmt.Printf("%-21s %14s %10s %11s\n", "amal", "mashhur jadval", "bu noutbuk", "nisbat")
		for _, q := range qatorlar {
			v, _, _, ok := biz(q)
			if !ok {
				continue
			}
			m, n := "—", ""
			if q.mashhur > 0 {
				m = qisqa(q.mashhur)
				if q.izoh != "" {
					m += " " + q.izoh
				}
				n = nisbat(q.mashhur, v)
			}
			fmt.Printf("%-21s %14s %10s %11s\n", q.nom, m, qisqa(v), n)
		}
	}
}

// oqi: "# tur" qatorlari natijani turlarga bo'ladi; har turda kalit → qiymatlar (ns).
func oqi() []map[string][]float64 {
	var turlar []map[string][]float64
	tur := map[string][]float64{}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		s := sc.Text()
		f := strings.Fields(s)
		switch {
		case strings.HasPrefix(s, "# tur"):
			if len(tur) > 0 {
				turlar = append(turlar, tur)
			}
			tur = map[string][]float64{}
		case strings.HasPrefix(s, "Benchmark") && len(f) >= 4 && f[3] == "ns/op":
			nom, _, _ := strings.Cut(f[0], "-") // BenchmarkMutex-12 → BenchmarkMutex
			v, err := strconv.ParseFloat(f[2], 64)
			if err == nil {
				tur[nom] = append(tur[nom], v)
			}
		case len(f) >= 6 && f[len(f)-2] == "…":
			// olcha qatori: nom  marta  odatda  p10 … p90
			d, err := time.ParseDuration(f[len(f)-4])
			if _, e2 := strconv.Atoi(f[len(f)-5]); err == nil && e2 == nil {
				nom := strings.Join(f[:len(f)-5], " ")
				tur[nom] = append(tur[nom], float64(d))
			}
		}
	}
	if len(tur) > 0 {
		turlar = append(turlar, tur)
	}
	return turlar
}

func median(v []float64) float64 {
	s := slices.Clone(v)
	slices.Sort(s)
	if len(s)%2 == 1 {
		return s[len(s)/2]
	}
	return (s[len(s)/2-1] + s[len(s)/2]) / 2
}

func round3(v float64) string { return strconv.FormatFloat(v, 'g', 3, 64) }

// xyson: chart uchun son, 3 ta muhim raqam, darajasiz (1.37e+05 emas, 137000).
func xyson(v float64) string {
	r, _ := strconv.ParseFloat(round3(v), 64)
	return strconv.FormatFloat(r, 'f', -1, 64)
}

// qisqa: ns → "1.5 ns", "75 µs", "1.5 ms", "150 ms".
func qisqa(ns float64) string {
	for _, u := range []struct {
		b float64
		s string
	}{{1e9, "s"}, {1e6, "ms"}, {1e3, "µs"}, {1, "ns"}} {
		if ns >= u.b || u.s == "ns" {
			return round3(ns/u.b) + " " + u.s
		}
	}
	return ""
}

// nisbat: biz mashhur jadvaldan necha barobar tez yoki sekin.
func nisbat(mashhur, biz float64) string {
	r := mashhur / biz
	switch {
	case r >= 1.15:
		return strconv.FormatFloat(r, 'f', 1, 64) + "× tez"
	case r <= 1/1.15:
		return strconv.FormatFloat(1/r, 'f', 1, 64) + "× sekin"
	default:
		return "≈"
	}
}
