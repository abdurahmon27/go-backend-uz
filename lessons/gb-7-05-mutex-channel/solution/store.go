// bozor: bitta umumiy holat — mahsulot qoldigʻi (ID → dona) — ikki usulda himoyalangan (GB-7.05).
// LockedStore: map + sync.Mutex (GB-4.07). OwnedStore: map'ning yagona egasi boʻlgan goroutine, qolganlar
// unga channel orqali soʻrov yuboradi. Ikkalasida Sell bir xil ishlaydi.
package main

import (
	"sync"
	"sync/atomic"
)

// LockedStore — qoldiq map'i va uni himoyalaydigan qulf. stock'ga har murojaat mu ostida.
type LockedStore struct {
	mu    sync.Mutex
	stock map[int64]int
}

// Sell id'li mahsulotdan n dona sotadi; qoldiq yetmasa false.
func (s *LockedStore) Sell(id int64, n int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stock[id] < n {
		return false
	}
	s.stock[id] -= n
	return true
}

// sellReq — egasiga yuboriladigan sotish soʻrovi: qaysi mahsulot, necha dona va javob qaytadigan channel.
type sellReq struct {
	id    int64
	n     int
	reply chan bool
}

// OwnedStore — qoldiq map'iga faqat loop goroutine'i tegadi; Sell unga soʻrov yuboradi va javobni kutadi.
type OwnedStore struct {
	sells  chan sellReq
	stocks chan stockReq
	calls  atomic.Int64 // Sell chaqiruvlari soni
}

// stockReq — qoldiqni soʻrash: qaysi mahsulot va javob channel'i.
type stockReq struct {
	id    int64
	reply chan int
}

// NewOwnedStore loop'ni yangi goroutine'da boshlaydi va stock map'ini unga topshiradi:
// bundan keyin chaqiruvchi bu map'ga tegmasligi kerak.
func NewOwnedStore(stock map[int64]int) *OwnedStore {
	s := &OwnedStore{sells: make(chan sellReq), stocks: make(chan stockReq)}
	go s.loop(stock)
	return s
}

// loop — stock'ning yagona egasi. Ikki xil soʻrovni select bilan kutadi va birma-bir bajaradi;
// sells yopilganda tugaydi.
func (s *OwnedStore) loop(stock map[int64]int) {
	for {
		select {
		case r, ok := <-s.sells:
			if !ok {
				return
			}
			sold := stock[r.id] >= r.n
			if sold {
				stock[r.id] -= r.n
			}
			r.reply <- sold
		case r := <-s.stocks:
			r.reply <- stock[r.id]
		}
	}
}

// Sell id'li mahsulotdan n dona sotadi; qoldiq yetmasa false.
func (s *OwnedStore) Sell(id int64, n int) bool {
	s.calls.Add(1)
	reply := make(chan bool)
	s.sells <- sellReq{id, n, reply}
	return <-reply
}

// Close loop goroutine'ini toʻxtatadi. Close'dan keyin Sell chaqirilmaydi; Close ikkinchi marta chaqirilsa,
// yopiq channel'ni qayta yopish panic beradi (GB-7.02).
func (s *OwnedStore) Close() { close(s.sells) }

// Stock id'li mahsulotning hozirgi qoldigʻini beradi (yoʻq mahsulot: 0).
func (s *OwnedStore) Stock(id int64) int {
	reply := make(chan int)
	s.stocks <- stockReq{id, reply}
	return <-reply
}

// Requests shu paytgacha nechta Sell chaqirilganini beradi (qoldiq yetmagan chaqiruvlar ham).
func (s *OwnedStore) Requests() int64 {
	return s.calls.Load()
}
