package main

import (
	"context"
	"fmt"
	"time"
)

func importStocks(ctx context.Context, ids []int64) ([]int, error) {
	stocks := make([]int, len(ids))
	for i, id := range ids {
		n, err := checkStock(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("mahsulot %d: %w", id, err)
		}
		stocks[i] = n
	}
	return stocks, nil
}

func main() {
	start := time.Now()
	stocks, err := importStocks(context.Background(), catalogIDs(1000))
	report(stocks, err, start)
}
