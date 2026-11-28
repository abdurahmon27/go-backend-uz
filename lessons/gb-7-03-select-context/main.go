// bozor: yetkazib beruvchidan qoldiq soʻrash, select va context bilan vaqtida toʻxtatish (GB-7.03).
// Videoning oxiridagi kod shu; vazifa (TASK.md): sumStock va reportHandler.
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
	// VAZIFA 1 (TASK.md): har id uchun checkStock(ctx, id); xato boʻlsa 0, err; oxirida yigʻindi, nil.
	return 0, errors.New("VAZIFA 1: sumStock hali yozilmagan")
}

// GET /report: reportIDs qoldiqlarining yigʻindisi, "jami: N".
func reportHandler(w http.ResponseWriter, r *http.Request) {
	// VAZIFA 2 (TASK.md): ctx'ni r.Context()'dan reportTimeout muddat bilan yasang (defer cancel()),
	// muddat tugasa 504 Gateway Timeout qaytaring, mijoz ketgan boʻlsa hech narsa yozmay qayting.
	total, err := sumStock(context.Background(), reportIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
