// olcham: maydonlari bir xil ikki struct, lekin oʻlchami har xil.
package main

import (
	"fmt"
	"unsafe"
)

// Buyurtma: maydonlar "mantiqiy" tartibda yozilgan.
type Buyurtma struct {
	Tolandi    bool  // toʻlandimi?
	Summa      int64 // tiyinda
	Yetkazildi bool  // yetkazildimi?
}

// Ixcham: xuddi shu uchta maydon, faqat tartibi boshqa.
type Ixcham struct {
	Summa      int64
	Tolandi    bool
	Yetkazildi bool
}

func main() {
	fmt.Println("Buyurtma:", unsafe.Sizeof(Buyurtma{}), "bayt")
	fmt.Println("Ixcham:  ", unsafe.Sizeof(Ixcham{}), "bayt")
}
