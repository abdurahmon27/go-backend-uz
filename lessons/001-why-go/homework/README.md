# 01 · Uy vazifasi javoblari

Avval o'zingiz urinib ko'ring. 🙂

## 🟢 Oson
[`easy/main.go`](easy/main.go). `fmt.Println` bir nechta argumentni bo'sh joy bilan ajratib chiqaradi:
`fmt.Println("Yoshim:", 17)` → `Yoshim: 17`.

## 🟡 O'rta
[`medium.sh`](medium.sh). `GOOS`/`GOARCH` maqsad tizimni bildiradi, compiler esa shu Linux'ning o'zida ishlaydi:

```
2409998 hello-linux-amd64
2473472 hello-windows-amd64.exe
2492482 hello-darwin-arm64
```

Hajmlar deyarli bir xil, chunki har birining ichida o'sha platforma uchun runtime bor.

## 🟠 Qiyin: `-ldflags=-s`
```
go build                     2,409,998 bayt
go build -ldflags=-w         1,733,756 bayt   (faqat DWARF olib tashlandi)
go build -ldflags=-s         1,601,698 bayt   (−33.5%)
go build -ldflags="-s -w"    1,601,698 bayt   (xuddi shu)
```
- `-s`: **symbol table** olib tashlanadi (`go tool nm` endi "no symbols" deydi). Go 1.22'dan beri `-s` o'zi `-w`ni
  ham yoqadi ([cmd/link/doc.go](https://github.com/golang/go/blob/go1.26.8/src/cmd/link/doc.go#L123)), shuning
  uchun `-s` va `-s -w` bir xil. Farqni ko'rish uchun `-s -w=0` bilan solishtiring.
- `-w`: **DWARF** debug ma'lumoti olib tashlanadi (debugger, ya'ni `dlv`, endi qatorlarni bilmaydi).

**Nega `.gopclntab` qoldi?** `size -A` bilan qarang: `-s` dan keyin ham `.gopclntab` 677 386 bayt. Bu jadvalni
runtime'ning o'zi ishlatadi: GC va stack o'sishi stack xaritalarini shu yerdan o'qiydi
([runtime/stkframe.go](https://github.com/golang/go/blob/go1.26.8/src/runtime/stkframe.go), `getStackMap`),
panic stack trace'dagi funksiya nomlari va qator raqamlari ham shundan. Ya'ni `-s` debug *qulayligini* oladi,
runtime'ga kerak narsani emas.

## 🔴 Guru: `main.main` ichida 20 ta instruction
To'liq chiqish: [`../asm/main.main.s`](../asm/main.main.s) (`go1.26.8 linux/amd64`).

| # | Instruction | Nima qiladi |
|---|---|---|
| 1 | `CMPQ SP, 0x10(R14)` | `R14` = joriy goroutine (`g`), `0x10(R14)` = `g.stackguard0`. Stack chegarasiga yaqinlashdikmi? |
| 2 | `JBE 0x49e18c` (manzil har build'da boshqacha bo'lishi mumkin) | Ha bo'lsa, oxiridagi `morestack`ga sakra (19-qator) |
| 3–4 | `PUSHQ BP` · `MOVQ SP, BP` | Frame pointer: profiler va debugger stack'ni shu orqali yuradi |
| 5 | `SUBQ $0x38, SP` | Funksiya uchun 56 bayt stack joy |
| 6–9 | `LEAQ … DX` · `MOVQ DX, …(SP)` | `"Salom, dunyo!"` uchun `any` qiymat: type pointer + data pointer |
| 10 | `MOVQ os.Stdout(SB), BX` | `Println` aslida `Fprintln(os.Stdout, …)`, va compiler uni **inline** qilgan |
| 11 | `LEAQ go:itab.*os.File,io.Writer(SB), AX` | `*os.File` → `io.Writer` uchun tayyor `itab` |
| 12–14 | `LEAQ 0x28(SP), CX` · `MOVL $1, DI` · `MOVQ DI, SI` | `[]any` slice: pointer, `len = 1`, `cap = 1` |
| 15 | `CALL fmt.Fprintln(SB)` | Haqiqiy chiqarish |
| 16–18 | `ADDQ $0x38, SP` · `POPQ BP` · `RET` | Stack'ni qaytarib, chiqish |
| 19 | `CALL runtime.morestack_noctxt.abi0(SB)` | Stack'ni kattalashtir (kattaroq stack'ga ko'chir) |
| 20 | `JMP main.main(SB)` | …va `main`ni **boshidan** qayta boshla |

Asosiy g'oya: goroutine kichik stack bilan boshlanadi, shuning uchun deyarli har bir funksiya boshida
1–2-qatordagi tekshiruv bor. Uni compiler o'zi qo'shadi.
