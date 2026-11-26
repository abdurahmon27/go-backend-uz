// abtest: tugma testini 10 000 marta simulyatsiya qiladi. Har guruhda n kishi, eski tugmada 5 % sotib oladi.
//
//	A/A: yangi tugma AYNAN BIR XIL (5 %). Har qanday "yutuq" shovqin.
//	A/B: yangi tugma haqiqatan +10 % (5.5 %). Test buni necha marta topadi?
//
// "Farq bor" qarori: ikki ulushning z-testi, |z| > 1.96 (95 % ishonch), NIST e-Handbook 7.3.3.
package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
)

const (
	p          = 0.05 // eski tugma konversiyasi
	tajribalar = 10_000
	ishchilar  = 10 // har biri o'z generatori bilan 1000 tajriba: natija har safar bir xil
)

// sotdi: n kishidan nechtasi q ehtimol bilan sotib oldi.
func sotdi(r *rand.Rand, n int, q float64) int {
	k := 0
	for range n {
		if r.Float64() < q {
			k++
		}
	}
	return k
}

type natija struct{ plus10, minus10, farqBor int }

// sina: yangi tugmaning haqiqiy konversiyasi q bo'lganda 10 000 tajriba.
func sina(n int, q float64) natija {
	var mu sync.Mutex
	var jami natija
	var wg sync.WaitGroup
	for w := range ishchilar {
		wg.Go(func() {
			r := rand.New(rand.NewPCG(1, uint64(w)))
			var t natija
			for range tajribalar / ishchilar {
				a, b := sotdi(r, n, p), sotdi(r, n, q)
				nisbiy := float64(b-a) / float64(max(a, 1))
				if nisbiy >= 0.10 {
					t.plus10++
				}
				if nisbiy <= -0.10 {
					t.minus10++
				}
				pa, pb, pp := float64(a)/float64(n), float64(b)/float64(n), float64(a+b)/float64(2*n)
				if se := math.Sqrt(pp * (1 - pp) * 2 / float64(n)); se > 0 && math.Abs(pb-pa)/se > 1.96 {
					t.farqBor++
				}
			}
			mu.Lock()
			jami.plus10, jami.minus10, jami.farqBor = jami.plus10+t.plus10, jami.minus10+t.minus10, jami.farqBor+t.farqBor
			mu.Unlock()
		})
	}
	wg.Wait()
	return jami
}

func foiz(k int) float64 { return 100 * float64(k) / tajribalar }

func main() {
	n := 2000 // har bir guruhda
	if len(os.Args) > 1 {
		k, err := strconv.Atoi(os.Args[1])
		if err != nil || k < 100 {
			fmt.Fprintln(os.Stderr, "foydalanish: go run ./abtest [n >= 100]")
			os.Exit(2)
		}
		n = k
	}
	aa, ab := sina(n, p), sina(n, p*1.10)
	fmt.Printf("har guruhda %d kishi, %d tajriba\n", n, tajribalar)
	fmt.Println("A/A, tugmalar bir xil:")
	fmt.Printf("  yangi tugma >= +10 %%: %4.1f %%\n", foiz(aa.plus10))
	fmt.Printf("  yangi tugma <= -10 %%: %4.1f %%\n", foiz(aa.minus10))
	fmt.Printf("  z-test \"farq bor\":    %4.1f %%\n", foiz(aa.farqBor))
	fmt.Println("A/B, yangi tugma rostdan +10 %:")
	fmt.Printf("  z-test \"farq bor\":    %4.1f %%\n", foiz(ab.farqBor))
	fmt.Printf("qo'lda hisob: shovqin ~ sqrt(2*0.05*0.95/%d) = %.2f foiz punkti\n", n, 100*math.Sqrt(2*p*(1-p)/float64(n)))
}
