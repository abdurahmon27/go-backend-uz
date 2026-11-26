package main

// Videodagi tuzatishni tekshiradi: afterOrder qoldiqni toʻgʻri hisoblaydi va ombordagi stock'ni oʻzgartirmaydi.

import (
	"slices"
	"testing"
)

func TestAfterOrder(t *testing.T) {
	stock := []int{40, 12, 7}
	sold := []int{2, 1, 0}
	left := afterOrder(stock, sold)
	// slices.Equal ikki slice'ning uzunligi va elementlari bir xilligini tekshiradi
	// (slice'larni == bilan solishtirib boʻlmaydi: kompilyator buni rad etadi).
	if want := []int{38, 11, 7}; !slices.Equal(left, want) {
		t.Errorf("afterOrder(%v, %v) = %v; kutilgan %v", []int{40, 12, 7}, sold, left, want)
	}
	if want := []int{40, 12, 7}; !slices.Equal(stock, want) {
		t.Errorf("afterOrder stock'ni oʻzgartirdi: endi %v, %v boʻlishi kerak", stock, want)
	}
}
