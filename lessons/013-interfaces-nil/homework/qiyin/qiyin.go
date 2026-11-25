// Uy vazifasi (qiyin): ikki interface'ni panic'siz solishtirish.
// a == b kompilyatsiya boʻladi, lekin ichidagi tur solishtirib boʻlmaydigan boʻlsa
// (slice, map, funksiya — yoki ularni ichida saqlagan struct) — runtime panic.
// go test -v ./homework/qiyin
package qiyin

// teng: a == b, lekin panic oʻrniga false.
// reflect.TypeOf(a).Comparable() yetarli emas: struct{ X any } "comparable" turi,
// ichidagi X esa slice boʻlishi mumkin. Shuning uchun panic'ni ushlaymiz.
func teng(a, b any) (natija bool) {
	defer func() {
		if recover() != nil {
			natija = false
		}
	}()
	return a == b
}
