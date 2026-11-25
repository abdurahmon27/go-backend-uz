package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func main() {
	matn := "non sut non choy sut non"
	sanoq := map[string]int{}
	for _, soz := range strings.Fields(matn) {
		sanoq[soz]++ // yoʻq kalit 0 qaytaradi, shuning uchun ++ darhol ishlaydi
	}
	for _, soz := range slices.Sorted(maps.Keys(sanoq)) {
		fmt.Println(soz, sanoq[soz])
	}
}
