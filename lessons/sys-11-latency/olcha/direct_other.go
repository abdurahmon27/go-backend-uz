//go:build !linux

package main

import "os"

// Linux bo'lmasa O_DIRECT yo'q: oddiy ochamiz (natija page cache'ni ham o'lchashi mumkin).
func ochDirect(nom string, n int) (*os.File, []byte, error) {
	f, err := os.Open(nom)
	return f, make([]byte, n), err
}
