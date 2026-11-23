# Dastur: 0 → senior backend engineer

Universitetdagi computer science, lekin oddiy tilda, hammasi Go'da. Har bir modul eng yaxshi universitet
kurslariga tayanadi. Biz ularni o'zbek tilida, kod bilan va "nega?" savoliga javob berib qayta quramiz.

**Asosiy tamoyil:** AI kod yozadi, engineer esa tushunadi. Har bir darsda: *nega* shunday ishlaydi, ichkarida nima
bo'ladi, pauza-savol, debugging mashqi va to'rt darajali uy vazifasi (oson / o'rta / qiyin / guru).

**Ritm:** har 5 ta darsdan keyin 1 ta LeetCode darsi, ya'ni masalaga *qanday yondashish*: pattern, brute force → optimal,
complexity. Modul oxirida intervyu savollari.

| # | Modul | Tayanch kurs | Darslar |
|---|---|---|---|
| 1 | Go asoslari | Go spec, *The Go Programming Language* | ~25 |
| 2 | Kompyuter qanday ishlaydi | CMU 15-213 (*CS:APP*) | ~14 |
| 3 | Algoritmlar va data structures | MIT 6.006 | ~35 |
| 4 | Operatsion tizim va concurrency | Wisconsin CS537 (*OSTEP*) | ~22 |
| 5 | Tarmoqlar | Stanford CS144 | ~16 |
| 6 | Ma'lumotlar bazalari | CMU 15-445 | ~22 |
| 7 | Backend engineering | ishlab chiqarish tajribasi | ~24 |
| 8 | Distributed systems | MIT 6.5840 (6.824) | ~14 |
| 9 | System design | intervyu va real tizimlar | ~14 |
| 10 | Intervyu | LeetCode patternlar, mock intervyular | doim |

---

## 1. Go asoslari
01. Go nega bor? `go build` ichida nima bo'ladi + Hello World ✅
02. Variables & types: `int` xotirada qanday saqlanadi, overflow
03. Constants va `iota`: untyped constant nima
04. Strings, bytes, runes: UTF-8 ichidan (`oʻ` necha byte?)
05. Control flow: `if`, `switch`, Go'da nega faqat `for`
LC-01. Two Sum: hash map pattern
06. Functions: bir nechta return, closures, stack frame
07. Arrays vs slices: backing array, `len`/`cap`, `append` qanday o'sadi
08. Maps: hash table ichidan (Go 1.24 Swiss tables)
09. Structs va xotira: alignment, padding
10. Pointers: stack vs heap, escape analysis
LC-02. Contains Duplicate / Valid Anagram: counting
11. Methods: value vs pointer receiver
12. Interfaces: `itab`, dynamic dispatch, nil interface tuzog'i
13. Errors: error as value, wrapping, `errors.Is`/`As`
14. `defer`, `panic`, `recover` ichidan
15. Packages va modules: `go.mod`, versiyalar, visibility
LC-03. Best Time to Buy and Sell Stock: bir o'tish
16. Testing: table tests, coverage
17. Benchmark va `pprof`: taxmin emas, o'lchash
18. Generics: qachon kerak, qachon kerak emas
19. `io.Reader`/`io.Writer`: Go'ning eng chiroyli abstraksiyasi
20. Fayllar, `bufio`, JSON (reflection ichidan)
LC-04. Valid Parentheses: stack
21. Tooling: `go vet`, `staticcheck`, race detector
22. Notanish kodni o'qish: standart kutubxona manbasini o'qiymiz
23. Loyiha: CLI tool, boshidan oxirigacha
24. AI yozgan Go kodini review qilish: 10 ta tipik xato
25. Modul 1: intervyu savollari

## 2. Kompyuter qanday ishlaydi
Binary va hex · two's complement · floating point (`0.1 + 0.2`) · CPU va assembly (Go asm) · cache va xotira
ierarxiyasi (qator vs ustun bo'yicha o'qish benchmark) · virtual memory · process va syscall (`strace` bilan Go binary) ·
linking va ELF · compiler optimizatsiyalari · garbage collection (Go'ning tri-color GC'si)

## 3. Algoritmlar va data structures
Big-O (o'lchab ko'ramiz) · dynamic array va amortized · linked list · stack/queue · hashing · recursion va call stack ·
sorting: insertion, merge, quick, heap, counting/radix · binary search (javob bo'yicha ham) · tree va BST · balanced
tree g'oyasi · heap va priority queue · trie · graph: BFS/DFS · topological sort · shortest path: Dijkstra,
Bellman-Ford · MST: Kruskal (union-find), Prim · dynamic programming (5 dars) · greedy · backtracking · bit
manipulation · patternlar: two pointers, sliding window, prefix sum, monotonic stack

## 4. Operatsion tizim va concurrency
Process vs thread · scheduling · goroutine va GMP scheduler · channel ichidan (`hchan`) · `select` · mutex va
atomics · race condition va race detector · deadlock · `WaitGroup`, `Once`, `Cond`, `RWMutex` · `context` ·
patternlar: worker pool, pipeline, fan-in/out, `errgroup`, semaphore · Go memory model · fayl tizimi va `fsync`

## 5. Tarmoqlar
Qatlamlar · IP va routing · TCP: handshake, flow va congestion control · UDP · DNS (o'z resolverimizni yozamiz) · TLS ·
HTTP/1.1 (raw TCP ustida server yozamiz) · HTTP/2 va HTTP/3 · `net/http` ichidan · WebSocket · gRPC va protobuf

## 6. Ma'lumotlar bazalari
Relational model va SQL · PostgreSQL amalda · storage: page, heap file · index: B+tree (o'zimiz yozamiz) ·
`EXPLAIN ANALYZE` · join algoritmlari · transaction va ACID · isolation levels va anomaliyalar · MVCC · lock va
deadlock · WAL va recovery · LSM tree (mini KV store yozamiz) · NoSQL trade-off'lari · Redis ichidan

## 7. Backend engineering
REST API dizayni · auth: session, JWT, OAuth · validation va xatolar · Go loyiha tuzilishi · `database/sql`, sqlc,
migrations · caching strategiyalari · message queue: Kafka, NATS, RabbitMQ · background jobs · rate limiting ·
idempotency va retry · observability: logs, metrics, tracing · testing: integration, testcontainers · Docker va
deploy · xavfsizlik: OWASP Top 10 · performance tuning

## 8. Distributed systems
Nega distributed · vaqt va soatlar · replication · consensus va Raft (Go'da yozamiz) · consistency modellari ·
CAP/PACELC · partitioning va sharding · distributed transactions: 2PC, saga · MapReduce (yozamiz)

## 9. System design
Freymvork · URL shortener · rate limiter · chat · news feed · notification tizimi · to'lov tizimi (idempotency) ·
autocomplete · video streaming · taksi tizimi

## 10. Intervyu
Har 5 darsda LeetCode patterni · mock intervyular: coding, system design, behavioral · rezyume · AI davrida
qanday ajralib turish
