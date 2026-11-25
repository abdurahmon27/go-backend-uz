// embed: embedding'da metodlar ham koʻtariladi (GO-1.10 vaʼdasi). Receiver baribir ichki struct.
package main

import "fmt"

type Manzil struct {
	Shahar string
	Kocha  string
}

func (m Manzil) Toliq() string         { return m.Shahar + ", " + m.Kocha }
func (m *Manzil) Kochir(shahar string) { m.Shahar = shahar }

type Mijoz struct {
	Ism    string
	Manzil // embedding
}

func main() {
	a := Mijoz{Ism: "Aziza", Manzil: Manzil{Shahar: "Toshkent", Kocha: "Navoiy"}}
	a.Kochir("Samarqand")  // (&a.Manzil).Kochir("Samarqand")
	fmt.Println(a.Toliq()) // a.Manzil.Toliq()
	fmt.Println(a.Ism, "—", a.Manzil.Shahar)
}
