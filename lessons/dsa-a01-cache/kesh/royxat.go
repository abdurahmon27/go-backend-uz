package kesh

import "math/rand/v2"

// Tugun: bog'langan ro'yxat (linked list) tuguni, 16 bayt.
type Tugun struct {
	Qiymat  int64
	Keyingi *Tugun
}

// Royxat: ro'yxat bo'ylab yurib qiymatlarni qo'shadi.
func Royxat(t *Tugun) int64 {
	var sum int64
	for ; t != nil; t = t.Keyingi {
		sum += t.Qiymat
	}
	return sum
}

// RoyxatYangi: har tugun new bilan, ketma-ket yaratiladi.
// Go allocator'i ularni odatda xotirada yonma-yon qo'yadi.
func RoyxatYangi(n int) *Tugun {
	var bosh, oxir *Tugun
	for i := range n {
		t := &Tugun{Qiymat: int64(i)}
		if bosh == nil {
			bosh = t
		} else {
			oxir.Keyingi = t
		}
		oxir = t
	}
	return bosh
}

// RoyxatSochilgan: o'sha ro'yxat, lekin tugunlar xotirada aralash joylarda turadi
// (uzoq ishlagan dasturda qo'shish/o'chirishlardan keyin shunday bo'ladi).
func RoyxatSochilgan(n int, r *rand.Rand) *Tugun {
	joy := make([]Tugun, n)
	t := Tartib(n, r)
	for k := range n {
		joy[t[k]].Qiymat = int64(k)
		if k+1 < n {
			joy[t[k]].Keyingi = &joy[t[k+1]]
		}
	}
	return &joy[t[0]]
}
