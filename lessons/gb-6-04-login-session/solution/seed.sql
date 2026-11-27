-- Namunaviy maʼlumot: ikki foydalanuvchi (GB-6.03'dagi Dilnoza va Jasur). Parollar faqat shu laboratoriya uchun (README):
--   dilnoza@bozor.test — tarvuz-osmon-qalam-ariq-27
--   jasur@bozor.test   — choy-non-qand-27
-- Bazada parolning oʻzi yoʻq, faqat argon2id satri (GB-6.03, internal/password.Hash: m=47104 KiB, t=1, p=1).
-- Yangi foydalanuvchi: POST /users (GB-6.03). Jadvallarni goose yaratadi (migrations/); bu fayl ularni tozalab, qaytadan toʻldiradi.
TRUNCATE sessions, users RESTART IDENTITY;
INSERT INTO users (email, password_hash)
VALUES ('dilnoza@bozor.test', '$argon2id$v=19$m=47104,t=1,p=1$nOpo1eo3/LqYO44wrQWXCw$krxLcr34r7HmNQnM3UHQF6NHozzxymB8qQeRiGq+wQA'),
       ('jasur@bozor.test', '$argon2id$v=19$m=47104,t=1,p=1$ecjvB7e6NhL2MXutDLGfXQ$OkrQiftpgFvIuhs1dpTtHdLt8mGImGy4bFew0fqLJAc');
