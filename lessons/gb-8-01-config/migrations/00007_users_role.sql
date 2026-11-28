-- +goose Up
-- Rol (GB-6.06): oddiy xaridor yoki admin. Yangi foydalanuvchi har doim xaridor: eng kam huquq
-- (OWASP Authorization Cheat Sheet, "Enforce Least Privileges"). Admin'ni faqat baza egasi tayinlaydi.
ALTER TABLE users
    ADD COLUMN role text NOT NULL
    DEFAULT 'xaridor'
    CHECK (role IN ('xaridor', 'admin'));

-- +goose Down
ALTER TABLE users DROP COLUMN role;
