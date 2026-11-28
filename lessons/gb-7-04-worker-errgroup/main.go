// bozor: yetkazib beruvchidan 1000 mahsulot qoldigʻini bir vaqtda koʻpi bilan 10 ta soʻrov bilan olish (GB-7.04).
// importStocks: videoning oxiridagi kod (errgroup.WithContext + SetLimit). fetchPrices: vazifa (TASK.md).
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// workers: yetkazib beruvchiga bir vaqtda koʻpi bilan shuncha soʻrov ketadi (u 10 tadan koʻpini rad etadi).
const workers = 10

// importStocks ids'dagi har mahsulot qoldigʻini yetkazib beruvchidan soʻraydi: bir vaqtda koʻpi bilan workers ta,
// natija ids tartibida. Birinchi xato qolgan soʻrovlarni bekor qiladi va qaytariladi; hammasi oʻtsa, natijani saqlaydi.
func importStocks(ctx context.Context, ids []int64) ([]int, error) {
	stocks := make([]int, len(ids))
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(workers)
	for i, id := range ids {
		g.Go(func() error {
			n, err := checkStock(gctx, id)
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
	if err := saveStocks(ctx, stocks); err != nil {
		return nil, err
	}
	return stocks, nil
}

// VAZIFA (TASK.md): fetchPrices'ni errgroup bilan qayta yozing.
// Hozir u har ID uchun alohida goroutine ochadi: limit yoʻq, xatoda qolganlari toʻxtamaydi, xatoda ID yoʻq.
// Kerak: bir vaqtda koʻpi bilan limit ta soʻrov; natija ids tartibida; birinchi xato «narx <ID>: <xato>»
// koʻrinishida qaytadi va qolgan soʻrovlarni bekor qiladi; limit 1 dan kichik boʻlsa, darhol errBadLimit.
var errBadLimit = errors.New("limit kamida 1 boʻlishi kerak")

func fetchPrices(ctx context.Context, ids []int64, limit int) ([]int64, error) {
	prices := make([]int64, len(ids))
	errs := make([]error, len(ids))
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Go(func() {
			prices[i], errs[i] = checkPrice(ctx, id)
		})
	}
	wg.Wait()
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return prices, nil
}

func main() {
	start := time.Now()
	stocks, err := importStocks(context.Background(), catalogIDs(1000))
	report(stocks, err, start)
}
