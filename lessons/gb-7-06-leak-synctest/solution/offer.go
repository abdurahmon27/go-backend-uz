package main

import (
	"context"
	"errors"
	"time"
)

// priceOffer — bitta yetkazib beruvchining olma uchun taklifi: qancha vaqtda javob beradi va qancha narx aytadi (tiyin).
type priceOffer struct {
	delay time.Duration
	price int64
}

// priceOffers: uch yetkazib beruvchi. Uchinchisi eng arzon, lekin juda sekin.
var priceOffers = []priceOffer{
	{100 * time.Millisecond, 1_250_000},
	{200 * time.Millisecond, 1_150_000},
	{3 * time.Second, 1_100_000},
}

// quoteTimeout: cheapestPrice takliflarni shuncha vaqtdan koʻp kutmaydi.
const quoteTimeout = time.Second

// errNoOffer: quoteTimeout ichida hech bir yetkazib beruvchi javob bermadi.
var errNoOffer = errors.New("hech bir yetkazib beruvchi vaqtida javob bermadi")

// askOffer bitta yetkazib beruvchidan taklif soʻraydi: o.delay oʻtguncha yoki ctx bekor qilinguncha kutadi.
func askOffer(ctx context.Context, o priceOffer) (int64, error) {
	select {
	case <-time.After(o.delay):
		return o.price, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// cheapestPrice hamma yetkazib beruvchidan birdan soʻraydi, quoteTimeout ichida kelgan takliflarning eng arzonini
// qaytaradi va qaytganda oʻzi boshlagan birorta goroutine'ni ham ortda qoldirmaydi.
func cheapestPrice(ctx context.Context) (int64, error) {
	tctx, cancel := context.WithTimeout(ctx, quoteTimeout)
	defer cancel()
	ch := make(chan int64, len(priceOffers))
	for _, o := range priceOffers {
		go func() {
			p, err := askOffer(tctx, o)
			if err == nil {
				ch <- p
			}
		}()
	}
	var best int64
	found := false
	take := func(p int64) {
		if !found || p < best {
			best, found = p, true
		}
	}
	for range priceOffers {
		select {
		case p := <-ch:
			take(p)
		case <-tctx.Done():
			if err := ctx.Err(); err != nil {
				return 0, err // chaqiruvchi bekor qildi yoki uning muddati tugadi
			}
			for len(ch) > 0 { // muddat bilan bir paytda kelib, buffer'da turgan javoblar
				take(<-ch)
			}
			if !found {
				return 0, errNoOffer
			}
			return best, nil
		}
	}
	return best, nil
}
