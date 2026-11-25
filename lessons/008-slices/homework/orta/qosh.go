package orta

// qosh s ning oxiriga v ni qoʻshadi va hech qachon chaqiruvchining massiviga yozmaydi:
// uch indeksli ifoda cap'ni len'ga tenglaydi, shuning uchun append yangi massiv ajratadi.
func qosh(s []string, v string) []string {
	return append(s[:len(s):len(s)], v)
}
