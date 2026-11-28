package main

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

// Videodagi testlar. synctest.Test funksiyani "pufak" ichida ishga tushiradi: soxta soat, va funksiya
// tugaganda pufakda kutib qolgan goroutine boʻlsa, test yiqiladi.

func TestFastestPrice(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		got, err := fastestPrice(t.Context(), 1)
		if err != nil || got != 1_200_000 {
			t.Errorf("fastestPrice(1) = %d, %v; kutilgan 1200000", got, err)
		}
	})
}

func TestFastestPriceTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := fastestPrice(ctx, 1)
		took := time.Since(start) // pufakda soat soxta: aniq tenglik mumkin (haqiqiy soatda emas)
		if !errors.Is(err, context.DeadlineExceeded) || took != 50*time.Millisecond {
			t.Errorf("xato %v, %v oʻtdi; kutilgan deadline, 50ms", err, took)
		}
	})
}

func TestOwnedStore(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := NewOwnedStore(map[int64]int{1: 5})
		defer s.Close()
		if !s.Sell(1, 2) {
			t.Errorf("Sell(1, 2) = false, kutilgan true (qoldiq 5)")
		}
	})
}
