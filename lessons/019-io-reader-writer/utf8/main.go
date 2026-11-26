// utf8: chelak baytlarni beradi, harflarni emas — boʻlak chegarasi "ʻ" ni (2 bayt) kesib oʻtadi.
package main

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

func main() {
	r := strings.NewReader("7000 soʻm")
	buf := make([]byte, 4) // 4 baytlik chelak
	for {
		n, err := r.Read(buf)
		if n > 0 {
			fmt.Printf("%q  toʻgʻri UTF-8: %v\n", buf[:n], utf8.Valid(buf[:n]))
		}
		if err == io.EOF {
			break
		}
	}
}
