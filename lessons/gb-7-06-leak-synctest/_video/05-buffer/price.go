package main

import "time"

// prices: yetkazib beruvchilar aytadigan narx, tiyinda (ID → tiyin): olma 1, non 2, oʻrik 3.
var prices = map[int64]int64{1: 1_200_000, 2: 500_000, 3: 3_000_000}

// supplierDelays: uch yetkazib beruvchining (chorsu, oloy, siyob) javob vaqti.
var supplierDelays = []time.Duration{
	100 * time.Millisecond,
	200 * time.Millisecond,
	300 * time.Millisecond,
}

func quote(delay time.Duration, id int64) int64 {
	time.Sleep(delay)
	return prices[id]
}

func fastestPrice(id int64) int64 {
	ch := make(chan int64, len(supplierDelays))
	for _, d := range supplierDelays {
		go func() {
			ch <- quote(d, id)
		}()
	}
	return <-ch
}
