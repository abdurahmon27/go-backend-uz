package qiyin

import "testing"

// Quti: turi comparable, lekin ichidagi X har qanday qiymat boʻlishi mumkin.
type Quti struct{ X any }

func TestPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("panic kutilgandi")
		}
		t.Log("panic:", r)
	}()
	var a, b any = []int{1}, []int{1}
	_ = a == b // kompilyatsiya boʻladi, lekin panic
}

func TestTeng(t *testing.T) {
	jadval := []struct {
		nom  string
		a, b any
		want bool
	}{
		{"ikki nil", nil, nil, true},
		{"int va int", 7, 7, true},
		{"int va int64", 0, int64(0), false},
		{"slice", []int{1}, []int{1}, false},
		{"nil va 0", nil, 0, false},
		{"struct ichida slice", Quti{[]int{1}}, Quti{[]int{1}}, false},
		{"struct ichida int", Quti{7}, Quti{7}, true},
	}
	for _, c := range jadval {
		if got := teng(c.a, c.b); got != c.want {
			t.Errorf("%s: teng = %v, kutilgan %v", c.nom, got, c.want)
		}
	}
}
