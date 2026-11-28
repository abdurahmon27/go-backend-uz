// bozor: yetkazib beruvchidan 1000 mahsulot qoldigʻini bir vaqtda koʻpi bilan 10 ta soʻrov bilan olish (GB-7.04).
// importStocks: videoning oxiridagi kod (errgroup.WithContext + SetLimit). Javob (solution): fetchPrices yozilgan.
package main

import (
	"context"
	"errors"
	"fmt"
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

// errBadLimit: limit 1 dan kichik. errgroup'da SetLimit(0) bilan g.Go hech qachon qaytmaydi (boshqa goroutine
// boʻlmasa, runtime «deadlock» bilan toʻxtatadi), manfiy limit esa limitni olib tashlaydi: ikkalasi ham bu yerda xato.
var errBadLimit = errors.New("limit kamida 1 boʻlishi kerak")

// fetchPrices ids'dagi har mahsulot narxini (tiyin) yetkazib beruvchidan soʻraydi: bir vaqtda koʻpi bilan limit ta,
// natija ids tartibida. Birinchi xato «narx <ID>: …» koʻrinishida qaytadi va qolgan soʻrovlarni bekor qiladi.
func fetchPrices(ctx context.Context, ids []int64, limit int) ([]int64, error) {
	if limit < 1 {
		return nil, errBadLimit
	}
	prices := make([]int64, len(ids))
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(limit)
	for i, id := range ids {
		g.Go(func() error {
			p, err := checkPrice(gctx, id)
			if err != nil {
				return fmt.Errorf("narx %d: %w", id, err)
			}
			prices[i] = p
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return prices, nil
}

func main() {
	start := time.Now()
	stocks, err := importStocks(context.Background(), catalogIDs(1000))
	report(stocks, err, start)
}
