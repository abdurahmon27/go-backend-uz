package store

import "testing"

// Jamoaning eski testi: Sell qoldiqni kamaytiradi.
func TestSell(t *testing.T) {
	var st Store
	st.Add("non", 400_000, 100)
	if err := st.Sell(1, 3); err != nil {
		t.Fatalf("Sell(1, 3) = %v; kutilgan nil", err)
	}
	if got := st.List()[0].Stock; got != 97 {
		t.Errorf("qoldiq %d; kutilgan 97", got)
	}
}
