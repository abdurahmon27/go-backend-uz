-- +goose Up
-- Har buyurtmaning egasi bor (GB-6.06): orders.user_id → users.id. Avvalgi customer matni kerak emas:
-- buyurtma kimniki ekanini endi session aytadi, mijoz yozgan ism emas.
-- NOT NULL default'siz qoʻshish faqat boʻsh jadvalda ishlaydi (laboratoriyada orders boʻsh, seed.sql keyin
-- toʻldiradi). Qatorlari bor bazada: avval ustun NULL bilan, eski qatorlarga egasi yoziladi, keyin SET NOT NULL.
ALTER TABLE orders
    ADD COLUMN user_id bigint NOT NULL REFERENCES users (id),
    DROP COLUMN customer;
-- "Shu foydalanuvchining buyurtmalari" (WHERE user_id = …) jadvalni oxirigacha oʻqimasin.
CREATE INDEX orders_user_id_idx ON orders (user_id);

-- +goose Down
DROP INDEX orders_user_id_idx;
ALTER TABLE orders
    ADD COLUMN customer text NOT NULL DEFAULT '',
    DROP COLUMN user_id;
