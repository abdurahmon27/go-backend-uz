package main

import (
	"context"
	"fmt"
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
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()
	n, err := checkStock(ctx, 3)
	fmt.Println(n, err, time.Since(start).Round(time.Second))
}
