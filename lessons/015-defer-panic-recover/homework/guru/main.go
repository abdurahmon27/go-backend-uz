// Guru: to'rtta variantdan qaysi birida recover panic'ni ushlaydi?
package main

import "fmt"

func qutqar() {
	if r := recover(); r != nil {
		fmt.Println("  ushlandi:", r)
	}
}

func a() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("  ushlandi:", r)
		}
	}()
	panic("a")
}

func b() {
	defer qutqar()
	panic("b")
}

func c() {
	defer func() { qutqar() }()
	panic("c")
}

func d() {
	defer recover()
	panic("d")
}

func sina(nom string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("  ushlanmadi, tashqariga chiqdi:", r)
		}
	}()
	fmt.Println(nom + ":")
	f()
}

func main() {
	sina("a", a)
	sina("b", b)
	sina("c", c)
	sina("d", d)
}
