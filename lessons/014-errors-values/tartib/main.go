// tartib: natija err tekshirilmasdan OLDIN ishlatiladi. ATAYLAB XATO.
// Go 1.21.5–1.24 kompilyatori nil tekshiruvini kechiktirib yuborardi (golang/go#72860;
// bug CL 537775 bilan 1.21.5 ga backport qilingan, #63743),
// shuning uchun bu dastur u versiyalarda jim ishlaydi. Go 1.21.4 gacha va 1.25+ da panic.
// ./versiyalar.sh
package main

import (
	"errors"
	"fmt"
)

type Buyurtma struct {
	ID    int
	Summa int64
}

var ErrTopilmadi = errors.New("buyurtma topilmadi")

func top(id int) (*Buyurtma, error) {
	if id != 7 {
		return nil, ErrTopilmadi
	}
	return &Buyurtma{ID: 7, Summa: 150000}, nil
}

func main() {
	b, err := top(42)
	summa := b.Summa // XATO: err hali tekshirilmagan, b esa nil
	if err != nil {
		fmt.Println("xato:", err)
		return
	}
	fmt.Println("summa:", summa)
}
