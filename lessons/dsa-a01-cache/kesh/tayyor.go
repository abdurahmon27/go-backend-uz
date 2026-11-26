package kesh

import "math/rand/v2"

// Qadam: har k-elementni qo'shadi (uy vazifasi va pauza savoli).
func Qadam(s []int64, k int) int64 {
	var sum int64
	for i := 0; i < len(s); i += k {
		sum += s[i]
	}
	return sum
}

// Sonlar: 0..n-1 sonlari bilan slice.
func Sonlar(n int) []int64 {
	s := make([]int64, n)
	for i := range s {
		s[i] = int64(i)
	}
	return s
}

// Tartib: 0..n-1 indekslarini tasodifiy aralashtiradi.
func Tartib(n int, r *rand.Rand) []int32 {
	t := make([]int32, n)
	for i := range t {
		t[i] = int32(i)
	}
	r.Shuffle(n, func(i, j int) { t[i], t[j] = t[j], t[i] })
	return t
}

// TartibliZanjir: zanjir, lekin keyingi element doim yonidagisi: s[i] = i+1.
// Har o'qish yana oldingisini kutadi, faqat manzillar ketma-ket.
func TartibliZanjir(n int) []int64 {
	s := make([]int64, n)
	for i := range s {
		s[i] = int64((i + 1) % n)
	}
	return s
}

// ZanjirSlice: 0 dan boshlanib har elementdan bir marta o'tadigan bitta aylana.
// Qiymatlar yana 0..n-1, faqat endi har biri "keyingisi kim" degan indeks.
func ZanjirSlice(n int, r *rand.Rand) []int64 {
	t := Tartib(n, r)
	for k, v := range t { // 0 ni boshiga olib kelamiz
		if v == 0 {
			t[0], t[k] = t[k], t[0]
			break
		}
	}
	s := make([]int64, n)
	for k := range t {
		s[t[k]] = int64(t[(k+1)%n])
	}
	return s
}
