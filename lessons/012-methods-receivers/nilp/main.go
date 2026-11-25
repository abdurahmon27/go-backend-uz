// nilp: nil pointer orqali metod. Pointer receiver nil'ni oʻzi tekshira oladi,
// value receiver esa chaqiruvdan oldin *p ni oʻqiydi va panic beradi (ATAYLAB).
package main

import "fmt"

type Buyurtma struct{ Summa int64 }

func (b Buyurtma) Jami() int64 { return b.Summa }

func (b *Buyurtma) JamiP() int64 {
	if b == nil {
		return 0
	}
	return b.Summa
}

func main() {
	var p *Buyurtma
	fmt.Println("JamiP:", p.JamiP()) // 0: receiver nil, metod buni tekshirdi
	fmt.Println("Jami: ", p.Jami())  // (*p).Jami() → panic
}
