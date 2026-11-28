#!/usr/bin/env python3
"""Terminaldagi Ctrl+C'ni takrorlaydi: buyruqni yangi jarayon guruhida ishga tushiradi, DELAY soniyadan keyin
butun guruhga SIGINT yuboradi (terminal ham Ctrl+C'da shunday qiladi) va buyruqning chiqish kodini qaytaradi.
    python3 _video/ctrlc.py 2.5 go run ./cmd/import katta.txt
"""
import os
import signal
import subprocess
import sys
import time

delay = float(sys.argv[1])
p = subprocess.Popen(sys.argv[2:], start_new_session=True)
time.sleep(delay)
print("^C", end="", flush=True)  # terminal Ctrl+C'ni shunday koʻrsatadi
os.killpg(p.pid, signal.SIGINT)
code = p.wait()
sys.exit(code if code >= 0 else 128 - code)
