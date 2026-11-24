package kabisa

// kabisa: yil kabisa yilimi? Grigorian qoidasi:
// 4 ga boʻlinadi, lekin 100 ga boʻlinsa, 400 ga ham boʻlinishi shart.
func kabisa(yil int) bool {
	return yil%4 == 0 && (yil%100 != 0 || yil%400 == 0)
}
