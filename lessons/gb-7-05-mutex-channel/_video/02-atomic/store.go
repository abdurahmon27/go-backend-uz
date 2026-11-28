// bozor: bitta umumiy holat — mahsulot qoldigʻi (ID → dona) — ikki usulda himoyalangan (GB-7.05).
// LockedStore: map + sync.Mutex (GB-4.07). OwnedStore: map'ning yagona egasi boʻlgan goroutine, qolganlar
// unga channel orqali soʻrov yuboradi. Ikkalasida Sell bir xil ishlaydi.
package main

import "sync"

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
	sells chan sellReq
}

// NewOwnedStore loop'ni yangi goroutine'da boshlaydi va stock map'ini unga topshiradi:
// bundan keyin chaqiruvchi bu map'ga tegmasligi kerak.
func NewOwnedStore(stock map[int64]int) *OwnedStore {
	s := &OwnedStore{sells: make(chan sellReq)}
	go s.loop(stock)
	return s
}

// loop — stock'ning yagona egasi. Soʻrovlarni birma-bir bajaradi; sells yopilganda tugaydi.
func (s *OwnedStore) loop(stock map[int64]int) {
	for r := range s.sells {
		ok := stock[r.id] >= r.n
		if ok {
			stock[r.id] -= r.n
		}
		r.reply <- ok
	}
}

// Sell id'li mahsulotdan n dona sotadi; qoldiq yetmasa false.
func (s *OwnedStore) Sell(id int64, n int) bool {
	reply := make(chan bool)
	s.sells <- sellReq{id, n, reply}
	return <-reply
}

// Close loop goroutine'ini toʻxtatadi. Close'dan keyin Sell chaqirilmaydi; Close ikkinchi marta chaqirilsa,
// yopiq channel'ni qayta yopish panic beradi (GB-7.02).
func (s *OwnedStore) Close() { close(s.sells) }
