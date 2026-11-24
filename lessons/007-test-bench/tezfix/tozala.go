package tozala

import "strings"

func tozala(s string) string {
	if !strings.ContainsAny(s, "'`´‘’ʻʼ") { // hech bir turi yoʻq
		return s
	}
	return tozalaEski(s)
}
