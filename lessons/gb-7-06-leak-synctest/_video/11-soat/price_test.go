package main

import (
	"testing"
	"testing/synctest"
	"time"
)

// Videodagi savol: pufak ichida bir soatlik time.Sleep. Soat soxta: test bir zumda tugaydi.
func TestHour(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		time.Sleep(time.Hour)
		t.Log(time.Since(start), start.UTC())
	})
}
