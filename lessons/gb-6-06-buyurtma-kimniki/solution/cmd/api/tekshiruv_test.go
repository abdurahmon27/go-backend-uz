// Vazifa tekshiruvi (TASK.md): PATCH /orders/{id} faqat oʻz buyurtmangizni oʻzgartirsin. Bu faylni oʻzgartirmang.
// seed.sql: 41 — Dilnozaniki, 42 — Jasurniki.
package main

import (
	"context"
	"net/http"
	"testing"
)

func TestUpdateNoteOwner(t *testing.T) {
	h, pool := testAPI(t)
	dilnoza := login(t, h, dilnozaLogin)
	jasur := login(t, h, jasurLogin)
	note := func(id int) string {
		var s string
		err := pool.QueryRow(context.Background(), "SELECT coalesce(note, '') FROM orders WHERE id = $1", id).Scan(&s)
		if err != nil {
			t.Fatalf("orders %d oʻqilmadi: %v", id, err)
		}
		return s
	}
	before := note(41)

	// Jasur Dilnozaning 41-buyurtmasiga oʻz manzilini yozmoqchi.
	rec := send(h, "PATCH", "/orders/41", `{"note":"Jasurning manzili"}`, jasur.Value)
	if rec.Code != http.StatusNotFound {
		t.Errorf("Jasur PATCH /orders/41 (Dilnozaniki): status %d, kutilgan 404 — begona buyurtma yoʻq buyurtma kabi (TASK.md, 1–2-qadam)", rec.Code)
	}
	if got := note(41); got != before {
		t.Fatalf("Jasur Dilnozaning 41-buyurtmasini oʻzgartirdi: izoh «%s», kutilgan «%s» — UPDATE'da egasi tekshirilmayapti (TASK.md, 1-qadam: WHERE id = $1 AND user_id = $2)", got, before)
	}

	// Oʻz buyurtmasi — oʻzgaradi.
	if rec := send(h, "PATCH", "/orders/41", `{"note":"Chilonzor 9-kvartal, 14-uy"}`, dilnoza.Value); rec.Code != http.StatusNoContent {
		t.Errorf("Dilnoza PATCH /orders/41 (oʻziniki): status %d %s, kutilgan 204 (TASK.md, 2-qadam)", rec.Code, rec.Body)
	}
	if got := note(41); got != "Chilonzor 9-kvartal, 14-uy" {
		t.Errorf("Dilnozaning izohi oʻzgarmadi: «%s» (TASK.md, 2-qadam)", got)
	}

	// Yoʻq buyurtma — oʻsha 404.
	if rec := send(h, "PATCH", "/orders/999", `{"note":"x"}`, dilnoza.Value); rec.Code != http.StatusNotFound {
		t.Errorf("PATCH /orders/999 (yoʻq): status %d, kutilgan 404 (TASK.md, 1-qadam: RowsAffected() == 0)", rec.Code)
	}
	if got := note(42); got != "Yunusobod 4-kvartal, 17-uy, ertalab" {
		t.Errorf("Jasurning 42-buyurtmasi oʻzgarib qoldi: «%s»", got)
	}
}
