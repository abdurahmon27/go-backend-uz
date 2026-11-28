// bozor: yetkazib beruvchidan qoldiq soʻrash, select va context bilan vaqtida toʻxtatish (GB-7.03).
// Javob (solution): sumStock va reportHandler yozilgan.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// checkStock yetkazib beruvchining javobini yoki ctx bekor qilinishini kutadi: qaysi biri oldin boʻlsa.
func checkStock(ctx context.Context, id int64) (int, error) {
	select {
	case n := <-supplier(id):
		return n, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// GET /stock: oʻrikning (ID 3) qoldigʻi. Mijoz ketsa, r.Context() bekor qilinadi va ish toʻxtaydi.
func stockHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	n, err := checkStock(r.Context(), 3)
	if err != nil {
		fmt.Println("toʻxtadi:", err, time.Since(start).Round(time.Second))
		return
	}
	fmt.Println("javob yozildi", time.Since(start).Round(time.Second))
	fmt.Fprintln(w, "qoldiq:", n)
}

// reportIDs: hisobotdagi mahsulotlar; reportTimeout: GET /report shuncha vaqtdan koʻp kutmaydi.
var (
	reportIDs     = []int64{1, 2, 3}
	reportTimeout = 5 * time.Second
)

// sumStock ids'dagi har mahsulot qoldigʻini checkStock bilan ketma-ket soʻrab, yigʻindisini qaytaradi.
// ctx bekor qilinsa yoki muddati tugasa, darhol 0 va ctx'ning xatosini qaytaradi.
func sumStock(ctx context.Context, ids []int64) (int, error) {
	total := 0
	for _, id := range ids {
		n, err := checkStock(ctx, id)
		if err != nil {
			return 0, err
		}
		total += n
	}
	return total, nil
}

// GET /report: reportIDs qoldiqlarining yigʻindisi, "jami: N".
func reportHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), reportTimeout)
	defer cancel()
	total, err := sumStock(ctx, reportIDs)
	if errors.Is(err, context.DeadlineExceeded) {
		http.Error(w, "yetkazib beruvchi javob bermadi", http.StatusGatewayTimeout)
		return
	}
	if err != nil {
		return // mijoz ketdi (context.Canceled): javob yozishning ma'nosi yoʻq
	}
	fmt.Fprintln(w, "jami:", total)
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /stock", stockHandler)
	mux.HandleFunc("GET /report", reportHandler)
	return mux
}

func main() {
	log.Fatal(http.ListenAndServe("localhost:8093", newMux()))
}
