import time

t = time.monotonic()
time.sleep(2)  # soniya
print(round(time.monotonic() - t, 1))  # 2.0
