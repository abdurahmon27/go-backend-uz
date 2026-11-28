package main

import (
	"sync"
	"testing"
	"testing/synctest"
)

// Videodagi "synctest chegarasi": mutex'da kutish durable emas. Goroutine qulfda abadiy qolsa,
// synctest deadlock demaydi: test -timeout tugaguncha osilib turadi.
func TestMutexHang(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var mu sync.Mutex
		mu.Lock()
		go func() {
			mu.Lock() // hech kim Unlock qilmaydi
		}()
	})
}
