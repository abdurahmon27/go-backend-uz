package main

import "time"

// stocks: yetkazib beruvchidagi qoldiqlar (ID → dona): olma 1, non 2, oʻrik 3.
var stocks = map[int64]int{1: 40, 2: 120, 3: 5}

// supplierDelay: yetkazib beruvchi shuncha vaqtda javob beradi (tekshiruv testi uni qisqartiradi).
var supplierDelay = 3 * time.Second

// supplier qoldiqni yetkazib beruvchidan soʻraydi: javob supplierDelay'dan keyin channel'ga keladi.
// Buffer 1: hech kim qabul qilmasa ham goroutine yuborib, tugaydi (videoning "Chuqurroq" qismi).
func supplier(id int64) chan int {
	ch := make(chan int, 1)
	go func() {
		time.Sleep(supplierDelay)
		ch <- stocks[id]
	}()
	return ch
}
