// oson: butun sonning manzilini qaytaradigan funksiya.
// go build -gcflags=-m ./homework/oson  ->  moved to heap: x
// x funksiya tugagach ham pointer orqali ishlatiladi, shuning uchun heap'da.
package oson

func yangiSon() *int {
	x := 42
	return &x
}
