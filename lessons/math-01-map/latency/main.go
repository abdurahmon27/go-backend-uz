// latency: ikki reliz, sintetik model. Qaysi biri tezroq?
// Javob qaysi o'lchovga qarashingizga bog'liq: o'rtacha, mediana yoki p99.
package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
)

const (
	n         = 100_000 // alohida so'rovlar
	sahifalar = 10_000  // har sahifa 100 so'rov yuborib, eng sekinini kutadi
)

// lognormal: mediana atrofida, o'ngga cho'zilgan "dum" bilan (ms).
func lognormal(r *rand.Rand, mediana, sigma float64) float64 {
	return mediana * math.Exp(sigma*r.NormFloat64())
}

// eski: har so'rov bazaga boradi, ~35 ms.
func eski(r *rand.Rand) float64 { return lognormal(r, 35, 0.3) }

// kesh: 96 % so'rov keshdan ~8 ms. 4 % keshda yo'q: sekin yo'l (sovuq disk, og'ir so'rov) ~400 ms.
func kesh(r *rand.Rand) float64 {
	if r.Float64() < 0.96 {
		return lognormal(r, 8, 0.3)
	}
	return lognormal(r, 400, 0.3)
}

// percentil: saralangan xs da "eng yaqin rang" usuli: p ulushdan kam bo'lmagan eng kichik qiymat.
func percentil(xs []float64, p float64) float64 {
	k := int(math.Ceil(p*float64(len(xs)))) - 1
	return xs[max(k, 0)]
}

func olcha(nom string, f func(*rand.Rand) float64) {
	r := rand.New(rand.NewPCG(1, 2))
	xs := make([]float64, n)
	var sum float64
	for i := range xs {
		xs[i] = f(r)
		sum += xs[i]
	}
	slices.Sort(xs)
	sh := make([]float64, sahifalar)
	for i := range sh {
		for range 100 {
			sh[i] = max(sh[i], f(r))
		}
	}
	slices.Sort(sh)
	fmt.Printf("%-4s so'rov:  o'rtacha %4.0f ms | p50 %4.0f ms | p99 %4.0f ms\n",
		nom, sum/n, percentil(xs, 0.50), percentil(xs, 0.99))
	fmt.Printf("%-4s sahifa: 100 so'rov   | p50 %4.0f ms | p99 %4.0f ms\n",
		nom, percentil(sh, 0.50), percentil(sh, 0.99))
}

func main() {
	olcha("eski", eski)
	olcha("kesh", kesh)
}
