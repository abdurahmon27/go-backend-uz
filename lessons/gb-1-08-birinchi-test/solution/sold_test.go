package main

import "testing"

func TestParseSold(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    int
		wantErr bool
	}{
		{"oddiy", "5", 5, false},
		{"eng_kam", "1", 1, false},
		{"chegara", "1000", 1000, false},
		{"nol", "0", 0, true},
		{"manfiy", "-1", 0, true},
		{"ortiq", "1001", 0, true},
		{"harf", "besh", 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseSold(tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("parseSold(%q) = %d, %v; xato kutilgan: %v", tc.in, got, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("parseSold(%q) = %d; kutilgan %d", tc.in, got, tc.want)
			}
		})
	}
}
