-- Namunaviy maʼlumot (GB-6.06): uch foydalanuvchi, toʻrt mahsulot, ikki buyurtma.
-- Parollar faqat shu laboratoriya uchun (README):
--   dilnoza@bozor.test — tarvuz-osmon-qalam-ariq-27  (xaridor, GB-6.03/6.04)
--   jasur@bozor.test   — choy-non-qand-27            (xaridor, GB-6.03/6.04)
--   admin@bozor.test   — qovun-anjir-tut-uzum-31     (admin)
-- Bazada parolning oʻzi yoʻq, faqat argon2id satri (GB-6.03, internal/password.Hash: m=47104 KiB, t=1, p=1).
-- Yangi foydalanuvchi: POST /users (GB-6.03) — rolini baza beradi (DEFAULT 'xaridor'), mijoz emas.
-- Jadvallarni goose yaratadi (migrations/); bu fayl ularni tozalab, qaytadan toʻldiradi.
TRUNCATE order_items, orders, sessions, users, products RESTART IDENTITY;
INSERT INTO users (email, password_hash, role)
VALUES ('dilnoza@bozor.test', '$argon2id$v=19$m=47104,t=1,p=1$nOpo1eo3/LqYO44wrQWXCw$krxLcr34r7HmNQnM3UHQF6NHozzxymB8qQeRiGq+wQA', 'xaridor'),
       ('jasur@bozor.test', '$argon2id$v=19$m=47104,t=1,p=1$ecjvB7e6NhL2MXutDLGfXQ$OkrQiftpgFvIuhs1dpTtHdLt8mGImGy4bFew0fqLJAc', 'xaridor'),
       ('admin@bozor.test', '$argon2id$v=19$m=47104,t=1,p=1$XUKfu6Xd7Fcse0J9GFRDuA$PO1/aBbjTlLgep39ILqaMPF8thT1IGnqNdRedkSTFjs', 'admin');
INSERT INTO products (name, price, stock)
VALUES ('olma', 1200000, 40),
       ('non', 400000, 3),
       ('choy', 2500000, 8),
       ('oʻrik', 3500000, 2);
-- Buyurtma ID'lari 41 dan: ishlab turgan doʻkondagidek, undan oldin ham buyurtmalar boʻlgan.
ALTER TABLE orders ALTER COLUMN id RESTART WITH 41;
INSERT INTO orders (user_id, note)
VALUES (1, 'Chilonzor 9-kvartal, 12-uy'),
       (2, 'Yunusobod 4-kvartal, 17-uy, ertalab');
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES (41, 1, 2, 1200000),
       (41, 2, 3, 400000),
       (42, 4, 1, 3500000),
       (42, 3, 1, 2500000);
