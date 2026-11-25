// stringer: Buyurtma hech qayerda "fmt.Stringer'ni bajaraman" demaydi.
// go run ./stringer
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func (b Buyurtma) String() string {
	return fmt.Sprintf("buyurtma: %d soʻm, toʻlandi: %v", b.Summa, b.Tolandi)
}

type Mijoz struct{ Ism string }

func (m Mijoz) String() string { return "mijoz: " + m.Ism }

func main() {
	royxat := []fmt.Stringer{
		Buyurtma{Summa: 150000},
		Mijoz{Ism: "Aziza"},
	}
	for _, s := range royxat {
		fmt.Println(s)
	}
}
