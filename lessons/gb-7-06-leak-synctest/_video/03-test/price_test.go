package main

import "testing"

func TestFastestPrice(t *testing.T) {
	got := fastestPrice(1)
	if got != 1_200_000 {
		t.Errorf("fastestPrice(1) = %d; kutilgan 1200000", got)
	}
}
