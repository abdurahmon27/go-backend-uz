package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

const workers = 10

func importStocks(ctx context.Context, ids []int64) ([]int, error) {
	stocks := make([]int, len(ids))
	errs := make([]error, len(ids))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for i := range jobs {
				stocks[i], errs[i] = checkStock(ctx, ids[i])
			}
		})
	}
	for i := range ids {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return stocks, errors.Join(errs...)
}

func main() {
	start := time.Now()
	stocks, err := importStocks(context.Background(), catalogIDs(1000))
	report(stocks, err, start)
}
