package main

// Vazifani tekshiradi (TASK.md). Oʻzgartirmang.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

const childEnv = "BOZOR_TEKSHIRUV_CHILD"

func TestVazifaRestockOddiy(t *testing.T) {
	st := &Store{products: map[int64]Product{3: {ID: 3, Stock: 5}}}
	if err := st.Restock(3, 10); err != nil {
		t.Fatalf("Restock(3, 10) xato qaytardi: %v", err)
	}
	if got := st.Stock(3); got != 15 {
		t.Errorf("Restock(3, 10)dan keyin qoldiq %d, kutilgan 15 (5 + 10)", got)
	}
	if err := st.Restock(99, 1); !errors.Is(err, ErrNotFound) {
		t.Errorf("yoʻq mahsulot: Restock(99, 1) xatosi %v, kutilgan ErrNotFound", err)
	}
}

func TestVazifaTotalOddiy(t *testing.T) {
	st := &Store{products: map[int64]Product{1: {Stock: 40}, 2: {Stock: 120}, 3: {Stock: 5}}}
	if got := st.Total(); got != 165 {
		t.Errorf("Total() = %d, kutilgan 165 (40 + 120 + 5): Total hamma mahsulotlarning qoldigʻini qoʻshishi kerak", got)
	}
}

// TestVazifaParallel Sell, Restock va Total'ni bir vaqtda chaqiradi. Qulfsiz kod ba'zan oʻzi qulaydi
// (fatal error: concurrent map …: uni test ichida ushlab boʻlmaydi), shuning uchun tekshiruv ikkinchi
// jarayonda ishlaydi va uning chiqishini oʻzbekcha xabarga aylantiradi.
func TestVazifaParallel(t *testing.T) {
	if os.Getenv(childEnv) == "" {
		for run := 1; run <= 3; run++ {
			cmd := exec.Command(os.Args[0], "-test.run=^TestVazifaParallel$", "-test.count=1", "-test.timeout=30s")
			cmd.Env = append(os.Environ(), childEnv+"=1")
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("bir vaqtda 900 ta chaqiruvda (Sell + Restock + Total) store buzildi yoki qotib qoldi (%d-urinish).\n"+
					"Tekshiring: Restock va Total ham s.mu.Lock() + defer s.mu.Unlock() ostida boʻlsin; "+
					"Total ichida s.Stock'ni chaqirmang (qulf ichida yana Lock: deadlock).\n"+
					"Tekshiruv jarayonining chiqishi (boshi):\n%s", run, head(string(out), 8))
			}
		}
		return
	}
	parallelChecks(t)
}

func parallelChecks(t *testing.T) {
	st := &Store{products: map[int64]Product{1: {ID: 1, Stock: 100}, 3: {ID: 3, Stock: 1000}}}
	var wg sync.WaitGroup
	for range 300 {
		wg.Go(func() { st.Sell(3, 1) })
		wg.Go(func() { st.Restock(3, 1) })
		wg.Go(func() { st.Total() })
	}
	wg.Wait()
	if got := st.Stock(3); got != 1000 {
		t.Fatalf("300 Sell va 300 Restock dan keyin qoldiq %d, kutilgan 1000: yangilanishlar yoʻqoldi", got)
	}
	if got := st.Total(); got != 1100 {
		t.Fatalf("Total() = %d, kutilgan 1100 (100 + 1000)", got)
	}
}

func head(s string, n int) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) > n {
		lines = append(lines[:n], fmt.Sprintf("… (yana %d qator)", len(lines)-n))
	}
	return strings.Join(lines, "\n")
}
