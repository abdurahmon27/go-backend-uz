package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

func importStocks(ctx context.Context, ids []int64) ([]int, error) {
	stocks := make([]int, len(ids))
	errs := make([]error, len(ids))
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Go(func() {
			stocks[i], errs[i] = checkStock(ctx, id)
		})
	}
	wg.Wait()
	return stocks, errors.Join(errs...)
}

func main() {
	start := time.Now()
	stocks, err := importStocks(context.Background(), catalogIDs(1000))
	report(stocks, err, start)
}
