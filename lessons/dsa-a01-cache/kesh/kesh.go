// Package kesh: bitta O(n) ish — n ta sonni qo'shish — uch xil xotira tartibida.
package kesh

// Ketma: elementlarni xotirada turgan tartibda qo'shadi.
func Ketma(s []int64) int64 {
	var sum int64
	for i := 0; i < len(s); i++ {
		sum += s[i]
	}
	return sum
}

// Aralash: o'sha elementlarni, har birini bir martadan, aralash tartibda qo'shadi.
func Aralash(s []int64, tartib []int32) int64 {
	var sum int64
	for _, i := range tartib {
		sum += s[i]
	}
	return sum
}

// Zanjir: har element keyingi elementning indeksini saqlaydi (pointer o'rniga indeks).
// Keyingi manzil faqat oldingi o'qish tugagach ma'lum bo'ladi.
func Zanjir(s []int64) int64 {
	var sum int64
	i := int64(0)
	for range len(s) {
		sum += s[i]
		i = s[i]
	}
	return sum
}
