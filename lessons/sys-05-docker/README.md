# SYS-0.05 · Docker: «mening kompyuterimda ishlaydi» muammosiga yechim

Oldingi dars: [SYS-0.02 · Jarayon va muhit: `ps`, `kill`, signallar, `$PATH` va env](../sys-02-process-env/) · Keyingi dars: SYS-0.06 · docker compose: Postgres va Go ilovasi bitta buyruq bilan

| Fayl | Nima |
|---|---|
| `server/` | kichik HTTP server: har bir mehmonni `$DATA/mehmonlar.txt` ga yozadi va nechanchi ekanini aytadi; `PORT` (standart 8077) va `DATA` (standart `data`) muhitdan; SIGTERM'da toza toʻxtaydi (`kut.go`, SYS-0.02 dan) |
| `Dockerfile.bir` | bitta bosqich: `golang:1.26.8` ustida yigʻish va ishga tushirish (969 MB) |
| `Dockerfile` | multi-stage: `golang:1.26.8` da yigʻish, `scratch` ga faqat binary (8.36 MB) |
| `boshqa/` | oʻlchov va xato variantlari: `Dockerfile.cgo` (`CGO_ENABLED=0` yoʻq), `Dockerfile.alpine`, `Dockerfile.distroless` |
| `tuzoq/` | PID 1 tuzogʻi: server `start.sh` (exec'siz) yoki `start-exec.sh` orqali ishga tushadigan alpine image |
| `yigish.sh` | videodagi hamma image'larni yigʻadi va hajmlarini koʻrsatadi |
| `toxtat.sh` | `./toxtat.sh IMAGE [buyruq]`: container'ni ishga tushirib, `docker stop` necha soniya kutganini va chiqish kodini aytadi |
| `homework/` | uy vazifasi javoblari: `bash homework/javoblar.sh`, `homework/guru/Dockerfile` |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/sys-05-docker`). Avval `./yigish.sh` (birinchi marta `golang:1.26.8` image'i yuklanadi,
diskda 874 MB). 8077 port boʻsh boʻlishi kerak.
```bash
./yigish.sh
docker image ls mehmon --format 'table {{.Repository}}:{{.Tag}}\t{{.Size}}'   # 969MB va 8.36MB
for t in bir kichik; do docker run -d --rm --name sinov -p 8078:8077 mehmon:$t >/dev/null; sleep 0.3; echo -n "$t: "; curl -s localhost:8078; docker rm -f sinov >/dev/null; done
docker run --rm mehmon:bir go version                                         # kompilyator image ichida qolgan
docker run --rm mehmon:bir sh -c 'gcc --version | head -1; head -1 /etc/os-release'
docker run --rm mehmon-boshqa:cgo                      # exec /server: no such file or directory
docker run --rm mehmon-boshqa:cgo-build ldd /server    # yigʻish bosqichi (--target build): libc.so.6, ld-linux
docker run -d --name mehmon mehmon:kichik >/dev/null; sleep 0.3; curl -sS localhost:8077          # -p yoʻq: ulanmaydi
docker rm -f mehmon >/dev/null; docker run -d --name mehmon -p 8077:8077 mehmon:kichik >/dev/null; sleep 0.3; curl -s localhost:8077; curl -s localhost:8077
docker ps --filter name=mehmon --format 'table {{.Names}}\t{{.Image}}\t{{.Ports}}'
docker inspect -f '{{.State.Pid}}' mehmon                 # xostdagi PID, masalan 199934
ps -o pid,ppid,comm -p 199934,199908                      # oʻz PID'ingiz va uning PPID'i (ps -o ppid= -p PID)
grep NSpid /proc/199934/status                            # xostdagi PID va 1
docker rm -f mehmon >/dev/null; docker run -d --name mehmon -p 8077:8077 mehmon:kichik >/dev/null; sleep 0.3; curl -s localhost:8077   # yana 1
docker rm -f mehmon >/dev/null; docker run -d --name mehmon -p 8077:8077 -v daftar:/data mehmon:kichik >/dev/null; sleep 0.3; curl -s localhost:8077; curl -s localhost:8077
docker rm -f mehmon >/dev/null; docker run -d --name mehmon -p 8077:8077 -v daftar:/data mehmon:kichik >/dev/null; sleep 0.3; curl -s localhost:8077   # 3
./toxtat.sh mehmon-tuzoq                     # 10.1 s, 137
./toxtat.sh mehmon-tuzoq /start-exec.sh      # 0.1 s, 0
docker rm -f mehmon; docker volume rm daftar # tozalash
```
Natijalar Docker Engine 29.8.1 (overlay2, containerd 2.3.3, runc 1.5.1), `golang:1.26.8` (Debian 13 trixie, gcc 14.2.0),
`alpine:3.24` (3.24.4), `gcr.io/distroless/static-debian13`, Linux 6.18 (NixOS), 2026-11-25. Hamma buyruqlar ikki marta
ishga tushirildi: PID'lar va `ldd` manzillaridan (ASLR) tashqari natijalar bir xil. `docker image ls` dagi SIZE — image
diskda egallaydigan (ochilgan) hajm, siqilgan yuklash hajmi emas. Bizda eski `overlay2` saqlash drayveri; Docker
Engine 29 ning yangi oʻrnatmalarida standart containerd image store, shuning uchun hajmlar sizda biroz farq qilishi mumkin.
Windows/macOS'da Docker Desktop container'larni Linux VM ichida ishga tushiradi: `ps` va `/proc` isboti Linux xostda ishlaydi.

**C kompilyatori haqida.** Bizning kompyuterda (NixOS) `gcc` yoʻq, shuning uchun mahalliy `go vet` / `go build`
`CGO_ENABLED=0` bilan ishlatiladi. Docker ichida esa `golang:1.26.8` da gcc bor: aynan shu farq `savol`/`javob`
sahnasidagi xatoni keltirib chiqaradi.

## Dars matni (videodagi soʻzlar)

<!-- matn: lesson.yaml dan -->
**hook.** Bitta Go server va ikkita image. Birining hajmi toʻqqiz yuz oltmish toʻqqiz megabayt, ikkinchisiniki sakkiz megabaytdan sal ortiq. Ikkalasini ishga tushirib, soʻrov yuboramiz: javob aynan bir xil. Farq Dockerfile'dagi bir necha qatorda. Bugun «mening kompyuterimda ishlaydi» muammosini Docker bilan yopamiz va docker stop nega oʻn soniya osilib qolishini ham koʻramiz.

**muammo.** Nega dastur sizda ishlaydi-yu, serverda yiqiladi? Birinchi Go darsida Go binary'si static ekanini, uni boshqa serverga koʻchirish oson ekanini koʻrgan edik. Lekin dasturni yigʻish uchun aniq Go versiyasi va ba'zan C kutubxonalari kerak. Ishga tushirish uchun esa port, muhit oʻzgaruvchilari va ma'lumot papkasi kerak. Sizda bularning hammasi bor, hamkasbingizda yoki serverda esa boshqacha boʻlishi mumkin. Docker dasturni shu muhit bilan birga bitta paketga yigʻadi, bu paketning nomi image.

**image.** Dockerfile, bu image uchun retsept. docker build undan image yigʻadi. Image ichida container uchun kerakli hamma fayllar va qaysi buyruq ishga tushishi yozilgan. Image oʻzgarmaydi va qatlamlardan iborat. docker run esa image'dan container ishga tushiradi va ustiga uning oʻz yoziladigan qatlamini qoʻshadi. Jarayonlar darsidagi dastur va jarayonni eslang: image dasturga, container esa jarayonga oʻxshaydi. Container virtual mashina emas, xost yadrosidan foydalanadigan oddiy jarayon, buni keyinroq ps bilan isbotlaymiz.

**bir.** Birinchi Dockerfile juda oddiy. FROM asos image'ni tanlaydi: rasmiy golang image, Go bir nuqta yigirma olti nuqta sakkiz. WORKDIR image ichidagi ish papkasi. COPY loyihani image ichiga koʻchiradi, RUN esa build paytida go build'ni bajaradi. CMD container ishga tushganda qaysi buyruq ishlashini aytadi.

**bir-run.** docker build shu retsept boʻyicha image yigʻadi, minus t esa unga nom beradi. Image ichida nima qoldi? Go kompilyatori shu yerda, C kompilyatori gcc ham shu yerda, ostida esa butun Debian. Ular faqat yigʻish uchun kerak edi. Serverni ishga tushirish uchun esa bitta binary yetadi.

**multi.** Yechimning nomi multi-stage build. Endi Dockerfile'da ikkita FROM bor. Birinchi bosqichga build deb nom beramiz. COPY faqat kerakli fayllarni, go.mod va server papkasini koʻchiradi. RUN qatoridagi CGO_ENABLED teng nolga e'tibor bering, unga hozir qaytamiz. Ikkinchi bosqich scratch'dan boshlanadi: bu boʻm-boʻsh image, ichida hatto shell ham yoʻq. COPY minus minus from build birinchi bosqichdan faqat tayyor binary'ni oladi, qolgan hammasi ortda qoladi. ENV serverga ma'lumot papkasini muhit orqali aytadi, EXPOSE esa faqat hujjat. CMD kvadrat qavsda, bu exec shakli: shell yoʻq, server toʻgʻridan-toʻgʻri ishga tushadi.

**savol.** Savol: CGO_ENABLED teng nolni olib tashlasak nima boʻladi? Birinchi Go darsida ldd binary'ni static deb koʻrsatgan edi, lekin cgo ishlatilsa, bu har doim ham shunday emas, degan edik. Server scratch'da ishga tushadimi? Videoni pauza qiling va oʻylab koʻring.

**javob.** Javob: ishlamaydi, xatosi esa juda chalgʻituvchi: no such file or directory. Holbuki fayl joyida turibdi. golang image'da gcc bor, shuning uchun Go cgoʻni oʻzi yoqadi va net paketi C kutubxonasiga bogʻlanadi. Yigʻish bosqichining oʻzida ldd ishga tushiramiz: binary libc'ga va dinamik yuklovchiga muhtoj. Scratch'da ular yoʻq, yadro esa yuklovchini topmasa, aynan shu xatoni qaytaradi. Bu oʻsha «menda ishlaydi» holati: libc bor kompyuterda binary ishlaydi, boʻm-boʻsh image'da esa yoʻq.

**olcham.** Endi toʻrt variantni oʻlchaymiz. Bitta bosqichli image toʻqqiz yuz oltmish toʻqqiz megabayt. Xuddi shu binary alpine ustida qariyb oʻn yetti megabayt: unda shell va paket menejeri bor. Distroless static oʻn bir megabaytga yaqin: shell yoʻq, lekin HTTPS uchun sertifikatlar va vaqt zonalari bor. Scratch esa sakkiz megabaytdan sal ortiq, ya'ni deyarli faqat binary'ning oʻzi. Serverga tashqi HTTPS yoki vaqt zonalari kerak boʻlmasa, scratch yetarli, kerak boʻlsa, distroless qulayroq. Raqamlar sizda biroz farq qilishi mumkin.

**run.** Endi container'ni ishga tushiramiz. docker run minus d uni fonda ishga tushiradi, lekin curl ulana olmadi. Container oʻz tarmogʻida ishlaydi, xostdagi sakkiz ming yetmish yettinchi portda esa hech kim yoʻq. Dockerfile'dagi EXPOSE port ochmaydi. Minus p esa xostdagi portni container ichidagi portga ulaydi, shundan keyin curl javob oldi, ikkinchi mehmon ham. docker ps ishlayotgan container'larni koʻrsatadi: image ls paketlarni, ps esa ishlayotganlarini. Serverda ehtiyot boʻling: minus p portni tashqi dunyoga ham ochadi, faqat oʻzingiz uchun boʻlsa, port oldiga localhost manzilini yozing.

**jarayon.** Endi va'da qilingan isbot. docker inspect container'ning xostdagi PID'ini beradi. ps uni oddiy jarayon sifatida koʻrsatadi, xuddi jarayonlar darsidagidek, ota jarayoni esa containerd-shim. Status faylining NSpid qatori ikkita raqam beradi: xostdagi PID va container ichidagi PID, ya'ni bir. Buni namespace ta'minlaydi: yadro jarayonga oʻz PID raqamlarini, oʻz tarmogʻini va oʻz fayllarini koʻrsatadi. Shuning uchun curl ham avval ulana olmagan edi. Windows va macOS'da Docker Desktop container'larni Linux virtual mashinasi ichida ishga tushiradi, shuning uchun bu isbotni Linux kompyuterda koʻrasiz.

**volume.** Endi ma'lumotni saqlashni koʻramiz. Server har bir mehmonni data papkasidagi faylga yozadi. Container'ni oʻchirib, yangisini ishga tushiramiz va hisob yana birdan boshlanadi. Fayl container'ning yoziladigan qatlamida edi va u bilan birga oʻchdi. Minus v daftar ikki nuqta data esa daftar nomli volume'ni container'ning data papkasiga ulaydi, server ikki mehmonni yozdi. Endi container'ni yana oʻchirib, yangisini ishga tushiramiz. Volume container'dan alohida yashaydi, shuning uchun hisob uchdan davom etdi.

**tuzoq.** Endi jarayonlar darsida va'da qilingan tuzoq. Koʻp loyihalarda server toʻgʻridan-toʻgʻri emas, start skripti orqali ishga tushadi: avval tayyorgarlik, keyin server. Bunday image'da container'ning birinchi jarayoni server emas, shell boʻladi. docker stop esa SIGTERM'ni aynan shu birinchi jarayonga yuboradi. Bizning scratch image'da bu muammo yoʻq: shell umuman yoʻq, server oʻzi birinchi. Repodagi toxtat skripti docker stop qancha kutganini oʻlchaydi.

**tuzoq-run.** docker stop oʻn soniya osilib qolib, bir yuz oʻttiz yetti qaytarsa, server signalni yo umuman olmagan, yo oʻn soniyada ulgurmagan. Bizda server start skripti orqali ishga tushgan, container'da PID bir esa shell. Jarayonlar darsida signalni ushlamagan dastur toʻxtaydi degan edik, PID bir bundan istisno: yadro unga faqat u ushlaydigan signalni yetkazadi. Shell SIGTERM'ni ushlamaydi va serverga uzatmaydi. Oʻn soniyadan keyin SIGKILL keladi: bir yuz yigirma sakkiz plyus toʻqqiz, ya'ni bir yuz oʻttiz yetti. Yechim bitta soʻz, exec: u shell'ni serverning oʻzi bilan almashtiradi. Endi serverdagi kut funksiyasi signalni ushlaydi va server bir zumda toza toʻxtaydi.

**homework.** Uy vazifasi, toʻrt darajada. Oson: serverni PORT teng toʻqqiz ming bilan ishga tushiring, muhit oʻzgaruvchisini minus e beradi, minus p'ni ham shunga moslang. Oʻrta: volume oʻrniga oʻz papkangizni ulang va faylni xostda oching. Qiyin: exec'siz image'ni minus minus init bilan toʻxtating. Tez toʻxtaydi, lekin server signal keldi deb yozadimi? Guru: server root emas, oddiy foydalanuvchi nomidan ishlasin, volume esa baribir yozilsin. Javoblar GitHub'dagi kurs repository'sida, link pastda.

**outro.** Bugun koʻrdik: image oʻzgarmas paket, container esa undan ishga tushgan oddiy jarayon. Multi-stage build kompilyatorni ortda qoldiradi, minus p portni ulaydi, volume esa ma'lumotni container'dan uzoq yashatadi. Keyingi darsda docker compose: Postgres va Go ilovasini bitta buyruq bilan koʻtaramiz. Unda koʻrishamiz, salomat boʻling.
<!-- /matn -->

## Chuqurroq

**Image va container.** Image — container ishga tushishi uchun kerak boʻlgan hamma narsa: fayllar, kutubxonalar va
qaysi buyruqni ishga tushirish. U oʻzgarmaydi (*immutable*) va qatlamlardan (*layers*) iborat. Container image
qatlamlarini oʻzgartirmasdan oladi va ustiga oʻzining yupqa yoziladigan qatlamini qoʻshadi; container oʻchirilsa,
shu qatlam ham oʻchadi. Container — virtual mashina emas: u xost yadrosini boshqa jarayonlar bilan baham koʻradigan
oddiy jarayon, faqat namespace'lar orqali oʻz fayl tizimi, jarayonlari, tarmogʻi va hostname'ini koʻradi
(Docker: *What is a container?*, `namespaces(7)`). `ps` da server oddiy jarayon, ota jarayoni `containerd-shim`;
`/proc/PID/status` dagi `NSpid` uning har bir PID namespace'dagi raqamini beradi: xostda 199934, ichkarida 1. `docker ps` — ishlayotgan container'lar, `docker image ls` — image'lar.

**Nega bitta bosqichli image 969 MB?** `golang:1.26.8` ning oʻzi 874 MB: Debian 13, Go (270 MB), gcc va git.
`RUN go build` qatlami yana 94.9 MB qoʻshdi (`docker history mehmon:bir`): binary va Go'ning build keshi
(`/root/.cache/go-build`, 88 MB). Bularning hech biri serverni ishga tushirish uchun kerak emas.

**Multi-stage.** Bir nechta `FROM`; `COPY --from=build` oldingi bosqichdan faqat kerakli faylni oladi, qolgani yakuniy
image'ga tushmaydi. `scratch` — Docker'ning maxsus boʻsh asosi (uni `pull` yoki `run` qilib boʻlmaydi). Real
loyihada `COPY go.mod go.sum ./` va `RUN go mod download` ni kodni nusxalashdan **oldin** yozing: kod oʻzgarganda
bogʻliqliklar qatlami keshdan olinadi. Bizda tashqi bogʻliqlik yoʻq, shuning uchun bu qadam tushirib qoldirilgan.

**cgo va scratch.** `cmd/cgo`: cgo C kompilyatori topilsa, odatiy yoqilgan (Go 1.20 dan beri kompilyator yoʻq
boʻlsa — oʻchiq). `golang:1.26.8` da gcc bor, `net` paketi esa cgo yoqilganda nomlarni tizim kutubxonasi orqali hal
qilishi mumkin, shuning uchun binary `libc.so.6` va dinamik yuklovchi `/lib64/ld-linux-x86-64.so.2` ga bogʻlanadi
(`ldd`). `scratch` da ular yoʻq; `execve(2)` ELF interpretatori topilmasa `ENOENT` qaytaradi, Docker buni
`exec /server: no such file or directory` deb koʻrsatadi — fayl esa bor. Yechim: `CGO_ENABLED=0` (yoki C
kerak boʻlsa `distroless/base` kabi glibc'li asos). GO-1.01 dagi "`ldd`: not a dynamic executable" — cgo'siz holat.

**Hajmlar va tanlov** (`./yigish.sh`): bitta bosqich 969 MB · alpine 16.8 MB (shell, `apk`) · distroless/static
10.7 MB (CA sertifikatlar, tzdata, `/etc/passwd`, `/tmp`; shell yoʻq) · scratch 8.36 MB (faqat binary).
Server tashqariga HTTPS soʻrov yuborsa, unga CA sertifikatlar kerak: scratch'da ular yoʻq, distroless'da bor.
Ichkarida shell kerak boʻlsa (debug), alpine qulay, lekin hujum yuzasi kattaroq.

**Port.** Container oʻz tarmoq namespace'ida. **Xavfsizlik:** `-p 8077:8077` portni xostning hamma manzillarida, yaʼni
tashqi dunyoga ham ochadi (Docker hujjati: "Publishing container ports is insecure by default"). Faqat oʻz kompyuteringiz
uchun: `-p 127.0.0.1:8077:8077`. `EXPOSE` portni ochmaydi, u faqat hujjat. `-p 8077:8077` xostning
8077-portini container'ning 8077-portiga ulaydi (`-p XOST:CONTAINER`). Docker hujjatiga koʻra, bridge tarmogʻidagi
container portlariga xostning oʻzidan container IP'si orqali ham ulanish mumkin; `localhost:8077` da esa `-p`
siz hech kim tinglamaydi.

**Volume.** Container'ning yoziladigan qatlamidagi maʼlumot container bilan birga yoʻqoladi. Volume'ni Docker
boshqaradi va u container hayotidan tashqarida yashaydi (`-v daftar:/data`). Papkani toʻgʻridan-toʻgʻri ulash —
bind mount (`-v "$PWD/papka:/data"`), uy vazifasiga qarang.

**PID 1 tuzogʻi.** `docker stop` asosiy jarayonga (PID 1) SIGTERM yuboradi, 10 soniyadan keyin SIGKILL.
`pid_namespaces(7)`: namespace'ning "init" jarayoniga (PID 1) faqat u handler oʻrnatgan signallar yetkaziladi
(SIGKILL va SIGSTOP bundan mustasno). `start.sh` dagi shell SIGTERM uchun handler oʻrnatmagan va uni bolasiga
uzatmaydi: server hech narsa bilmaydi, 10 soniyadan keyin SIGKILL, 137 = 128 + 9 (SYS-0.02). `exec /server` shell'ni
server bilan almashtiradi (yangi jarayon yaratmasdan, `bash(1)`: exec), server PID 1 boʻladi va SIGTERM'ni oʻzi ushlaydi
(Docker *best practices*: "uses the exec builtin so that the final running application becomes the container's PID 1").
Bizning asosiy image'da `CMD ["/server"]` (exec shakli, shell yoʻq) — server oʻzi PID 1, 0.1 soniyada toza toʻxtaydi.
Eslatma: alpine'ning `sh` i oddiy `CMD /server` (shell shakli) da serverni `exec` qilib yuboradi (oʻlchadik:
0.1 s, 0), tuzoq esa skriptda bir nechta buyruq boʻlganda chiqadi. Bizning server `signal.NotifyContext` orqali
SIGTERM uchun handler oʻrnatadi, shuning uchun PID 1 boʻlsa ham signalni oladi (oʻlchadik: "signal keldi", 0).

## Manbalar
- Docker: [What is a container?](https://docs.docker.com/get-started/docker-concepts/the-basics/what-is-a-container/),
  [Multi-stage builds](https://docs.docker.com/build/building/multi-stage/), [Base images: scratch](https://docs.docker.com/build/building/base-images/),
  [Dockerfile reference](https://docs.docker.com/reference/dockerfile/) (exec/shell shakli, `EXPOSE`, ENTRYPOINT va PID 1),
  [Port publishing](https://docs.docker.com/engine/network/port-publishing/), [Storage](https://docs.docker.com/engine/storage/),
  [Volumes](https://docs.docker.com/engine/storage/volumes/), [docker stop](https://docs.docker.com/reference/cli/docker/container/stop/),
  [docker run `--init`](https://docs.docker.com/reference/cli/docker/container/run/), [Best practices: ENTRYPOINT/exec, USER](https://docs.docker.com/build/building/best-practices/)
- Go: [cmd/cgo](https://pkg.go.dev/cmd/cgo), [Go 1.20 release notes: cgo](https://go.dev/doc/go1.20#cgo), [net: Name Resolution](https://pkg.go.dev/net#hdr-Name_Resolution),
  [os/signal](https://pkg.go.dev/os/signal), [net/http Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown)
- Linux: [namespaces(7)](https://man7.org/linux/man-pages/man7/namespaces.7.html), [pid_namespaces(7)](https://man7.org/linux/man-pages/man7/pid_namespaces.7.html),
  [proc_pid_status(5)](https://man7.org/linux/man-pages/man5/proc_pid_status.5.html) (`NSpid`), [execve(2)](https://man7.org/linux/man-pages/man2/execve.2.html) (`ENOENT`: ELF interpreter), [bash(1)](https://man7.org/linux/man-pages/man1/bash.1.html) (`exec`)
- [distroless](https://github.com/GoogleContainerTools/distroless) ([static tarkibi](https://github.com/GoogleContainerTools/distroless/blob/main/base/README.md)),
  [tini](https://github.com/krallin/tini) (`--init`), [Yelp: dumb-init](https://engineeringblog.yelp.com/2016/01/dumb-init-an-init-for-docker.html)
- Liz Rice, [Containers From Scratch (GOTO 2018)](https://www.youtube.com/watch?v=8fi7uSYlOdc) — container'ni Go'da noldan yozish

## Uy vazifasi
- 🟢 **Oson:** `mehmon:kichik` ni `PORT=9000` bilan ishga tushiring (muhit oʻzgaruvchisini `-e` beradi) va `-p` ni shunga moslang.
- 🟡 **Oʻrta:** volume oʻrniga oʻz papkangizni ulang (`-v "$PWD/papka:/data"`) va faylni xostda oching. Fayl kimniki boʻldi?
- 🟠 **Qiyin:** exec'siz `mehmon-tuzoq` ni `--init` bilan ishga tushirib toʻxtating. Tez toʻxtaydimi? Server "signal keldi" deb yozadimi?
- 🔴 **Guru:** server root emas, `65532` foydalanuvchisi nomidan ishlasin, volume'ga esa baribir yozilsin.

Avval oʻzingiz urinib koʻring. Javoblar: `bash homework/javoblar.sh` (avval `./yigish.sh`):
- Oson: `docker run -d -e PORT=9000 -p 9000:9000 mehmon:kichik` → `server: :9000 da tinglayapman`, `Salom! Siz 1-mehmonsiz`.
  Ilova portni muhitdan oladi (SYS-0.02), `-p` ning ikkinchi raqami container ichidagi port bilan bir xil boʻlishi kerak.
- Oʻrta: `-v "$PWD/papka:/data"` (bind mount) — fayl xostdagi papkada paydo boʻladi. `--user "$(id -u):$(id -g)"`
  bermasangiz, server container ichida root, fayl ham xostda root'niki boʻladi va uni oddiy foydalanuvchi oʻchira olmaydi.
  Javobda `--user` bor: fayl sizniki.
- Qiyin: `./toxtat.sh --init mehmon-tuzoq` → `docker stop: 0.1 s, chiqish kodi: 143`, lekin "signal keldi" **yoʻq**.
  `--init` PID 1 qilib `docker-init` (tini) ni qoʻyadi; u SIGTERM'ni faqat bevosita bolasiga — shell'ga — uzatadi.
  Shell endi PID 1 emas, SIGTERM uni oʻldiradi (128 + 15 = 143), init chiqadi va yadro namespace'dagi qolgan
  jarayonlarni, jumladan serverni, SIGKILL bilan oʻldiradi (`pid_namespaces(7)`). Tez, lekin toza emas: yechim baribir `exec`
  (yoki tini'ning `-g` rejimi — signalni butun jarayon guruhiga yuborish).
- Guru ([homework/guru/Dockerfile](homework/guru/Dockerfile)): shunchaki `--user 65532:65532` bersangiz,
  `open /data/mehmonlar.txt: permission denied` — yangi volume'ning papkasi root'niki. Yechim: build bosqichida
  `RUN mkdir /data`, yakuniy image'ga `COPY --from=build --chown=65532:65532 /data /data` va `USER 65532:65532`.
  Boʻsh volume birinchi marta ulanganda Docker image'dagi papka tarkibini (egasi bilan) unga koʻchiradi:
  hisob `1`, `2`, `docker top` da `USER 65532`. scratch'da `/etc/passwd` yoʻq, shuning uchun raqamli UID ishlatiladi.
