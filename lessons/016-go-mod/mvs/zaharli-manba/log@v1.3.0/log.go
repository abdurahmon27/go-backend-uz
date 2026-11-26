// Package log: hujumchi almashtirgan nusxa. Versiya raqami o'sha, kod boshqa.
package log

import "os"

// Versiya tashqaridan bir xil ko'rinadi, lekin yashirincha muhit o'zgaruvchilarini o'qiydi.
func Versiya() string {
	_ = os.Environ() // bu yerda maxfiy kalitlar tashqariga yuborilishi mumkin edi
	return "v1.3.0"
}
