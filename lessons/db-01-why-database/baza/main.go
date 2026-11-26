// baza: xuddi shu buyurtmalar, lekin Postgres'da (DB-1.01).
//
//	./bin/baza -tayyorla                     jadvallarni yaratadi va boʻshatadi
//	./bin/baza -usul json -mijoz A           JSON'ni bitta ustunda saqlaydi: oʻqi → qoʻsh → yoz (Go'da)
//	./bin/baza -usul insert -mijoz A         har buyurtma — bitta INSERT, qolganini baza hal qiladi
//	./bin/baza -sana                         bazada nechta buyurtma borligini sanaydi
//	./bin/baza -usul insert -n 0 -tasdiq A.ids   tasdiq olgan har bir id'ni faylga yozadi
//	./bin/baza -tekshir A.ids B.ids          tasdiqlangan id'larning hammasi bazadami?
//
// Ulanish manzili DATABASE_URL muhit oʻzgaruvchisida (.env, SYS-0.02).
package main

import (
	"bufio"
	"bytes"
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

const sxema = `
CREATE TABLE IF NOT EXISTS ombor (id int PRIMARY KEY, doc jsonb NOT NULL);
CREATE TABLE IF NOT EXISTS buyurtmalar (
	id    bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	mijoz text   NOT NULL,
	summa bigint NOT NULL
);
TRUNCATE ombor, buyurtmalar RESTART IDENTITY;
INSERT INTO ombor VALUES (1, '{"buyurtmalar": []}');`

// jsonUsul: fayldagi kod bilan bir xil, faqat fayl oʻrnida jadvaldagi bitta qator.
func jsonUsul(ctx context.Context, db *pgx.Conn, mijoz string) (int64, error) {
	return 0, jsonYoz(ctx, db, mijoz)
}

func jsonYoz(ctx context.Context, db *pgx.Conn, mijoz string) error {
	var b []byte
	if err := db.QueryRow(ctx, `SELECT doc FROM ombor WHERE id = 1`).Scan(&b); err != nil {
		return err
	}
	var o Ombor
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}
	o.Buyurtmalar = append(o.Buyurtmalar, Buyurtma{ID: len(o.Buyurtmalar) + 1, Mijoz: mijoz, Summa: 25000})
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, `UPDATE ombor SET doc = $1 WHERE id = 1`, b)
	return err
}

// insertUsul: oʻqish yoʻq — oʻzgarishni bazaning oʻziga beramiz. Baza yangi qatorning id'sini qaytaradi.
func insertUsul(ctx context.Context, db *pgx.Conn, mijoz string) (int64, error) {
	var id int64
	err := db.QueryRow(ctx, `INSERT INTO buyurtmalar (mijoz, summa) VALUES ($1, 25000) RETURNING id`, mijoz).Scan(&id)
	return id, err
}

// tekshir: tasdiq olgan id'lar (fayllardan) bazada bormi? Har bir id alohida tekshiriladi, faqat soni emas.
func tekshir(ctx context.Context, db *pgx.Conn, fayllar []string) error {
	bazada := map[int64]bool{}
	rows, err := db.Query(ctx, `SELECT id FROM buyurtmalar`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		bazada[id] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	tasdiq, yoq := 0, 0
	for _, f := range fayllar {
		b, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		sc := bufio.NewScanner(bytes.NewReader(b))
		for sc.Scan() {
			var id int64
			fmt.Sscan(sc.Text(), &id)
			tasdiq++
			if !bazada[id] {
				yoq++
			}
		}
	}
	fmt.Printf("tasdiqlangan: %d · bazada: %d · tasdiqlanib yoʻqolgan: %d\n", tasdiq, len(bazada), yoq)
	return nil
}

func main() {
	mijoz := flag.String("mijoz", "A", "kim buyurtma beryapti")
	n := flag.Int("n", 1000, "nechta buyurtma (0 = xato boʻlguncha)")
	usul := flag.String("usul", "insert", "json | insert")
	tayyorla := flag.Bool("tayyorla", false, "jadvallarni yaratish va boʻshatish")
	sana := flag.Bool("sana", false, "faqat sanash")
	tasdiqFayl := flag.String("tasdiq", "", "tasdiq olgan id'lar yoziladigan fayl")
	tekshirish := flag.Bool("tekshir", false, "argumentdagi id fayllarini bazadagi qatorlar bilan solishtirish")
	flag.Parse()

	ctx := context.Background()
	db, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(*mijoz+":", err)
		os.Exit(1)
	}
	defer db.Close(ctx)

	switch {
	case *tayyorla:
		if _, err := db.Exec(ctx, sxema); err != nil {
			fmt.Println("tayyorla:", err)
			os.Exit(1)
		}
		return
	case *tekshirish:
		if err := tekshir(ctx, db, flag.Args()); err != nil {
			fmt.Println("tekshir:", err)
			os.Exit(1)
		}
		return
	case *sana:
		var js, a, b int
		err := db.QueryRow(ctx, `
			SELECT (SELECT jsonb_array_length(doc->'buyurtmalar') FROM ombor WHERE id = 1),
			       count(*) FILTER (WHERE mijoz = 'A'),
			       count(*) FILTER (WHERE mijoz = 'B')
			FROM buyurtmalar`).Scan(&js, &a, &b)
		if err != nil {
			fmt.Println("sana:", err)
			os.Exit(1)
		}
		fmt.Printf("json ustunida: %d ta · buyurtmalar jadvalida: %d ta (A: %d, B: %d)\n", js, a+b, a, b)
		return
	}

	yoz := insertUsul
	if *usul == "json" {
		yoz = jsonUsul
	}
	var idlar *bufio.Writer
	if *tasdiqFayl != "" {
		f, err := os.Create(*tasdiqFayl)
		if err != nil {
			fmt.Println(*mijoz+":", err)
			os.Exit(1)
		}
		defer f.Close()
		idlar = bufio.NewWriter(f)
		defer idlar.Flush()
	}
	i := 0
	for ; *n == 0 || i < *n; i++ {
		id, err := yoz(ctx, db, *mijoz)
		if err != nil {
			fmt.Printf("%s: %d ta tasdiqlandi, keyin xato: %v\n", *mijoz, i, err)
			if idlar != nil {
				idlar.Flush() // os.Exit defer'larni bajarmaydi
			}
			os.Exit(1)
		}
		if idlar != nil {
			fmt.Fprintln(idlar, id) // faqat baza "tayyor" degandan keyin
		}
	}
	fmt.Printf("%s: %d ta yozdim\n", *mijoz, i)
}
