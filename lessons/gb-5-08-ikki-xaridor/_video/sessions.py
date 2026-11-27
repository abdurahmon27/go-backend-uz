#!/usr/bin/env python3
"""GB-5.08: ikki psql seansini bosqichma-bosqich ishga tushiradi (A va B) va har buyruqning haqiqiy chiqishini yozadi.
Ishlatish (compose papkasida): python3 sessions.py <stsenariy>   stsenariy: naive | shartli | forupdate | check
Har seans: docker compose exec -T db psql. Buyruqdan keyin \\echo belgisi yuboriladi: chiqish shu belgigacha oʻqiladi.
Kutadigan buyruq (qulf) uchun javob darhol kelmasligi tekshiriladi, keyin boshqa seans tugatgach oʻqiladi."""
import os, select, subprocess, sys, time

class Sess:
    def __init__(self, name):
        self.name, self.n = name, 0
        self.p = subprocess.Popen(["docker", "compose", "exec", "-T", "db", "psql"], stdin=subprocess.PIPE,
                                  stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True, bufsize=0)
        self.buf = ""
        self.mark = None

    def send(self, sql):
        self.n += 1
        self.mark = f"@@{self.name}{self.n}"
        self.p.stdin.write(sql + "\n\\echo " + self.mark + "\n")
        self.p.stdin.flush()

    def _pull(self, timeout):
        r, _, _ = select.select([self.p.stdout], [], [], timeout)
        if r:
            chunk = os.read(self.p.stdout.fileno(), 65536).decode()
            self.buf += chunk
        return bool(r)

    def done(self):
        return self.mark in self.buf

    def result(self, timeout=15):
        end = time.time() + timeout
        while not self.done() and time.time() < end:
            self._pull(0.2)
        if not self.done():
            raise SystemExit(f"{self.name}: javob kelmadi")
        out, _, rest = self.buf.partition(self.mark)
        self.buf = rest.lstrip("\n")
        return out.strip("\n")

    def waits(self, secs=2.0):
        """buyruq javobsiz osilib turganini tekshiradi"""
        end = time.time() + secs
        while time.time() < end:
            self._pull(0.2)
            if self.done():
                return False
        return True

def show(who, sql, out):
    print(f"[{who}] {sql}")
    print(out if out else "(javob yoʻq)")
    print()

def admin(sql):
    return subprocess.run(["docker", "compose", "exec", "-T", "db", "psql", "-c", sql], capture_output=True, text=True).stdout.strip()

SQL_STOCK = "SELECT stock FROM products WHERE id = 4;"

def step(s, sql):
    s.send(sql)
    out = s.result()
    show(s.name, sql, out)
    return out

def waiting(s, sql):
    s.send(sql)
    w = s.waits()
    print(f"[{s.name}] {sql}")
    print(f"(2 soniya javob yoʻq: kutyapti = {w})\n")
    if not w:
        raise SystemExit("buyruq kutishi kerak edi")

def finish_wait(s):
    out = s.result()
    print(f"[{s.name}] (kutgan buyruqning javobi)")
    print(out + "\n")

def main(sc):
    print(admin("UPDATE products SET stock = 1 WHERE id = 4"))
    print(admin(SQL_STOCK) + "\n")
    a, b = Sess("A"), Sess("B")
    if sc == "naive":
        step(a, "BEGIN;"); step(a, SQL_STOCK)
        step(b, "BEGIN;"); step(b, SQL_STOCK)
        step(a, "UPDATE products SET stock = 0 WHERE id = 4;"); step(a, "COMMIT;")
        step(b, "UPDATE products SET stock = 0 WHERE id = 4;"); step(b, "COMMIT;")
        step(a, SQL_STOCK)
    elif sc == "shartli":
        step(a, "BEGIN;")
        step(a, "UPDATE products SET stock = stock - 1 WHERE id = 4 AND stock >= 1;")
        waiting(b, "UPDATE products SET stock = stock - 1 WHERE id = 4 AND stock >= 1;")
        step(a, "COMMIT;")
        finish_wait(b)
    elif sc == "forupdate":
        step(a, "BEGIN;"); step(b, "BEGIN;")
        step(a, "SELECT stock FROM products WHERE id = 4 FOR UPDATE;")
        waiting(b, "SELECT stock FROM products WHERE id = 4 FOR UPDATE;")
        step(a, "UPDATE products SET stock = 0 WHERE id = 4;")
        step(a, "COMMIT;")
        finish_wait(b)
        step(b, "COMMIT;")
    elif sc == "check":
        admin("UPDATE products SET stock = 0 WHERE id = 4")
        step(a, "UPDATE products SET stock = stock - 1 WHERE id = 4;")
    print("oxirgi:", admin(SQL_STOCK))
    for s in (a, b):
        s.p.stdin.close(); s.p.wait(timeout=10)

main(sys.argv[1])
