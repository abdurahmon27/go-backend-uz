package main

import "testing"

func TestParsePrice(t *testing.T) {
	got, err := parsePrice("12000")
	if err != nil || got != 1_200_000 {
		t.Errorf("parsePrice(%q) = %d, %v; kutilgan %d, xatosiz", "12000", got, err, 1_200_000)
	}
	got, err = parsePrice("500")
	if err != nil || got != 50_000 {
		t.Errorf("parsePrice(%q) = %d, %v; kutilgan %d, xatosiz", "500", got, err, 50_000)
	}
}
