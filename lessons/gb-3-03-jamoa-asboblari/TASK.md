# Vazifa · GB-3.03: jamoadoshning PR'ini uchta asbob bilan tekshiring

Videoda bozor'dagi uchta xatoni uchta asbob topdi: `go vet` takrorlangan JSON tag'ini, `staticcheck` tekshirilmay
qolgan `err`'ni, `govulncheck` esa `golang.org/x/text`'ning eski versiyasidagi zaiflikni. Papka videoning oxiridagi
holatda, ustiga **jamoadoshingizning PR'i** qoʻshilgan:

- `internal/store/rename.go` — yangi `Rename(id int64, name string) error` metodi: mahsulot nomini almashtiradi,
  nom chetidagi boʻsh joylarni olib tashlaydi, boʻsh nom boʻlsa `error` qaytaradi;
- PR eski branch'dan kelgani uchun `go.mod`'da `golang.org/x/text` yana `v0.3.7` boʻlib qolgan.

Kod kompilyatsiya boʻladi va eski test (`store_test.go`) oʻtadi. Commit'dan oldingi tekshiruvni oʻzingiz qiling.

**Nima qilish kerak** (papkaning ildizida):

1. `go vet ./...` — bitta xabar chiqadi. Uni oʻqing va `rename.go`'dagi qatorni tuzating.
2. `staticcheck ./...` — vet koʻrsatgan qatorni u ham koʻrsatadi, lekin yana bitta xato topadi, uni vet koʻrmaydi.
   Ikkalasini tuzating.
3. `govulncheck ./...` — qaysi zaiflik, qaysi qator uni chaqiradi? Modulni yangilang.

Uchala buyruq jim boʻlsin: `go vet` va `staticcheck` hech narsa chiqarmasin, `govulncheck` faqat bitta qator chiqarsin:
`No vulnerabilities found.` (Agar uning ostida yana "This scan also found … vulnerabilities in modules you require"
qatori boʻlsa, versiya hali eski.) Keyin:

```bash
go test ./...
```

`ok` chiqquncha ishlang. Commit'dan oldingi toʻliq qator videodagidek:

```bash
go vet ./... && staticcheck ./... && govulncheck ./... && go test ./...
```

Test fayllari `internal/store/tekshiruv_test.go` va `cmd/bozor/tekshiruv_test.go`: ularni oʻzgartirmang. Hozir
`go test ./...` shunday yiqiladi:

```
--- FAIL: TestVazifaXText (0.00s)
    tekshiruv_test.go:35: go.mod: golang.org/x/text v0.3.7 — govulncheck bu versiyada zaiflik topadi (main language.MatchStrings'ni chaqiradi); kamida v0.41.0 kerak. govulncheck ./... nima deydi? (TASK.md, 3-qadam)
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/bozor	0.002s
--- FAIL: TestVazifaRename (0.00s)
    tekshiruv_test.go:29: Rename(1, «  oʻrik »)'dan keyin nom hali «olma»: nom almashmadi. go vet ./... nima deydi? (TASK.md, 1-qadam)
--- FAIL: TestVazifaRenameBosh (0.00s)
    tekshiruv_test.go:42: Rename(2, «   ») = nil; kutilgan error: boʻsh joylar olib tashlansa, nom boʻsh qoladi (TASK.md, 2-qadam)
FAIL
FAIL	github.com/abdurahmon27/bozor/internal/store	0.002s
FAIL
```

**Maslahatlar**

- Asboblar hali oʻrnatilmagan boʻlsa, **shu papkaning ichida** oʻrnating (videodagidek, versiyasi bilan):
  `go install honnef.co/go/tools/cmd/staticcheck@2026.2.2` va `go install golang.org/x/vuln/cmd/govulncheck@v1.9.0`.
  Papka ichida `go` buyrugʻi `go.mod`'dagi Go 1.27'ga oʻtadi va asbob shu Go bilan quriladi (sababi README'da).
  `command not found` chiqsa: `export PATH="$PATH:$(go env GOPATH)/bin"`.
- staticcheck xabari oxirida qavs ichida tekshiruv nomi turadi. `staticcheck -explain <nom>` uning hujjatini
  terminalda koʻrsatadi (videodagi `SA4006` kabi). Funksiya nima qaytarishini koʻrish: `go doc strings.TrimSpace`.
- Modulni yangilash: `go get golang.org/x/text@latest`. `go.mod` va `go.sum` oʻzi yangilanadi; keyin yana
  `govulncheck ./...` va `go test ./...`.
- `govulncheck` bazani vuln.go.dev'dan oladi: internet kerak.

Javob: [`solution/`](solution) (`diff -ru internal solution/internal`, `diff go.mod solution/go.mod`).
