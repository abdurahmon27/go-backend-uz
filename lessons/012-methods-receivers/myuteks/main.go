// myuteks: sync.Mutex bor struct'da value receiver — ATAYLAB xato. go vet ./myuteks buni ushlaydi.
// Qosh (pointer receiver) lock'ni ushlab turib Balansi (value receiver) ni chaqiradi:
// Balansi mutex'ning nusxasini oladi — nusxa ham qulflangan, Lock abadiy kutadi.
package main

import (
	"fmt"
	"sync"
)

type Hisob struct {
	mu     sync.Mutex
	balans int64
}

func (h *Hisob) Qosh(n int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.balans += n
	fmt.Println("yangi balans:", h.Balansi())
}

// Balansi: value receiver — mutex ham nusxalanadi.
func (h Hisob) Balansi() int64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.balans
}

func main() {
	var h Hisob
	h.Qosh(100)
}
