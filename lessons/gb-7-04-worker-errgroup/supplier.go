package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// supplier.go: yetkazib beruvchining API'si, bu laboratoriyada soxta (tarmoqsiz). Bekor qilishdagi soʻrovlar soni
// va cheklovsiz parallel ishlar rejalashtirishga bogʻliq. GB-7.03'dagi yetkazib beruvchining tezroq va katta nusxasi:
// olma (1), non (2) va oʻrik (3) qoldigʻi
// oʻsha, qolgan ID'lar uchun qoldiq = ID % 50. Har soʻrov supplierDelay davom etadi. Bir vaqtda supplierLimit'dan koʻp soʻrov kelsa, ortigʻi
// rad etiladi (ErrTooMany, HTTP'dagi 429 Too Many Requests kabi: GB-6.07). Yetkazib beruvchi faqat
// 1..supplierProducts ID'larni biladi. stats hamma raqamni yozib boradi: nechta soʻrov keldi, eng koʻpi bilan
// nechtasi bir vaqtda ishladi, nechtasi rad etildi.

const (
	supplierDelay    = 20 * time.Millisecond
	supplierLimit    = 10
	supplierProducts = 1000
)

var (
	ErrTooMany  = errors.New("429 too many requests")
	ErrNotFound = errors.New("mahsulot topilmadi")
)

// supplierStats: yetkazib beruvchi tomonidagi hisoblagichlar; mu ularni himoyalaydi (GB-4.07).
type supplierStats struct {
	mu            sync.Mutex
	calls         int // kelgan soʻrovlar
	active        int // hozir ishlayotganlari
	peak          int // eng koʻpi bilan bir vaqtda
	rejected      int // rad etilganlari (ErrTooMany)
	peakGoroutine int // dasturdagi goroutine'lar soni, eng koʻpi
}

var stats supplierStats

// call bitta soʻrovni taqlid qiladi: delay kutadi yoki ctx bekor qilinsa darhol qaytadi (GB-7.03).
// ctx allaqachon bekor qilingan boʻlsa, soʻrov umuman yuborilmaydi va hisoblanmaydi.
func call(ctx context.Context, id int64, delay time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	stats.mu.Lock()
	stats.calls++
	stats.active++
	stats.peak = max(stats.peak, stats.active)
	stats.peakGoroutine = max(stats.peakGoroutine, runtime.NumGoroutine())
	tooMany := stats.active > supplierLimit
	stats.mu.Unlock()
	defer func() {
		stats.mu.Lock()
		stats.active--
		stats.mu.Unlock()
	}()

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return ctx.Err()
	}
	if tooMany {
		stats.mu.Lock()
		stats.rejected++
		stats.mu.Unlock()
		return ErrTooMany
	}
	if id < 1 || id > supplierProducts {
		return ErrNotFound
	}
	return nil
}

// stocks: GB-7.03'dagi qoldiqlar (ID → dona): olma 1, non 2, oʻrik 3.
var stocks = map[int64]int{1: 40, 2: 120, 3: 5}

// checkStock id'li mahsulotning yetkazib beruvchidagi qoldigʻini soʻraydi (dona); ctx'ni hurmat qiladi (GB-7.03).
func checkStock(ctx context.Context, id int64) (int, error) {
	if err := call(ctx, id, supplierDelay); err != nil {
		return 0, err
	}
	if n, ok := stocks[id]; ok {
		return n, nil
	}
	return int(id % 50), nil
}

// checkPrice id'li mahsulotning yetkazib beruvchidagi narxini soʻraydi, tiyinda (vazifa uchun; tezroq: 5 ms).
func checkPrice(ctx context.Context, id int64) (int64, error) {
	if err := call(ctx, id, 5*time.Millisecond); err != nil {
		return 0, err
	}
	return id * 1000 * tiyinPerSom, nil // ID × 1000 soʻm
}

const tiyinPerSom = 100 // 1 soʻm = 100 tiyin (GB-1.03)

// saveStocks qoldiqlarni bazaga yozadi. Bu laboratoriyada soxta: faqat ctx'ni tekshiradi
// (haqiqiy Postgres'ga yozish GB-7.07'da).
func saveStocks(ctx context.Context, stocks []int) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("saqlash: %w", err)
	}
	return nil
}

// catalogIDs import faylidagi mahsulot ID'larini qaytaradi: 1, 2, …, n.
func catalogIDs(n int) []int64 {
	ids := make([]int64, n)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	return ids
}

// resetStats hisoblagichlarni nolga qaytaradi (testlar uchun).
func resetStats() {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.calls, stats.active, stats.peak, stats.rejected, stats.peakGoroutine = 0, 0, 0, 0, 0
}

// report import natijasini, yetkazib beruvchi hisoblagan raqamlarni va dastur sarflagan resursni chiqaradi:
// eng koʻp goroutine soni va Go runtime operatsion tizimdan olgan xotira (runtime.MemStats.Sys).
func report(stocks []int, err error, start time.Time) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	fmt.Printf("soʻrov: %d · bir vaqtda: %d · rad: %d · vaqt: %v\n",
		stats.calls, stats.peak, stats.rejected, time.Since(start).Round(time.Millisecond))
	if err != nil {
		fmt.Println("xato:", firstLine(err))
	} else {
		fmt.Println("qoldiqlar:", len(stocks))
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("goroutine: %d · xotira: %d MB\n", stats.peakGoroutine, m.Sys>>20)
}

// firstLine koʻp qatorli xatoning (errors.Join, GB-2.08) faqat birinchi qatorini va qolganlar sonini qaytaradi.
func firstLine(err error) string {
	lines := strings.Split(err.Error(), "\n")
	if len(lines) == 1 {
		return lines[0]
	}
	return fmt.Sprintf("%s (va yana %d ta)", lines[0], len(lines)-1)
}
