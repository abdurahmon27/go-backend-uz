package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	// Unix vaqti: 1970-01-01 dan beri o'tgan sekundlar. int32'da saqlansa:
	last := time.Unix(math.MaxInt32, 0).UTC()
	fmt.Println(last) // 2038-01-19 03:14:07 +0000 UTC

	var t int32 = math.MaxInt32
	t++
	fmt.Println(time.Unix(int64(t), 0).UTC()) // 1901-12-13 20:45:52 +0000 UTC
}
