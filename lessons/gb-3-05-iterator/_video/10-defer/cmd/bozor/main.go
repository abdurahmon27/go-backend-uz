// Iterator ichidagi defer: sikl break qilsa ham, iterator return qilganda fayl yopiladi.
// (Demo: ochish xatosi tashlab yuborilgan; ishda xatoni ham qaytaring.)
package main

import (
	"bufio"
	"fmt"
	"iter"
	"os"
)

// lines faylni qatorma-qator beradi (GB-3.01'dagi bufio.Scanner).
func lines(name string) iter.Seq[string] {
	return func(yield func(string) bool) {
		f, err := os.Open(name)
		if err != nil {
			return
		}
		defer func() {
			f.Close()
			fmt.Println("iterator: fayl yopildi")
		}()
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if !yield(sc.Text()) {
				return
			}
		}
	}
}

func main() {
	for line := range lines("katalog.txt") {
		fmt.Println("sikl:", line)
		break
	}
	fmt.Println("sikldan keyin")
}
