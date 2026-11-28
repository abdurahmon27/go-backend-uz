package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func checkStock(ctx context.Context, id int64) (int, error) {
	select {
	case n := <-supplier(id):
		return n, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := checkStock(ctx, 3)
	time.Sleep(4 * time.Second)
	fmt.Println(err, "goroutine'lar:", runtime.NumGoroutine())
}
