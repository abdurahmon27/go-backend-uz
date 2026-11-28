// sizdir — videodagi «sizib chiqqan» jadvallarni va nomzodlar roʻyxatini yasaydi (GB-6.03).
// Hammasi oʻylab topilgan: .test domeni (RFC 2606) va bozor'ga oid soʻzlar, hech kimning haqiqiy paroli emas.
//
//	go run ./_video/sizdir      (papkaning ildizida)
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/abdurahmon27/bozor/internal/password"
)

var words = []string{"olma", "non", "choy", "anor", "uzum", "qovun", "tarvuz", "behi", "gilos", "shaftoli",
	"bodom", "yongoq", "asal", "sut", "qatiq", "guruch", "palov", "somsa", "lagmon", "manti",
	"toshkent", "samarqand", "buxoro", "xiva", "andijon", "namangan", "fargona", "qarshi", "termiz", "nukus",
	"salom", "parol", "qwerty", "admin", "bozor", "dunyo", "sevgi", "baxt", "oila", "vatan"}

var suffixes = []string{"", "1", "12", "123", "1234", "12345", "123456", "!", "01", "007", "2000", "2005",
	"2010", "2015", "2018", "2019", "2020", "2021", "2022", "2023", "2024", "2025", "777", "99", "_uz"}

var users = []struct{ email, pw string }{
	{"dilnoza@bozor.test", "choy2024"},
	{"jasur@bozor.test", "choy2024"},
	{"malika@bozor.test", "tarvuz-osmon-qalam-ariq-27"}, // roʻyxatda yoʻq uzun ibora
}

func main() {
	var list strings.Builder
	for _, w := range words {
		for _, s := range suffixes {
			list.WriteString(w + s + "\n")
		}
	}
	write("nomzodlar.txt", list.String())

	var plain, sha, a2 strings.Builder
	for _, u := range users {
		sum := sha256.Sum256([]byte(u.pw))
		fmt.Fprintf(&plain, "%s %s\n", u.email, u.pw)
		fmt.Fprintf(&sha, "%s %s\n", u.email, hex.EncodeToString(sum[:]))
		fmt.Fprintf(&a2, "%s %s\n", u.email, password.Hash(u.pw))
	}
	write("sizdi/ochiq.txt", plain.String())
	write("sizdi/sha256.txt", sha.String())
	write("sizdi/argon2id.txt", a2.String())
}

func write(path, s string) {
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
