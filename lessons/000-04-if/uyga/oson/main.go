package main

import "fmt"

// Vazifa: chegirma 100 000 soʻmdan emas, 50 000 soʻmdan boshlansin.
// Oʻzgartirishdan oldin taxmin qiling: hozir nima chiqadi, oʻzgartirgandan keyin-chi?
func main() {
	xarid := 70000
	if xarid >= 100000 {
		fmt.Println("chegirma bor")
	} else {
		fmt.Println("chegirma yoʻq")
	}
}
