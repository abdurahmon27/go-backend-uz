`password.go.txt` — `internal/password/password.go`ning nusxasi, bitta qator farq bilan: `decode` salt'ni
`base64.StdEncoding` bilan oʻqiydi, `Hash` esa `RawStdEncoding` (oxirida `=` yoʻq) bilan yozgan. Natija videoda:
`illegal base64 data at input byte 20`. `record.sh` uni vaqtinchalik nusxaga qoʻyib `go test` qiladi.
