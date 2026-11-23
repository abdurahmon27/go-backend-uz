# 01 · Go nega bor? `go build` ichida nima bo'ladi + Hello World

| Papka | Nima |
|---|---|
| `hello/` | videodagi dastur |
| `notmain/` | pauza-savol: `package main` → `package hello` |
| `bug/` | debugging mashqi: ikkita xato |
| `asm/main.main.s` | `go tool objdump -s main.main` to'liq chiqishi (20 ta instruction) |
| `homework/` | uy vazifasi javoblari |

## Videodagi buyruqlar
```bash
go run ./hello
go run ./notmain                       # xato: not a main package
go build -o hello/hello ./hello && ls -lh hello/hello
ldd hello/hello                        # not a dynamic executable
GOOS=windows go build -o hello/hello.exe ./hello
go run ./bug                           # 2 ta compile error
go tool objdump -s 'main.main$' hello/hello
```

## Uy vazifasi
- 🟢 **Oson:** ismingiz va yoshingizni chiqaradigan dastur yozing.
- 🟡 **O'rta:** dasturni Windows, Linux va macOS uchun build qiling va hajmlarini solishtiring.
- 🟠 **Qiyin:** `-ldflags="-s -w"` bilan build qiling. Hajm qancha kamaydi, va nima olib tashlandi?
- 🔴 **Guru:** `go tool objdump` bilan `main.main`ni oching va har bir instruction nima qilishini tushuntiring.

Avval o'zingiz urinib ko'ring. Javoblar: [homework/](homework/).
