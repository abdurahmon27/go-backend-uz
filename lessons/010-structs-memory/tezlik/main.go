// tezlik: n ta buyurtmaning toʻlangan summasi (standart: bir million). Ish bir xil,
// faqat xotiradan oʻqiladigan baytlar soni har xil.
// Ikkala variant navbatma-navbat 100 marta oʻlchanadi va eng yaxshi
// natija olinadi: kompyuterda boshqa ish ketayotgan boʻlsa ham shovqin kam.
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Buyurtma struct {
	Tolandi    bool
	Summa      int64
	Yetkazildi bool
}

type Ixcham struct {
	Summa      int64
	Tolandi    bool
	Yetkazildi bool
}

var natija int64 // kompilyator hisobni oʻchirib tashlamasin

func jami1(s []Buyurtma) (j int64) {
	for i := range s {
		if s[i].Tolandi {
			j += s[i].Summa
		}
	}
	return j
}

func jami2(s []Ixcham) (j int64) {
	for i := range s {
		if s[i].Tolandi {
			j += s[i].Summa
		}
	}
	return j
}

// vaqt: f ni bir marta ishlatib, ketgan vaqtni qaytaradi.
func vaqt(f func()) time.Duration {
	t := time.Now()
	f()
	return time.Since(t)
}

func main() {
	n := 1_000_000 // go run ./tezlik 100000 — boshqa oʻlcham
	if len(os.Args) > 1 {
		n, _ = strconv.Atoi(os.Args[1])
	}
	s1 := make([]Buyurtma, n)
	s2 := make([]Ixcham, n)
	for i := range n {
		s1[i] = Buyurtma{Tolandi: i%2 == 0, Summa: int64(i)}
		s2[i] = Ixcham{Tolandi: i%2 == 0, Summa: int64(i)}
	}
	eng1, eng2 := time.Hour, time.Hour
	for range 100 {
		eng1 = min(eng1, vaqt(func() { natija += jami1(s1) }))
		eng2 = min(eng2, vaqt(func() { natija += jami2(s2) }))
	}
	fmt.Printf("n = %d\n", n)
	fmt.Printf("Buyurtma: %.1f µs\n", eng1.Seconds()*1e6)
	fmt.Printf("Ixcham:   %.1f µs\n", eng2.Seconds()*1e6)
	fmt.Printf("Ixcham / Buyurtma: %.2f\n", eng2.Seconds()/eng1.Seconds())
}
