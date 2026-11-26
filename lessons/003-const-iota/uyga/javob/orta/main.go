// 🟡 Javob: nom ↔ Holat. Bazada "paid" saqlansa, const blokidagi tartib oʻzgarsa ham maʼno buzilmaydi.
package main

import "fmt"

type Holat int

const (
	Nomalum Holat = iota
	Yangi
	Tolangan
	Bekor
)

func HolatNomi(h Holat) string {
	if h == Yangi {
		return "new"
	}
	if h == Tolangan {
		return "paid"
	}
	if h == Bekor {
		return "cancelled"
	}
	return "unknown"
}

func HolatNomdan(s string) (Holat, error) {
	for _, h := range []Holat{Nomalum, Yangi, Tolangan, Bekor} {
		if HolatNomi(h) == s {
			return h, nil
		}
	}
	return Nomalum, fmt.Errorf("nomaʼlum holat: %q", s)
}

func main() {
	h, err := HolatNomdan(HolatNomi(Tolangan))
	fmt.Println(HolatNomi(h), h, err)

	h, err = HolatNomdan("refund")
	fmt.Println(h, err)
}
