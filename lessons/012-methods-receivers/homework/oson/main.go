// oson (javob): Hisoblagich nega 0 da qoladi? Value receiver nusxani oshiradi.
// Tuzatish: pointer receiver. staticcheck xato variantni SA4005 bilan ushlaydi.
package main

import "fmt"

type Hisoblagich struct{ n int }

// Xato edi: func (h Hisoblagich) Oshir() { h.n++ }
func (h *Hisoblagich) Oshir() { h.n++ }

func (h *Hisoblagich) Qiymat() int { return h.n } // izchillik: hamma metod pointer receiver

func main() {
	var h Hisoblagich
	h.Oshir()
	h.Oshir()
	h.Oshir()
	fmt.Println(h.Qiymat()) // 3
}
