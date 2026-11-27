package store

// GB-4.08 vazifasi (Store.Delete) testlari: endi oddiy store testlari.

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"
)

// deleteVaqt bilan: Delete 2 soniyada qaytmasa, test osilib qolmaydi, nima boʻlganini aytadi.
func deleteVaqt(t *testing.T, st *Store, id int64) error {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- st.Delete(id) }()
	select {
	case err := <-done:
		return err
	case <-time.After(2 * time.Second):
		t.Fatalf("Delete(%d) 2 soniyada qaytmadi: qulf ichida Get, Sell yoki All kabi oʻzi qulf oladigan metodni chaqirmang — "+
			"ikkinchi Lock birinchisini kutib qotib qoladi (deadlock, GB-4.07)", id)
		return nil
	}
}

// qulfBosh: Delete'dan keyin qulf boʻsh boʻlishi kerak (defer s.mu.Unlock() unutilmagan).
func qulfBosh(t *testing.T, st *Store) {
	t.Helper()
	if !st.mu.TryLock() {
		t.Fatal("Delete qaytdi, lekin qulf hali band: Lock'dan keyin defer s.mu.Unlock() qoʻying")
	}
	st.mu.Unlock()
}

func TestVazifaDelete(t *testing.T) {
	st := katalog()
	if err := deleteVaqt(t, st, 2); err != nil {
		t.Fatalf("Delete(2) xato qaytardi: %v — non (ID 2) katalogda bor, xato boʻlmasligi kerak", err)
	}
	qulfBosh(t, st)
	if _, err := st.Get(2); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete(2) dan keyin Get(2) xatosi %v, kutilgan ErrNotFound: mahsulot katalogdan oʻchmagan", err)
	}
	var ids []int64
	for p := range st.All() {
		ids = append(ids, p.ID)
	}
	if !slices.Equal(ids, []int64{1, 3, 4}) {
		t.Fatalf("Delete(2) dan keyin katalogdagi ID'lar %v, kutilgan [1 3 4] (qolganlari oʻz tartibida qoladi)", ids)
	}
	if err := deleteVaqt(t, st, 2); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete(2) ikkinchi marta: xato %v, kutilgan ErrNotFound (errors.Is bilan tekshiriladi: %%w bilan oʻrang)", err)
	}
	qulfBosh(t, st)
	if err := deleteVaqt(t, st, 99); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete(99): xato %v, kutilgan ErrNotFound", err)
	}
}

// Delete products'ni oʻzgartiradi, demak qulfni olishi shart: qulf band boʻlsa, u kutishi va
// katalogga tegmasligi kerak.
func TestVazifaDeleteQulf(t *testing.T) {
	st := katalog()
	st.mu.Lock() // qulfni test ushlab turadi
	done := make(chan error, 1)
	go func() { done <- st.Delete(1) }()
	select {
	case err := <-done:
		st.mu.Unlock()
		t.Fatalf("Delete qulf band boʻlsa ham darhol qaytdi (xato: %v): boshiga s.mu.Lock() va defer s.mu.Unlock() qoʻying (GB-4.07)", err)
	case <-time.After(200 * time.Millisecond):
		// toʻgʻri: Delete qulfni kutyapti
	}
	if n := len(st.products); n != 4 {
		st.mu.Unlock()
		t.Fatalf("qulf testda boʻlganda Delete katalogni oʻzgartirdi (%d ta mahsulot qoldi): Lock metodning eng boshida boʻlsin", n)
	}
	st.mu.Unlock()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("qulf boʻshagandan keyin ham Delete 2 soniyada qaytmadi (deadlock?)")
	}
	qulfBosh(t, st)
}

// Bir vaqtda Delete, Add va All: katalog buzilmasin. go test -race bilan ishga tushirsangiz,
// qulfsiz oʻqish/yozishni race detector ham koʻrsatadi (TASK.md).
func TestVazifaDeleteParallel(t *testing.T) {
	st := &Store{}
	const n = 400
	for i := range n {
		if _, err := st.Add(fmt.Sprintf("m%d", i+1), 100, 1); err != nil {
			t.Fatal(err)
		}
	}
	var wg sync.WaitGroup
	for id := int64(1); id <= n; id += 2 { // toq ID'lar oʻchadi
		wg.Go(func() { st.Delete(id) })
		wg.Go(func() {
			for range st.All() {
			}
		})
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("bir vaqtdagi Delete va All 5 soniyada tugamadi (deadlock?)")
	}
	var ids []int64
	for p := range st.All() {
		ids = append(ids, p.ID)
	}
	want := []int64{}
	for id := int64(2); id <= n; id += 2 {
		want = append(want, id)
	}
	if !slices.Equal(ids, want) {
		t.Fatalf("%d ta parallel Delete'dan keyin katalogda %d ta mahsulot, kutilgan %d ta (faqat juft ID'lar): "+
			"Delete'ning hamma qismi — qidirish ham, oʻchirish ham — bitta Lock ostida boʻlsin", n/2, len(ids), len(want))
	}
}
