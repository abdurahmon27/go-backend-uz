package main

// Chek (GB-7.01'dagi vaʼda). Buyurtma yozilgach xaridorga chek yuboriladi. Email servisi bu laboratoriyada
// simulyatsiya: receiptDelay kutib, log'ga bitta qator yozadi. Javob chekni kutib qolmasin deb chek fonda,
// alohida goroutine'da ketadi.
//
// Har goroutine'ning egasi bor (GB-7.01): chek goroutine'larining egasi — server. createOrder ularni bg.Go bilan
// boshlaydi, run esa toʻxtashda Shutdown'dan keyin waitBG bilan kutadi. «go sendReceipt(…)» boʻlsa, Shutdown
// handler'lar tugashini kutadi, lekin ular boshlagan goroutine'larni bilmaydi: jarayon chiqqanda chek yoʻqoladi.
//
// Goroutine'ga r.Context() emas, logger beriladi: soʻrov context'i handler qaytganda bekor boʻladi
// (go doc net/http.Request.Context), chek esa handler'dan keyin ham ishlaydi.

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/abdurahmon27/bozor/internal/store"
)

// receiptDelay — email servisining javob vaqti (simulyatsiya). Testlar uni qisqartiradi.
var receiptDelay = 2 * time.Second

// sendReceipt — chekni yuborish (simulyatsiya): kutadi va log'ga yozadi.
// Fon goroutine'idagi panic'ni recoverPanic middleware ushlamaydi (u boshqa goroutine'da) va butun server
// yiqiladi (GB-7.01, GB-6.01): shuning uchun chek oʻz panic'ini oʻzi ushlab, log'ga yozadi.
func sendReceipt(logger *slog.Logger, o store.Order) {
	defer func() {
		if v := recover(); v != nil {
			logger.Error("chek: panic", "order_id", o.ID, "panic", v)
		}
	}()
	time.Sleep(receiptDelay)
	logger.Info("chek yuborildi", "order_id", o.ID)
}

// waitBG — bg'dagi hamma fon ishi tugaguncha kutadi, lekin ctx tugasa (Shutdown'ning muddati) ctx.Err() bilan qaytadi.
// bg.Wait()'ning muddati yoʻq, shuning uchun u alohida goroutine'da kutadi va tugaganda done'ni yopadi (GB-7.02);
// select done va ctx.Done()'dan qaysi biri oldin tayyor boʻlsa, shuni oladi (GB-7.03).
// Muddat tugasa, kutayotgan goroutine bg.Wait()'da qoladi: jarayon baribir chiqyapti, u ham birga tugaydi.
// Ikkala holat bir vaqtda tayyor boʻlsa (ish yoʻq, ctx ham tugagan), select ulardan birini tasodifiy tanlaydi
// (Go spec): nil ham, ctx.Err() ham toʻgʻri javob.
func waitBG(ctx context.Context, bg *sync.WaitGroup) error {
	done := make(chan struct{})
	go func() {
		bg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
