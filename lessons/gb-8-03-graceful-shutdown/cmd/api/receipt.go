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

// waitBG — bg'dagi hamma fon ishi tugaguncha kutadi.
// VAZIFA (TASK.md): ctx tugasa (Shutdown'ning muddati), kutishni toʻxtatib ctx.Err() qaytarsin.
// Hozir ctx eʼtiborsiz: chek osilib qolsa, run ham osilib qoladi va docker stop jarayonni SIGKILL bilan oʻldiradi.
func waitBG(ctx context.Context, bg *sync.WaitGroup) error {
	bg.Wait()
	return nil
}
