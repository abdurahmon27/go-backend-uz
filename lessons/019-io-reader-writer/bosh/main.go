// bosh: hech narsa qilmaydigan Go dasturi — taqqoslash uchun: Go runtime'ning oʻzi qancha xotira oladi?
package main

import (
	"fmt"

	"github.com/abdurahmon27/go-backend-uz/lessons/019-io-reader-writer/ozim"
)

func main() {
	fmt.Printf("boʻsh dastur · eng koʻp xotira: %.1fMiB\n", ozim.RSS())
}
