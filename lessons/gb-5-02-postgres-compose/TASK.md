# Vazifa: test bazasi uchun ikkinchi servis

`compose.yaml`'da `db` servisi tayyor (videodagidek): `postgres:18-alpine`, `127.0.0.1:5482:5432` porti va `pgdata`
volume'i. Testlar asosiy bazani buzmasligi uchun **ikkinchi, alohida Postgres** kerak.

Kerak: Docker Engine va Docker Compose v2 ishlab turishi ([oʻrnatish](https://docs.docker.com/engine/install/); Linux'da
foydalanuvchi `docker` guruhida boʻlishi kerak) va kompyuteringizda `5482` hamda `5483` portlari boʻsh boʻlishi.

## Nima qilish kerak

1. `compose.yaml`'ga `db-test` servisini qoʻshing (`services:` ichida, `db` bilan bir xil daraja). `db` blokini
   nusxalab, **toʻrt narsani** oʻzgartiring:
   - servis nomi: `db-test`;
   - `POSTGRES_DB`: `bozor_test` (foydalanuvchi va parol `bozor` boʻlib qolaveradi);
   - port: `"127.0.0.1:5483:5432"` (chap raqam kompyuterniki, boshqa boʻlishi shart);
   - volume: oʻz nomi bilan, masalan `pgtest:/var/lib/postgresql`, va uni fayl oxiridagi `volumes:` ga ham yozing.
2. Ikkala servisni ishga tushiring: `docker compose up -d`.
3. `bozor.sql`'ni asosiy bazaga yuklang (videodagi buyruq):
   `docker compose exec -T db psql -U bozor -d bozor < bozor.sql`
4. Tekshiring: `go test` oʻtishi kerak.

```bash
docker compose up -d
docker compose exec -T db psql -U bozor -d bozor < bozor.sql
go test ./...
```

Postgres ishga tushguncha bir necha soniya ketadi: `go test` ulanishni 30 soniyagacha oʻzi kutadi. Serverning tayyorligini
oʻzingiz koʻrmoqchi boʻlsangiz: `docker compose logs db | tail -n 1` (`database system is ready to accept connections`).

Ishingiz tugagach hammasini oʻchiring (volume'lar bilan): `docker compose down -v`.

## Tekshiruv

`tekshiruv_test.go` ikki narsani tekshiradi: `docker compose config` chiqishida `db` va `db-test` toʻgʻri
sozlangan (image, port, volume), va ikkala servisga `psql` bilan ulanib soʻrov yuborish mumkin. Boshlangʻich holatda test
Oʻzbekcha xabar bilan yiqiladi: u sizga nima yetishmayotganini aytadi. Testni oʻzgartirmang.

## Maslahatlar

1. Xato xabaridan boshlang: `go test` nima yoʻq deydi?
2. `docker compose config` faylingizni tekshiradi: xatosi boʻlsa aytadi, yoʻq boʻlsa yakuniy shaklini chiqaradi (test ham shuni
   oʻqiydi). YAML'da qatorning ichkariga yozilishi (boʻsh joylar soni) ma'no beradi: `db-test` `services:` ichida, `db` bilan bir
   xil chuqurlikda turishi kerak.
3. `db-test` ham `db` kabi **nomli volume**'ga yozishi kerak, lekin boshqa nomli: aks holda ikki baza bitta fayllarni
   boʻlishadi.

Javob `solution/` papkasida.
