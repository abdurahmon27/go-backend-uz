// Uy vazifasi (oʻrta): tur ham Error, ham String metodiga ega. fmt.Println qaysi birini chaqiradi?
// go run ./homework/orta
package main

import "fmt"

type Ikkalasi struct{}

func (Ikkalasi) Error() string  { return "Error() chaqirildi" }
func (Ikkalasi) String() string { return "String() chaqirildi" }

func main() {
	fmt.Println(Ikkalasi{})
}
