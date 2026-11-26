package main

import "testing"

func testParsePrice(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    int64
		wantErr bool
	}{
		{"oddiy", "12000", 1_200_000, false},
		{"chegara", "10000000", 1_000_000_000, false},
		{"harf", "12k", 0, true},
		{"nol", "0", 0, true},
		{"ortiq", "10000001", 0, true},
		{"ulkan", "100000000000000000", 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parsePrice(tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("parsePrice(%q) = %d, %v; xato kutilgan: %v", tc.in, got, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("parsePrice(%q) = %d; kutilgan %d", tc.in, got, tc.want)
			}
		})
	}
}
