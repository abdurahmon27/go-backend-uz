// orta: README dagi 6 qatordan kompilyatsiya boʻladiganlari (1, 2, 5, 6).
// 3 va 4 — ATAYLAB yoʻq: ular kompilyatsiya boʻlmaydi.
// go run ./homework/orta
package main

import "fmt"

type Quti[T any] struct{ v T }

// 1: generic turdagi metod, receiver'dagi T bilan (Go 1.18+)
func (q Quti[T]) Ol() T { return q.v }

// 2: generic turdagi metod, oʻz type parametri bilan (Go 1.27+)
func (q Quti[T]) Juft[U any](u U) (T, U) { return q.v, u }

type Oluvchi[T any] interface{ Ol() T }

func main() {
	q := Quti[int]{7}
	// 5: generic interface'ni oddiy metod bajaradi
	var o Oluvchi[int] = q
	fmt.Println(o.Ol())
	// 6: generic metoddan aniq tur bilan method value
	f := q.Juft[string]
	fmt.Println(f("yetti"))
}
