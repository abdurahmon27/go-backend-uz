// Uy vazifasi javoblari (Oʻrta va Qiyin): Postgres'dagi JSON ustunini yoʻqotishsiz yangilash.
//
//	go run ./homework/baza-tuzat -usul sql -mijoz A        Oʻrta: oʻzgarishni bitta UPDATE bilan bazaga beramiz
//	go run ./homework/baza-tuzat -usul forupdate -mijoz A  Qiyin: transaction + SELECT … FOR UPDATE
//
// Jadvalni avval ./bin/baza -tayyorla yaratadi; natijani ./bin/baza -sana sanaydi.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Buyurtma struct {
	ID    int    `json:"id"`
	Mijoz string `json:"mijoz"`
	Summa int    `json:"summa"`
}

type Ombor struct {
	Buyurtmalar []Buyurtma `json:"buyurtmalar"`
}

// sqlUsul: Go hech narsani oʻqimaydi. Oʻqish, qoʻshish va yozish bitta UPDATE ichida:
// Postgres qatorni yangilash paytida qulflaydi, ikkinchi UPDATE birinchisi tugashini kutadi
// va yangi qiymatning ustiga qoʻshadi (Read Committed, PostgreSQL docs 13.2.1).
func sqlUsul(ctx context.Context, db *pgx.Conn, mijoz string) error {
	_, err := db.Exec(ctx, `
		UPDATE ombor
		SET doc = jsonb_set(doc, '{buyurtmalar}', doc->'buyurtmalar' || jsonb_build_object(
			'id', jsonb_array_length(doc->'buyurtmalar') + 1, 'mijoz', $1::text, 'summa', 25000))
		WHERE id = 1`, mijoz)
	return err
}

// forUpdateUsul: Go kodi oʻzgarmaydi, lekin oʻqish va yozish bitta transaction'da,
// SELECT … FOR UPDATE esa qatorni commit'gacha qulflaydi — ikkinchi jarayon kutadi.
func forUpdateUsul(ctx context.Context, db *pgx.Conn, mijoz string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // commit'dan keyin hech narsa qilmaydi
	var b []byte
	if err := tx.QueryRow(ctx, `SELECT doc FROM ombor WHERE id = 1 FOR UPDATE`).Scan(&b); err != nil {
		return err
	}
	var o Ombor
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}
	o.Buyurtmalar = append(o.Buyurtmalar, Buyurtma{ID: len(o.Buyurtmalar) + 1, Mijoz: mijoz, Summa: 25000})
	if b, err = json.Marshal(o); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE ombor SET doc = $1 WHERE id = 1`, b); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func main() {
	mijoz := flag.String("mijoz", "A", "kim buyurtma beryapti")
	n := flag.Int("n", 1000, "nechta buyurtma")
	usul := flag.String("usul", "sql", "sql | forupdate")
	flag.Parse()

	ctx := context.Background()
	db, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(*mijoz+":", err)
		os.Exit(1)
	}
	defer db.Close(ctx)

	yoz := sqlUsul
	if *usul == "forupdate" {
		yoz = forUpdateUsul
	}
	for i := 0; i < *n; i++ {
		if err := yoz(ctx, db, *mijoz); err != nil {
			fmt.Printf("%s: %d ta yozildi, keyin xato: %v\n", *mijoz, i, err)
			os.Exit(1)
		}
	}
	fmt.Printf("%s: %d ta yozdim\n", *mijoz, *n)
}
