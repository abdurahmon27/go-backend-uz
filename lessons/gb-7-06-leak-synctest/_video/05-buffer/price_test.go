package main

import (
	"testing"
	"testing/synctest"
)

func TestFastestPrice(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		got := fastestPrice(1)
		if got != 1_200_000 {
			t.Errorf("fastestPrice(1) = %d; kutilgan 1200000", got)
		}
	})
}
