// stringerxato: kompilyator method set'ni tekshiradi. Bu papka ATAYLAB kompilyatsiya boʻlmaydi:
// go build -o /dev/null ./stringerxato
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func (b *Buyurtma) String() string { return fmt.Sprint(b.Summa) }

func main() {
	b := Buyurtma{Summa: 50_000}
	var s fmt.Stringer = b // ❌ String pointer receiver'da
	fmt.Println(s)
}
