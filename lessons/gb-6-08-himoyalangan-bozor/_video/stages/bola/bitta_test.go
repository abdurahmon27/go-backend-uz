// XATO USUL (video, «xato» bobi): bitta foydalanuvchi bilan yozilgan himoya testi.
// U faqat Dilnozaning oʻz buyurtmasini va umuman yoʻq buyurtmani soʻraydi — begona buyurtmani hech kim
// soʻramaydi, shuning uchun BOLA'li store'da ham oʻtadi. Papkada emas: chiqishlar.sh uni vaqtincha qoʻyadi.
package main

import "testing"

func TestOrdersOneUser(t *testing.T) {
	h, _ := testAPI(t)
	dilnoza := login(t, h, dilnozaLogin).Value
	if rec := send(h, "GET", "/orders/41", "", dilnoza); rec.Code != 200 {
		t.Errorf("Dilnoza GET /orders/41 (oʻziniki): %d, kutilgan 200", rec.Code)
	}
	if rec := send(h, "GET", "/orders/999", "", dilnoza); rec.Code != 404 {
		t.Errorf("Dilnoza GET /orders/999 (yoʻq): %d, kutilgan 404", rec.Code)
	}
	if rec := send(h, "GET", "/orders/41", "", ""); rec.Code != 401 {
		t.Errorf("cookie'siz GET /orders/41: %d, kutilgan 401", rec.Code)
	}
}
