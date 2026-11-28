package main

// Vazifani tekshiradi (TASK.md). Oʻzgartirmang.
//
// Har chaqiruv vaqt chegarasi bilan ishlaydi: loop soʻrovni olmasa yoki javob bermasa, Stock abadiy kutadi,
// shuning uchun test 2 soniyadan keyin oʻzbekcha xabar bilan yiqiladi (select va time.After: GB-7.03).
// reflect (OwnedStore maydonlarini koʻrish) va runtime.NumGoroutine (goroutine'lar soni) — tekshiruvning ichki
// qismi, ularni oʻqish shart emas.

import (
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
)

// within f'ni 2 soniya kutadi; ulgurmasa false.
func within(f func()) bool {
	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(2 * time.Second):
		return false
	}
}

const qotdi = "2 soniyada javob kelmadi: Stock soʻrov yuborib javob kutyapti, lekin loop uni olmayapti yoki javob bermayapti.\n" +
	"\tNewOwnedStore yangi channel'ni make bilan yaratdimi? loop select'da shu channel'ni ham kutadimi va reply'ga yozadimi?"

func TestVazifaStock(t *testing.T) {
	s := NewOwnedStore(map[int64]int{1: 40, 3: 5})
	defer s.Close()

	var got1, got3, got99 int
	if !within(func() { got1, got3, got99 = s.Stock(1), s.Stock(3), s.Stock(99) }) {
		t.Fatal(qotdi)
	}
	if got1 == -1 {
		t.Fatalf("Stock(1) = -1: Stock hali yozilmagan. TASK.md 1-qadam: Stock soʻrovini loop'ga yuborib, javobni qaytaring")
	}
	if got1 != 40 || got3 != 5 || got99 != 0 {
		t.Fatalf("Stock(1), Stock(3), Stock(99) = %d, %d, %d; kutilgan 40, 5, 0 (yoʻq mahsulot: 0)", got1, got3, got99)
	}

	// 8 goroutine oʻrikdan (ID 3, 5 dona) ikki martadan sotadi, yana 8 tasi shu paytda qoldiqni soʻraydi.
	var wg sync.WaitGroup
	var mu sync.Mutex
	sold, bad := 0, 0
	ok := within(func() {
		for range 8 {
			wg.Go(func() {
				for range 2 {
					if s.Sell(3, 1) {
						mu.Lock()
						sold++
						mu.Unlock()
					}
				}
			})
			wg.Go(func() {
				for range 50 {
					if n := s.Stock(3); n < 0 || n > 5 {
						mu.Lock()
						bad = n
						mu.Unlock()
					}
				}
			})
		}
		wg.Wait()
	})
	if !ok {
		t.Fatal(qotdi + "\n\tSell va Stock bir vaqtda chaqirilganda qotdi.")
	}
	if bad != 0 {
		t.Fatalf("Stock(3) %d qaytardi: qoldiq 0 bilan 5 orasida boʻlishi kerak", bad)
	}
	var left int
	if !within(func() { left = s.Stock(3) }) {
		t.Fatal(qotdi)
	}
	if sold != 5 || left != 0 {
		t.Fatalf("16 ta Sell(3, 1): %d tasi oʻtdi, Stock(3) = %d; kutilgan 5 va 0", sold, left)
	}

	// Holat faqat loop'da yashaydi: OwnedStore maydonlarida (koʻrsatkich orqali ham) map ham, qulf ham boʻlmasligi kerak.
	mutexT, rwT := reflect.TypeOf(sync.Mutex{}), reflect.TypeOf(sync.RWMutex{})
	lockerT := reflect.TypeOf((*sync.Locker)(nil)).Elem()
	tp := reflect.TypeOf(OwnedStore{})
	for i := range tp.NumField() {
		f := tp.Field(i)
		ft := f.Type
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}
		switch {
		case ft.Kind() == reflect.Map:
			t.Errorf("OwnedStore'da map maydoni bor (%s): map faqat loop goroutine'ida yashashi kerak, struct'da emas", f.Name)
		case ft == mutexT || ft == rwT || ft == lockerT:
			t.Errorf("OwnedStore'da qulf bor (%s): bu vazifa qulfsiz — holatni loop goroutine'i himoyalaydi", f.Name)
		}
	}
}

func TestVazifaClose(t *testing.T) {
	before := runtime.NumGoroutine()
	stores := make([]*OwnedStore, 20)
	ok := within(func() {
		for i := range stores {
			stores[i] = NewOwnedStore(map[int64]int{1: 1})
			stores[i].Sell(1, 1)
			stores[i].Stock(1)
		}
	})
	if !ok {
		t.Fatal(qotdi)
	}
	for _, s := range stores {
		s.Close()
	}
	deadline := time.Now().Add(2 * time.Second)
	for runtime.NumGoroutine() > before {
		if time.Now().After(deadline) {
			t.Fatalf("Close'dan keyin loop goroutine'lari tugamadi (goroutine'lar: oldin %d, hozir %d).\n"+
				"\tYopiq channel'dan oʻqish darhol nol qiymat beradi (GB-7.02): select'da `r, ok := <-s.sells` va `if !ok { return }`",
				before, runtime.NumGoroutine())
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestVazifaRequests(t *testing.T) {
	s := NewOwnedStore(map[int64]int{1: 3})
	defer s.Close()
	var wg sync.WaitGroup
	ok := within(func() {
		for range 100 {
			wg.Go(func() { s.Sell(1, 1) })
		}
		wg.Wait()
	})
	if !ok {
		t.Fatal(qotdi)
	}
	n := s.Requests()
	if n == 0 {
		t.Skip("🟡 Requests hali yozilmagan (ixtiyoriy, TASK.md 2-qadam)")
	}
	if n != 100 {
		t.Fatalf("100 ta Sell (3 tasi oʻtdi, 97 tasi qoldiq yetmadi), Requests() = %d; kutilgan 100 (har chaqiruv sanaladi)", n)
	}
}
