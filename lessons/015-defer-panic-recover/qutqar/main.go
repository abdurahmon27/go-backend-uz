// recover: panic'ni ushlab, oddiy xatoga aylantirish (nomlangan natija orqali).
package main

import "fmt"

func narx(narxlar []int, i int) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("narx(%d): %v", i, r)
		}
	}()
	return narxlar[i], nil
}

func main() {
	narxlar := []int{1200, 3400, 560}
	for _, i := range []int{1, 5, 2} {
		n, err := narx(narxlar, i)
		fmt.Println(n, err)
	}
	fmt.Println("dastur davom etdi")
}
