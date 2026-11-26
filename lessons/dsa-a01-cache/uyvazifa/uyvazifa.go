// Package uyvazifa: DSA-A.01 uy vazifasi javoblari (o'lchov: uyvazifa_test.go, natija: README.md).
package uyvazifa

// Qadam: har k-elementni qo'shadi.
func Qadam(s []int64, k int) int64 {
	var sum int64
	for i := 0; i < len(s); i += k {
		sum += s[i]
	}
	return sum
}

// Yig64 va Yig32: bir xil ish, element 8 yoki 4 bayt.
func Yig64(s []int64) int64 {
	var sum int64
	for i := 0; i < len(s); i++ {
		sum += s[i]
	}
	return sum
}

func Yig32(s []int32) int64 {
	var sum int64
	for i := 0; i < len(s); i++ {
		sum += int64(s[i])
	}
	return sum
}

// Nuqta: 32 bayt. Bitta kesh qatoriga ikkitasi sig'adi.
type Nuqta struct{ X, Y, Z, W int64 }

// YigX: struct'lar slice'idan (AoS) faqat X ni yig'adi.
func YigX(s []Nuqta) int64 {
	var sum int64
	for i := 0; i < len(s); i++ {
		sum += s[i].X
	}
	return sum
}
