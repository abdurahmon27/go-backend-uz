package main

import (
	"fmt"
	"math"
)

func main() {
	var f32 float32 = math.Pi
	var f64 float64 = math.Pi
	fmt.Printf("float32: %.20f\n", f32)
	fmt.Printf("float64: %.20f\n", f64)
	fmt.Println("haqiqiy: 3.14159265358979323846...")
}
