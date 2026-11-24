package tozala

import "strings"

func tozala(s string) string {
	if !strings.Contains(s, "'") { // tez yoʻl: apostrof yoʻq
		return s
	}
	return tozalaEski(s)
}
