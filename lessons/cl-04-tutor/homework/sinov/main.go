// sinov: Claude ustoz rejimida bergan uchta tekshiruv savoliga javobni ishga tushirib koʻrsatadi.
// Avval savollarga oʻzingiz javob yozing (README → Uy vazifasi), keyin: go run ./homework/sinov
package main

import (
	"errors"
	"fmt"
)

type MaydonXato struct{ Maydon, Sabab string }

func (e *MaydonXato) Error() string { return e.Maydon + ": " + e.Sabab }

// QiymatXato: xuddi shunday, lekin Error() qiymat qabul qiluvchida va pointersiz qaytariladi
type QiymatXato struct{ Maydon, Sabab string }

func (e QiymatXato) Error() string { return e.Maydon + ": " + e.Sabab }

func main() {
	err := fmt.Errorf("tekshir: %w", &MaydonXato{Maydon: "narx", Sabab: "musbat boʻlishi kerak"})

	// 1-savol: errors.Is qiymatni (bu yerda pointer manzilini) solishtiradi, AsType esa turni qidiradi
	_, ok := errors.AsType[*MaydonXato](err)
	fmt.Println("1. errors.Is(err, &MaydonXato{}):", errors.Is(err, &MaydonXato{}), "· AsType[*MaydonXato]:", ok)

	// 2-savol: yana bir qavat %w bilan oʻrasak zanjir saqlanadi, %v bilan esa uziladi
	w := fmt.Errorf("handler: %w", err)
	v := fmt.Errorf("handler: %v", err)
	_, okW := errors.AsType[*MaydonXato](w)
	_, okV := errors.AsType[*MaydonXato](v)
	fmt.Println("2. %w bilan yana oʻralgan:", okW, "· %v bilan:", okV)

	// 3-savol: xato pointersiz qaytsa, *QiymatXato turi zanjirda yoʻq
	q := fmt.Errorf("tekshir: %w", QiymatXato{Maydon: "narx", Sabab: "musbat boʻlishi kerak"})
	_, okP := errors.AsType[*QiymatXato](q)
	_, okQ := errors.AsType[QiymatXato](q)
	fmt.Println("3. pointersiz xato: AsType[*QiymatXato]:", okP, "· AsType[QiymatXato]:", okQ)
}
