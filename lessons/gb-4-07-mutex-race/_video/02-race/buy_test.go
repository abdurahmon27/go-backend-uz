package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// orderParallel n ta "POST /buy" soʻrovini bir vaqtda yuboradi.
func orderParallel(st *Store, n int) {
	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			st.buy(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/buy", nil))
		})
	}
	wg.Wait()
}

func TestBuyParallel(t *testing.T) {
	st := &Store{products: map[int64]Product{3: {Stock: 1000}}}
	orderParallel(st, 1000)
	if got := st.products[3].Stock; got != 0 {
		t.Errorf("qoldiq %d, kutilgan 0", got)
	}
}
