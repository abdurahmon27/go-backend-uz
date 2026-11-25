// guru (javob): embedding va method set. Kochiruvchi interface'ini kim qanoatlantiradi?
//
//	type A struct{ Manzil }   — A ning method set'ida Toliq bor, Kochir yoʻq; *A da ikkalasi bor.
//	type B struct{ *Manzil }  — B va *B ning method set'ida ikkalasi ham bor.
//
// var _ Kochiruvchi = A{} kompilyatsiya boʻlmaydi:
//
//	cannot use A{} (value of struct type A) as Kochiruvchi value in variable declaration:
//	A does not implement Kochiruvchi (method Kochir has pointer receiver)
package main

import "fmt"

type Manzil struct{ Shahar string }

func (m Manzil) Toliq() string         { return m.Shahar }
func (m *Manzil) Kochir(shahar string) { m.Shahar = shahar }

type Kochiruvchi interface{ Kochir(string) }

type A struct{ Manzil }
type B struct{ *Manzil }

var (
	_ Kochiruvchi = &A{} // ✅
	_ Kochiruvchi = B{}  // ✅ (lekin B{}.Manzil == nil: Kochir chaqirilsa panic)
	_ Kochiruvchi = &B{} // ✅
	// _ Kochiruvchi = A{} // ❌ yuqoridagi xato
)

func main() {
	b := B{Manzil: &Manzil{Shahar: "Toshkent"}}
	var k Kochiruvchi = b // B nusxalanadi, lekin ichidagi pointer bitta Manzil'ga qaraydi
	k.Kochir("Buxoro")
	fmt.Println(b.Toliq()) // Buxoro
}
