package oldindan

// Ikkala funksiya natijani chaqiruvchiga qaytaradi, shuning uchun slice heap'da yashaydi.
// go:noinline kompilyatorga ularni benchmark ichiga koʻchirishni taqiqlaydi:
// aks holda u slice'ni stekka joylashi mumkin va biz haqiqiy dasturdagi holatni oʻlchamaymiz.

// yigNil: boʻsh slice'dan boshlaydi, append oʻzi oʻstiradi.
//
//go:noinline
func yigNil(n int) []int {
	var s []int
	for i := range n {
		s = append(s, i)
	}
	return s
}

// yigMake: kerakli sigʻimni oldindan soʻraydi.
//
//go:noinline
func yigMake(n int) []int {
	s := make([]int, 0, n)
	for i := range n {
		s = append(s, i)
	}
	return s
}
