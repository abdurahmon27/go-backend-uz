package main

import (
	"testing"
	"testing/synctest"
)

func TestOwnedStore(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := NewOwnedStore(map[int64]int{1: 5})
		if !s.Sell(1, 2) {
			t.Errorf("Sell(1, 2) = false, kutilgan true (qoldiq 5)")
		}
	})
}
