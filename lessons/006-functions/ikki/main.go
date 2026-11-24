package main

import (
	"errors"
	"fmt"
	"strconv"
)

func bolish(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("nolga boʻlib boʻlmaydi")
	}
	return a / b, nil
}

func main() {
	q, err := bolish(10, 3)
	fmt.Println(q, err)
	_, err = bolish(1, 0) // birinchi qiymat kerak emas: _
	fmt.Println(err)
	n, err := strconv.Atoi("12a")
	fmt.Println(n, err)
}
