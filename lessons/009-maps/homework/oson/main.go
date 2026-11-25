package main

import "fmt"

// narxMatn: "bepul", "yoʻq" yoki narx. Nol narx va yoʻq mahsulot har xil narsa.
func narxMatn(narx map[string]int, nom string) string {
	n, bor := narx[nom]
	switch {
	case !bor:
		return nom + ": bunday mahsulot yoʻq"
	case n == 0:
		return nom + ": bepul"
	default:
		return fmt.Sprintf("%s: %d soʻm", nom, n)
	}
}

func main() {
	narx := map[string]int{"non": 5000, "suv": 0}
	for _, nom := range []string{"non", "suv", "olma"} {
		fmt.Println(narxMatn(narx, nom))
	}
}
