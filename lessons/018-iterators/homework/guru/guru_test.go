package guru

import (
	"iter"
	"slices"
	"testing"
)

func TestZip(t *testing.T) {
	ismlar := slices.Values([]string{"Aziz", "Bobur", "Dilnoz"})
	summalar := slices.Values([]int64{150000, 99000})
	var got []string
	for ism, s := range Zip(ismlar, summalar) {
		got = append(got, ism)
		_ = s
	}
	if !slices.Equal(got, []string{"Aziz", "Bobur"}) {
		t.Fatalf("got %v", got)
	}
}

// yopildi: iterator'ning defer'i ishladimi?
func manba(yopildi *bool) iter.Seq[int] {
	return func(yield func(int) bool) {
		defer func() { *yopildi = true }()
		for i := 0; ; i++ { // cheksiz
			if !yield(i) {
				return
			}
		}
	}
}

// stop chaqirilmasa, Pull qilingan iterator oʻrtada osilib qoladi va uning defer'i ishlamaydi.
func TestStop(t *testing.T) {
	var a, b bool
	for x, y := range Zip(manba(&a), manba(&b)) {
		if x+y > 4 {
			break
		}
	}
	if !a || !b {
		t.Fatalf("defer ishlamadi: a=%v b=%v", a, b)
	}

	var c bool
	next, _ := iter.Pull(manba(&c)) // stop yoʻq
	next()
	if c {
		t.Fatal("kutilmagan: stop'siz defer ishladi")
	}
}
