package tezlik

import "testing"

func malumot() []int64 {
	ss := make([]int64, 1000)
	for i := range ss {
		ss[i] = int64(i * 1000)
	}
	return ss
}

var natija int64

func BenchmarkOddiy(b *testing.B) {
	ss := malumot()
	for b.Loop() {
		natija = JamiOddiy(ss)
	}
}

func BenchmarkValues(b *testing.B) {
	ss := malumot()
	for b.Loop() {
		natija = JamiValues(ss)
	}
}

func BenchmarkOz(b *testing.B) {
	ss := malumot()
	for b.Loop() {
		natija = JamiOz(ss)
	}
}

func BenchmarkSeq(b *testing.B) {
	ss := malumot()
	for b.Loop() {
		natija = JamiSeq(Barchasi(ss))
	}
}

func TestBirXil(t *testing.T) {
	ss := malumot()
	want := JamiOddiy(ss)
	for nom, got := range map[string]int64{
		"Values": JamiValues(ss),
		"Oz":     JamiOz(ss),
		"Seq":    JamiSeq(Barchasi(ss)),
	} {
		if got != want {
			t.Errorf("%s = %d, kutilgan %d", nom, got, want)
		}
	}
}
