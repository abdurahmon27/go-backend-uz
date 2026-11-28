// bozor: bitta mahsulot narxini uch yetkazib beruvchidan birdan soʻrab, eng tez javobni olish (GB-7.06).
// Videoning oxiridagi holat: goroutine'lar fastestPrice bilan birga tugaydi (ctx + cancel), leak yoʻq.
package main

import (
	"context"
	"time"
)

// prices: yetkazib beruvchilar aytadigan narx, tiyinda (ID → tiyin): olma 1, non 2, oʻrik 3.
var prices = map[int64]int64{1: 1_200_000, 2: 500_000, 3: 3_000_000}

// supplierDelays: uch yetkazib beruvchining (chorsu, oloy, siyob) javob vaqti.
var supplierDelays = []time.Duration{
	100 * time.Millisecond,
	200 * time.Millisecond,
	300 * time.Millisecond,
}

// quote bitta yetkazib beruvchidan narx soʻrashni taqlid qiladi: delay oʻtguncha yoki ctx bekor qilinguncha kutadi.
func quote(ctx context.Context, delay time.Duration, id int64) (int64, error) {
	select {
	case <-time.After(delay):
		return prices[id], nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// fastestPrice hamma yetkazib beruvchidan birdan soʻraydi va birinchi kelgan javobni qaytaradi.
// Qaytishdan oldin cancel() qolgan soʻrovlarni toʻxtatadi: oʻzi boshlagan goroutine'lar uning bilan tugaydi.
func fastestPrice(ctx context.Context, id int64) (int64, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ch := make(chan int64, len(supplierDelays))
	for _, d := range supplierDelays {
		go func() {
			p, err := quote(ctx, d, id)
			if err == nil {
				ch <- p
			}
		}()
	}
	select {
	case p := <-ch:
		return p, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
