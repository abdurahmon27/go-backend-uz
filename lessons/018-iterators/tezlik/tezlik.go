// Tezlik: bir xil yigʻindi toʻrt usulda.
package tezlik

import (
	"iter"
	"slices"
)

// 1. Oddiy range slice ustida.
func JamiOddiy(ss []int64) int64 {
	var j int64
	for _, s := range ss {
		j += s
	}
	return j
}

// 2. Standart iterator: slices.Values.
func JamiValues(ss []int64) int64 {
	var j int64
	for s := range slices.Values(ss) {
		j += s
	}
	return j
}

// Barchasi — oʻzimiz yozgan iterator.
func Barchasi(ss []int64) iter.Seq[int64] {
	return func(yield func(int64) bool) {
		for _, s := range ss {
			if !yield(s) {
				return
			}
		}
	}
}

// 3. Oʻz iteratorimiz.
func JamiOz(ss []int64) int64 {
	var j int64
	for s := range Barchasi(ss) {
		j += s
	}
	return j
}

// 4. Iterator qiymat sifatida keladi, funksiya inline qilinmaydi:
// endi har bir element uchun yield haqiqiy chaqiruv boʻladi.
//
//go:noinline
func JamiSeq(seq iter.Seq[int64]) int64 {
	var j int64
	for s := range seq {
		j += s
	}
	return j
}
