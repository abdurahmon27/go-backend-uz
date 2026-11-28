package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

const workers = 10

func importStocks(ctx context.Context, ids []int64) ([]int, error) {
	stocks := make([]int, len(ids))
	var g errgroup.Group
	g.SetLimit(workers)
	for i, id := range ids {
		g.Go(func() error {
			n, err := checkStock(ctx, id)
			if err != nil {
				return fmt.Errorf("mahsulot %d: %w", id, err)
			}
			stocks[i] = n
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return stocks, nil
}

func main() {
	start := time.Now()
	ids := catalogIDs(1000)
	ids[499] = 9999
	stocks, err := importStocks(context.Background(), ids)
	report(stocks, err, start)
}
