package main

import "fmt"

type Perm uint8

const (
	Read  Perm = 1 << iota // 001
	Write                  // 010
	Exec                   // 100
)

func main() {
	p := Read | Write
	fmt.Printf("%03b\n", p)        // 011
	fmt.Println(p&Write != 0)      // true: yozish huquqi bor
	fmt.Println(p&Exec != 0)       // false: ishga tushirish yo'q
	fmt.Printf("%03b\n", p&^Write) // 001: Write olib tashlandi
}
