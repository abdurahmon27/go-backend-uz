package store

// Bu fayl vazifani (TASK.md) tekshiradi: go test ./... uni store.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang. Test store paketining ichida, shuning uchun products maydonini ham koʻradi.
// Bu yerda iterator ikki xil chaqiriladi:
//   for p := range seq { … }  — videodagidek;
//   seq(func(p Product) bool { … }) — iterator oddiy funksiya, uni yield bilan oʻzimiz ham chaqira olamiz:
//   shu yoʻl bilan test yield false qaytargandan keyin iterator yana chaqiradimi, yoʻqmi, sanaydi.
//   slices.Equal(a, b) — ikki slice uzunligi va elementlari bir xilmi (GB-2.01).
//   func(V) bool { … } — parametri ishlatilmaydigan nomsiz funksiya: parametr nomini yozmasa ham boʻladi.
//   max(a, b) — Go'ga oʻrnatilgan funksiya: kattasini qaytaradi.

import (
	"iter"
	"slices"
	"testing"
)

// katalog videodagi toʻrt mahsulotli katalogni qaytaradi.
func katalog() *Store {
	var st Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 3)
	st.Add("choy", 2_500_000, 8)
	st.Add("oʻrik", 3_500_000, 2)
	return &st
}

// names iterator bergan mahsulotlarning nomlarini yigʻadi.
func names(seq iter.Seq[Product]) []string {
	var out []string
	for p := range seq {
		out = append(out, p.Name)
	}
	return out
}

// stopsAfterFalse seq'ni yield bilan oʻzimiz chaqiradi: birinchi qiymatdan keyin false qaytaramiz
// va iterator yield'ni yana necha marta chaqirganini qaytaramiz (toʻgʻri iterator: 0).
func stopsAfterFalse[V any](seq iter.Seq[V]) (calls int) {
	seq(func(V) bool {
		calls++
		return false
	})
	return max(calls-1, 0)
}

// Videodagi All: katalog tartibida beradi va break'da toʻxtaydi (bu test boshidanoq oʻtadi).
func TestVideoAll(t *testing.T) {
	st := katalog()
	want := []string{"olma", "non", "choy", "oʻrik"}
	if got := names(st.All()); !slices.Equal(got, want) {
		t.Fatalf("All() %v berdi, kutilgan %v", got, want)
	}
	if extra := stopsAfterFalse(st.All()); extra != 0 {
		t.Fatalf("All() yield false qaytargandan keyin yana %d marta chaqirdi", extra)
	}
}

// 1-qadam: LowStock(limit) — qoldigʻi limit'dan kam mahsulotlar, katalog tartibida, break'da toʻxtaydi.
func TestVazifaLowStock(t *testing.T) {
	st := katalog()
	tests := []struct {
		name  string
		limit int
		want  []string
	}{
		{"beshdan_kam", 5, []string{"non", "oʻrik"}},
		{"ondan_kam", 10, []string{"non", "choy", "oʻrik"}},
		{"hech_biri", 2, nil},
		{"hammasi", 100, []string{"olma", "non", "choy", "oʻrik"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := names(st.LowStock(tc.limit))
			if !slices.Equal(got, tc.want) {
				t.Fatalf("LowStock(%d) %v berdi, kutilgan %v: qoldigʻi limit'dan kam (<) mahsulotlarni katalog tartibida bering (TASK.md, 1-qadam)", tc.limit, got, tc.want)
			}
		})
	}

	// yield false qaytargach, iterator toʻxtashi kerak (aks holda quyidagi break'li range panic qiladi).
	if extra := stopsAfterFalse(st.LowStock(10)); extra != 0 {
		t.Fatalf("LowStock yield false qaytargandan keyin yana %d marta chaqirdi: if !yield(p) { return } yozing, aks holda range'da panic boʻladi (TASK.md, 1-qadam)", extra)
	}

	// break: birinchi kam qolgan mahsulotda toʻxtaymiz.
	var first string
	for p := range st.LowStock(5) {
		first = p.Name
		break
	}
	if first != "non" {
		t.Fatalf("for range st.LowStock(5) + break: birinchi mahsulot %q, kutilgan \"non\" (TASK.md, 1-qadam)", first)
	}

	// Iterator nusxa emas, jonli koʻrinish: LowStock'ni chaqirgandan keyin qoʻshilgan mahsulot ham
	// range paytida koʻrinadi. Avval slice'ga yigʻib, keyin bergan iterator buni koʻrmaydi.
	seq := st.LowStock(5)
	st.Add("anor", 3_000_000, 1)
	if got, want := names(seq), []string{"non", "oʻrik", "anor"}; !slices.Equal(got, want) {
		t.Fatalf("LowStock(5) %v berdi, kutilgan %v: iterator katalogni range paytida oʻqisin, oldindan nusxa yigʻmasin (TASK.md, 1-qadam)", got, want)
	}
}

// 2-qadam: SortedNames — nomlar slices.Sort tartibida (baytlar boʻyicha), yangi slice'da.
func TestVazifaSortedNames(t *testing.T) {
	st := katalog()
	want := []string{"choy", "non", "olma", "oʻrik"}
	got := st.SortedNames()
	if !slices.Equal(got, want) {
		t.Fatalf("SortedNames() %v qaytardi, kutilgan %v: Names() iteratoridan slices.Sorted bilan yigʻing (TASK.md, 2-qadam)", got, want)
	}
	got[0] = "buzildi"
	if again := st.SortedNames(); !slices.Equal(again, want) {
		t.Fatalf("SortedNames natijasiga yozish katalogni oʻzgartirdi: %v; har safar yangi slice qaytaring (TASK.md, 2-qadam)", again)
	}
	var empty Store
	if got := empty.SortedNames(); len(got) != 0 {
		t.Fatalf("boʻsh katalogda SortedNames() %v, kutilgan boʻsh (TASK.md, 2-qadam)", got)
	}
}
