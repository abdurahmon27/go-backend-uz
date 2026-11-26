package main

import (
	"os"
	"syscall"
)

// ochDirect: faylni O_DIRECT bilan ochadi. Bunda o'qish page cache'dan o'tmaydi, to'g'ridan-to'g'ri
// qurilmaga boradi. Shart: bufer manzili, o'lchami va offset blok o'lchamiga karrali bo'lsin.
// mmap qaytargan xotira sahifa (4096 bayt) chegarasidan boshlanadi.
func ochDirect(nom string, n int) (*os.File, []byte, error) {
	f, err := os.OpenFile(nom, os.O_RDONLY|syscall.O_DIRECT, 0)
	if err != nil {
		return nil, nil, err
	}
	buf, err := syscall.Mmap(-1, 0, n, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	return f, buf, nil
}
