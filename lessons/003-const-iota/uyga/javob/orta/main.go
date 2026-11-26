package main

import "fmt"

func main() {
	const (
		Nomalum    = iota // 0: nomaʼlum kun (nol qiymat)
		Dushanba          // 1
		Seshanba          // 2
		Chorshanba        // 3
		Payshanba         // 4
		Juma              // 5
		Shanba            // 6
		Yakshanba         // 7
	)
	fmt.Println(Dushanba, Juma, Yakshanba)
}
