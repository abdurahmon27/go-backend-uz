// fanout: bitta sahifa N ta so'rovni parallel yuboradi va eng sekin javobni kutadi.
// Har bir so'rov yuztadan bitta holatda 1 soniyadan sekin. Sahifalarning qanchasi sekin?
// So'rovlar bir-biriga bog'liq emas (mustaqil) deb olinadi: hayotda bu har doim ham rost emas.
package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
)

const sahifalar = 100_000

// sorov: 99 % hollarda 10..50 ms, 1 % hollarda 1..3 soniya.
func sorov(r *rand.Rand) float64 {
	if r.Float64() < 0.01 {
		return 1000 + 2000*r.Float64()
	}
	return 10 + 40*r.Float64()
}

func main() {
	N := 100
	if len(os.Args) > 1 {
		if os.Args[1] == "jadval" { // 1 - 0.99^N, turli N uchun
			for _, n := range []int{1, 10, 20, 50, 100, 200, 500} {
				fmt.Printf("N = %3d: %5.1f %% sahifa sekin\n", n, 100*(1-math.Pow(0.99, float64(n))))
			}
			return
		}
		n, err := strconv.Atoi(os.Args[1])
		if err != nil || n < 1 {
			fmt.Fprintln(os.Stderr, "foydalanish: go run ./fanout [N >= 1 | jadval]")
			os.Exit(2)
		}
		N = n
	}
	r := rand.New(rand.NewPCG(1, 2))
	sekinSorov, sekinSahifa := 0, 0
	for range sahifalar {
		eng := 0.0
		for range N {
			t := sorov(r)
			if t >= 1000 {
				sekinSorov++
			}
			eng = max(eng, t) // sahifa eng sekin javobni kutadi
		}
		if eng >= 1000 {
			sekinSahifa++
		}
	}
	fmt.Printf("so'rovlar: 1 s dan sekini %5.2f %%\n", 100*float64(sekinSorov)/float64(sahifalar*N))
	fmt.Printf("sahifalar (%d so'rovli): 1 s dan sekini %4.1f %%\n", N, 100*float64(sekinSahifa)/sahifalar)
	fmt.Printf("qo'lda hisob: 1 - 0.99^%d = %.1f %%\n", N, 100*(1-math.Pow(0.99, float64(N))))
}
