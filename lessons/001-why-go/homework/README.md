# 01 · Uy vazifasi javoblari

Avval o'zingiz urinib ko'ring. 🙂

## 🟢 Oson
[`easy/main.go`](easy/main.go). `fmt.Println` bir nechta argumentni bo'sh joy bilan ajratib chiqaradi:
`fmt.Println("Yoshim:", 17)` → `Yoshim: 17`.

## 🟡 O'rta
[`medium.sh`](medium.sh). `GOOS`/`GOARCH` maqsad tizimni bildiradi, compiler esa shu Linux'ning o'zida ishlaydi:

```
2.3M hello-linux-amd64
2.4M hello-windows-amd64.exe
2.4M hello-darwin-arm64
```

Hajmlar deyarli bir xil, chunki har birining ichida o'sha platforma uchun runtime bor.

## 🟠 Qiyin: `-ldflags="-s -w"`
```
go build                     2,410,014 bayt
go build -ldflags="-s -w"    1,601,698 bayt   (−33.5%)
```
- `-s`: **symbol table** olib tashlanadi (`go tool nm` endi hech narsa ko'rsatmaydi).
- `-w`: **DWARF** debug ma'lumoti olib tashlanadi (debugger, ya'ni `dlv`, endi qatorlarni bilmaydi).

**Guru tuzog'i:** `panic` bo'lsa, stack trace baribir funksiya nomlari va qator raqamlarini ko'rsatadi. Nega?
Runtime'ning o'zi `pclntab` jadvalidan foydalanadi (stack trace va GC uchun), va linker uni olib tashlamaydi.
Ya'ni `-s -w` debug *qulayligini* oladi, runtime'ga kerak narsani emas.

## 🔴 Guru: `main.main` ichida 20 ta instruction
To'liq chiqish: [`../asm/main.main.s`](../asm/main.main.s) (`go1.26.8 linux/amd64`).

| # | Instruction | Nima qiladi |
|---|---|---|
| 1 | `CMPQ SP, 0x10(R14)` | `R14` = joriy goroutine (`g`), `0x10(R14)` = `g.stackguard0`. Stack chegarasiga yaqinlashdikmi? |
| 2 | `JBE 0x49e18c` | Ha bo'lsa, oxiridagi `morestack`ga sakra (19-qator) |
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
