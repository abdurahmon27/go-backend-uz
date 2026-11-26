// Uy vazifasi javoblari: go run ./homework
package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
)

func main() {
	// Oson: 20 so'rovli sahifa, har so'rov 100 tadan 1 holatda 1 s dan sekin.
	fmt.Printf("oson:  1 - 0.99^20 = %.1f %% (tekshiring: go run ./fanout 20)\n", 100*(1-math.Pow(0.99, 20)))

	// O'rta: har so'rov 1000 tadan 1 holatda 1 s dan sekin bo'lsa, 100 so'rovli sahifa.
	fmt.Printf("o'rta: 1 - 0.999^100 = %.1f %% (100 tadan 1 bo'lganda 63.4 %% edi)\n", 100*(1-math.Pow(0.999, 100)))

	// Qiyin: Kohavi formulasi n = 16*sigma^2/delta^2, 5 % konversiya, +5 % nisbiy (5 % -> 5.25 %).
	sigma2, delta := 0.05*0.95, 0.05*0.05
	fmt.Printf("qiyin: n = 16*0.0475/0.0025^2 = %.0f kishi har guruhda\n", 16*sigma2/(delta*delta))

	// Guru: kesh relizida p99 < 100 ms bo'lishi uchun kesh qancha so'rovni ushlashi kerak?
	for _, hit := range []float64{0.96, 0.98, 0.99, 0.995} {
		fmt.Printf("guru:  kesh %.1f %% -> p99 %4.0f ms\n", 100*hit, p99(hit))
	}
}

func p99(hit float64) float64 {
	r := rand.New(rand.NewPCG(1, 2))
	xs := make([]float64, 100_000)
	for i := range xs {
		m := 400.0
		if r.Float64() < hit {
			m = 8
		}
		xs[i] = m * math.Exp(0.3*r.NormFloat64())
	}
	slices.Sort(xs)
	return xs[int(math.Ceil(0.99*float64(len(xs))))-1]
}
