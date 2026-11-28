package main

// OwnedStore — GB-7.05'dagi qoldiq doʻkoni (qisqartirilgan): stock map'iga faqat loop goroutine'i tegadi,
// Sell unga channel orqali soʻrov yuboradi. loop sells yopilganda (Close) tugaydi.
type OwnedStore struct {
	sells chan sellReq
}

// sellReq — egasiga yuboriladigan sotish soʻrovi: qaysi mahsulot, necha dona va javob qaytadigan channel.
type sellReq struct {
	id    int64
	n     int
	reply chan bool
}

// NewOwnedStore loop'ni yangi goroutine'da boshlaydi va stock map'ini unga topshiradi.
func NewOwnedStore(stock map[int64]int) *OwnedStore {
	s := &OwnedStore{sells: make(chan sellReq)}
	go s.loop(stock)
	return s
}

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

// Close loop goroutine'ini toʻxtatadi. Close'dan keyin Sell chaqirilmaydi.
func (s *OwnedStore) Close() { close(s.sells) }
