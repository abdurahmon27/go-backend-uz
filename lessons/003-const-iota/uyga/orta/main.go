package main

import "fmt"

func main() {
	// VAZIFA: hafta kunlari uchun iota bilan konstantalar yozing.
	// 0 ni "nomaʼlum kun" qoldiring (videodagi Nomalum kabi), keyin Dushanba … Yakshanba.
	const (
		Nomalum = iota
		Dushanba
		// … qolgan kunlarni shu yerga yozing
	)
	fmt.Println(Dushanba) // keyin Juma ni ham chiqaring: Juma necha?
}
