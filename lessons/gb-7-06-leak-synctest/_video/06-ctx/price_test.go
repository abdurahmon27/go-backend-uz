package main

import (
	"testing"
	"testing/synctest"
)

func TestFastestPrice(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		got, err := fastestPrice(t.Context(), 1)
		if err != nil || got != 1_200_000 {
			t.Errorf("fastestPrice(1) = %d, %v; kutilgan 1200000", got, err)
		}
	})
}
