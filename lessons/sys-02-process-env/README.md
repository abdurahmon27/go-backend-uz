# SYS-0.02 · Jarayon va muhit: `ps`, `kill`, signallar, `$PATH` va env

Oldingi dars: [SYS-0.01 · Terminal: pipe, `grep`](../sys-01-terminal/) · Keyingi dars: SYS-0.03 · git: commit — snapshot, branch — bitta hash

| Fayl | Nima |
|---|---|
| `server/` | kichik HTTP server: portni `PORT` muhit oʻzgaruvchisidan oladi (band boʻlsa, xato bilan chiqadi), har soʻrov 1 soniya (`sekin.go`); SIGTERM/Ctrl+C da boshlangan soʻrovlarni tugatib, toza toʻxtaydi (`kut.go`); `KUT=0` boʻlsa signal ushlanmaydi |
| `salom/` | `ISM` muhit oʻzgaruvchisini oʻqiydigan dastur |
| `sinov.sh` | `./sinov.sh TERM\|KILL`: serverni yigʻadi (`./srv`), unga sekin soʻrov yuborib, oʻrtasida signal joʻnatadi |
| `homework/` | uy vazifasi javoblari (`bash homework/javoblar.sh`) |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/sys-02-process-env`) ishga tushiriladi. `./sinov.sh` `./srv` ni ham yigʻadi.
```bash
./sinov.sh TERM            # soʻrov tugaydi, javob keladi, chiqish kodi 0
./sinov.sh KILL            # curl: (52) Empty reply from server, chiqish kodi 137
./srv & sleep 0.3; ps -o pid,ppid,stat,comm -p $$,$!; kill $!; wait
./srv & sleep 0.3; ./srv; ss -ltnpH 'sport = :8077' | grep -o 'users.*'; kill $(pgrep -x srv); wait   # port band: kim egallagan?
sleep 9 & kill $!; wait $!; echo $?          # 143 = 128 + 15
sleep 9 & kill -9 $!; wait $!; echo $?       # 137 = 128 + 9
kill -l 143 137                              # TERM, KILL
KUT=0 ./sinov.sh TERM      # signal ushlanmasa: Go darhol chiqadi, 143, soʻrov uziladi
go run ./salom
ISM=Aziz go run ./salom; echo "[$ISM]"       # faqat shu buyruq uchun
PORT=9090 ./srv & sleep 0.3; tr '\0' '\n' < /proc/$!/environ | grep ^PORT=; kill $!; wait
ISM=Vali; go run ./salom                     # dunyo: export qilinmagan shell oʻzgaruvchisi
export ISM=Vali; go run ./salom              # Vali
bash -c 'export SHAHAR=Toshkent'; echo "[$SHAHAR]"   # []: bola ota-onaning muhitini oʻzgartira olmaydi
type kill curl                               # kill: shell builtin, curl: fayl
salom; echo $?                               # command not found, 127
GOBIN=$PWD/bin go install ./salom; ls bin
PATH="$PATH:$PWD/bin" salom                  # faqat shu buyruq uchun
```
Natijalar `go1.26.8`, GNU bash 5.3, procps `ps`, curl, Linux (`bash -c` orqali) bilan olingan. PID raqamlari har
safar boshqacha, qolgan natijalar bir xil (ikki marta tekshirildi).

**C kompilyatori haqida.** `net/http` odatda cgo bilan yigʻiladi. Bizning kompyuterda (NixOS) `gcc` yoʻq va Go
`CGO_ENABLED=1` bilan oʻrnatilgan, shuning uchun `sinov.sh` serverni `CGO_ENABLED=0 go build` bilan yigʻadi (sof Go,
natija bir xil). Oddiy Linux/macOS'da `go build -o srv ./server` ham yetarli. Bu ham muhit oʻzgaruvchisi bilan
sozlashga misol.

## Dars matni

**Jarayon.** CMU 15-213 taʼrifi: jarayon (*process*) — ishlayotgan dastur ("an instance of a running program"). Bitta dasturdan bir nechta
jarayon boʻlishi mumkin. Har birining raqami (PID), ota jarayoni (PPID), oʻz xotirasi va oʻz muhiti (env) bor.
`cmd &` dasturni fonda ishga tushiradi, `$!` oxirgi fon jarayonining PID'i, `$$` shell'ning oʻzi. `ps -o pid,ppid,stat,comm`
ustunlarni tanlaydi; `STAT` dagi `S` = uxlayapti (hodisa kutyapti), `l` = koʻp oqimli (Go runtime bir nechta OS
thread ishlatadi).

**Jarayonni topish.** Port band boʻlsa (`bind: address already in use`), `ss -ltnp 'sport = :8077'` portni kim
egallaganini PID bilan koʻrsatadi (macOS'da `lsof -i :8077`). Nomini bilsangiz: `pgrep -x srv` (PID), `ps -C srv`.
Koʻp uchraydigan usul `ps aux | grep srv`, lekin u `grep srv` ning oʻzini ham topadi (u ham jarayon, nomida `srv`
bor) — `pgrep` shu sababli qulay. Interaktiv koʻrish uchun `top`/`htop`.

**Signal.** `kill` signal yuboradi; keyin nima boʻlishini signal va dastur hal qiladi. Signal — jarayonga kelgan kichik
xabar, asosan raqamdan iborat (Linux'da `SA_SIGINFO` bilan yuboruvchining PID'i ham keladi, sigaction(2)).
Ctrl+C → terminal oldingi plandagi jarayon(lar) guruhiga SIGINT (2) yuboradi; `kill` odatda SIGTERM (15),
`kill -9` → SIGKILL (9). SIGTERM'ning standart harakati ham "tugatish": dastur uni **ushlasa**, toza chiqishi
mumkin, ushlamasa, natija SIGKILL bilan bir xil (`KUT=0` misoli). SIGKILL'ni ushlab ham, eʼtiborsiz qoldirib ham
boʻlmaydi. (Raqamlar x86/ARM Linux uchun; `kill -l` roʻyxatni beradi.)

**Chiqish kodlari.** Jarayon signal bilan tugasa, bash chiqish kodini `128 + N` qiladi: 130 = Ctrl+C, 137 = SIGKILL,
143 = SIGTERM. Buyruq topilmasa 127. `kill -l 137` kodni nomga qaytaradi. 137 deyarli har doim SIGKILL: yo `kill -9`, yo `docker stop`/Kubernetes'da
kutish vaqti tugagan, yo xotira tugaganda Linux'ning OOM killer'i (`mm/oom_kill.c`: `do_send_sig_info(SIGKILL, …)`; Docker'da `docker inspect` →
`OOMKilled: true`). Dastur oʻzi `exit(137)` qilsa ham 137 chiqadi, lekin bu kam uchraydi.

**Toza toʻxtash.** Go dasturi odatda SIGTERM, SIGINT yoki SIGHUP kelganda darhol chiqadi (`os/signal` hujjati).
`signal.NotifyContext` bu signallarni ushlab, context'ni yopadi; `Server.Shutdown` yangi ulanishlarni qabul qilmaydi
va faol soʻrovlar tugashini kutadi (biz 5 soniya chegara qoʻydik). `KUT=0 ./sinov.sh TERM` buni isbotlaydi: SIGTERM'da 143 va
`curl: (52) Empty reply from server` (server javob bermay ulanishni yopdi), ya'ni foydalanuvchi uchun `kill -9` bilan
bir xil. Server portni `net.Listen` bilan oldin egallaydi: port band boʻlsa, "tinglayapman" deb yolgʻon gapirmay,
xato bilan chiqadi (06-darsdagi yolgʻon log'ni eslang).
Nozik joy: fon jarayonlari interaktiv boʻlmagan shell'da SIGINT'ni eʼtiborsiz qoldiradi (bash(1)); Go `Notify`
bunday signalni qayta yoqadi, lekin `stop()` dan keyin u yana eʼtiborsiz holatiga qaytadi. Shuning uchun
`homework/javoblar.sh` `set -m` ishlatadi.

**Muhit.** Muhit oʻzgaruvchilari — jarayonga beriladigan `NOM=qiymat` juftliklari; Go'da `os.Getenv`.
`ISM=Aziz cmd` qiymatni faqat shu buyruqqa beradi. `ISM=Vali; cmd` esa faqat shell oʻzgaruvchisi: u muhitga
tushmaydi va bolaga oʻtmaydi. `export` uni shu shell'dan keyin ishga tushadigan jarayonlarga beradi. (Videoda bola uchun alohida `SHAHAR` oʻzgaruvchisi olingan: har bir qadam
alohida `bash -c` da ishlaydi; bitta terminalda ketma-ket yozsangiz, oldingi `export ISM=Vali` ham koʻrinib qoladi.) Bola jarayon ota-onasining muhitidan **nusxa** oladi (`environ(7)`), shuning uchun `bash -c 'export ISM=Vali'`
tashqaridagi shell'ga taʼsir qilmaydi. Ishlayotgan jarayonning boshlangʻich muhiti `/proc/PID/environ` da, yozuvlar
nol bayt bilan ajratilgan (SYS-0.01 dagi `-print0` kabi), shuning uchun `tr '\0' '\n'`. Bu fayl jarayon *boshlangandagi*
muhitni koʻrsatadi, keyingi oʻzgarishlarni emas. Faylni jarayon egasi (va root) oʻqiy oladi (ptrace ruxsati):
muhitdagi parol ular uchun sir emas. [06-dars · Funksiyalar](../006-functions/) dagi `PORT=9000` misoli ham shu yoʻl.

**`$PATH`.** Buyruq nomida `/` boʻlmasa, bash avval alias, funksiya va oʻz ichki buyruqlarini (builtin: `kill`,
`echo`, `export`, `type`) tekshiradi, keyin `PATH` dagi papkalardan (`:` bilan ajratilgan) chapdan oʻngga qidiradi
va birinchi bajariladigan faylni ishga tushiradi (topilgan yoʻlni eslab qoladi: `hash`). Shuning uchun
`kill -l 143 137` bash'ning `kill` i bilan ishlaydi (util-linux 2.42 dagi `/usr/bin/kill -l 143 137` esa `kill: too many arguments` deydi; tashqi `kill`
bir nechta boʻlishi mumkin, `type -a kill` hammasini koʻrsatadi); topilmasa — `command not found`, 127. Nomda `/` boʻlsa (`./srv`),
qidiruv boʻlmaydi (POSIX XBD 8.3). Joriy papka `PATH` da odatda yoʻq. `go install` dasturni `GOBIN` ga (odatda
`$HOME/go/bin`) oʻrnatadi, shuning uchun Go oʻrnatish hujjati ([go.dev/doc/install](https://go.dev/doc/install)) bu papkani `PATH` ga
qoʻshishni maslahat beradi (`$HOME/.profile` yoki `$HOME/.bashrc`; macOS zsh'da `~/.zshrc`):
`export PATH="$PATH:$(go env GOPATH)/bin"` (`~/.bashrc` yoki `~/.zshrc` ga). Papkani **oxiriga** qoʻshing:
boshiga qoʻyilgan papkadagi fayl tizimning `ls` yoki `sudo` sini "yashirib" qoʻyishi mumkin. Videodagi
`PATH="…" salom` yozuvi faqat shu bitta buyruq uchun.

**Prodakshnda.** `docker stop`: SIGTERM, 10 soniyadan keyin SIGKILL. Kubernetes: SIGTERM, standart 30 soniya,
keyin SIGKILL. Shuning uchun server SIGTERM'ni ushlab, ishini toza tugatishi kerak. Sozlamani muhitda saqlash —
Twelve-Factor App'ning III qoidasi: kod oʻzgarmaydi, muhit oʻzgaradi (sir kalitlar kodda yoʻq). Docker'da signal faqat
asosiy jarayonga (PID 1) boradi: shell-shakldagi `CMD` da u shell boʻladi va signal'ni dasturingizga uzatmaydi,
natijada `docker stop` 10 soniya kutib, 137 bilan oʻldiradi (Yelp, dumb-init, 2016). Bu SYS-0.05 da.

## Manbalar
- CMU 15-213: [Processes](https://www.cs.cmu.edu/~213/lectures/16-processes.pdf), [Exceptional Control Flow](https://www.cs.cmu.edu/~213/lectures/17-ecf.pdf)
- MIT Missing Semester: [Command-line Environment (2026)](https://missing.csail.mit.edu/2026/command-line-environment/)
- man sahifalari: [signal(7)](https://man7.org/linux/man-pages/man7/signal.7.html), [kill(1)](https://man7.org/linux/man-pages/man1/kill.1.html),
  [bash(1)](https://man7.org/linux/man-pages/man1/bash.1.html) (EXIT STATUS, ENVIRONMENT), [ps(1)](https://man7.org/linux/man-pages/man1/ps.1.html),
  [environ(7)](https://man7.org/linux/man-pages/man7/environ.7.html), [execve(2)](https://man7.org/linux/man-pages/man2/execve.2.html),
  [proc_pid_environ(5)](https://man7.org/linux/man-pages/man5/proc_pid_environ.5.html)
- [POSIX.1-2024 XBD 8.3: PATH](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap08.html)
- [ss(8)](https://man7.org/linux/man-pages/man8/ss.8.html), [sigaction(2)](https://man7.org/linux/man-pages/man2/sigaction.2.html), [termios(3)](https://man7.org/linux/man-pages/man3/termios.3.html) (VINTR → SIGINT),
  [pgrep(1)](https://man7.org/linux/man-pages/man1/pgrep.1.html), Linux [mm/oom_kill.c](https://github.com/torvalds/linux/blob/master/mm/oom_kill.c)
- Go: [os/signal](https://pkg.go.dev/os/signal), [net/http Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown),
  [os.Getenv / LookupEnv](https://pkg.go.dev/os#Getenv), `go help install`, [go.dev/doc/install](https://go.dev/doc/install), [go.dev/doc/code](https://go.dev/doc/code)
- [docker stop](https://docs.docker.com/reference/cli/docker/container/stop/), [Kubernetes: Pod lifecycle](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/),
  [The Twelve-Factor App: III. Config](https://12factor.net/config), [Yelp: dumb-init](https://engineeringblog.yelp.com/2016/01/dumb-init-an-init-for-docker.html)

## Uy vazifasi
- 🟢 **Oson:** `sleep 100` ni Ctrl+C bilan toʻxtating va `echo $?`. Qaysi raqam va nega?
- 🟡 **Oʻrta:** serverni `PORT=9191` bilan ishga tushiring, ikkinchi terminaldan `curl localhost:9191`, keyin `kill`.
- 🟠 **Qiyin:** `salom` `ISM` boʻsh qiymat bilan berilganini (`ISM=`) va umuman berilmaganini farqlasin.
- 🔴 **Guru:** birinchi Ctrl+C da server soʻrovlarni kutsin, ikkinchi Ctrl+C da darhol chiqsin.

Avval oʻzingiz urinib koʻring. Javoblar: `bash homework/javoblar.sh`:
- Oson: **130** = 128 + 2 (SIGINT).
- Oʻrta: `server: :9191 da tinglayapman`, `javob: Salom!`, `kill` dan keyin toza toʻxtash xabarlari.
- Qiyin ([homework/qiyin](homework/qiyin/main.go)): `os.LookupEnv` ikkinchi qiymat qaytaradi: oʻzgaruvchi bormi.
  `os.Getenv` ikkala holatda ham `""` beradi.
- Guru ([homework/guru](homework/guru/main.go)): `<-ctx.Done()` dan keyin darhol `stop()`. `NotifyContext` hujjatiga
  koʻra `stop` signal'ni roʻyxatdan chiqaradi va standart xatti-harakatni qaytarishi mumkin: ikkinchi SIGINT Go
  dasturini darhol toʻxtatadi (chiqish kodi 130, `curl: (52)`). Haqiqiy serverlarda shu usul koʻp ishlatiladi:
  operator ikkinchi marta bossa, kutmaydi.
