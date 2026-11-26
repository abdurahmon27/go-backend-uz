# Python tomoni (GO-X.04): for bergan user — ro'yxatdagi O'SHA dict (nom obyektga bog'lanadi).
users = [{"id": "42", "name": "Dilnoza", "balance": 100}]
for user in users:
    if user["id"] == "42":
        user["balance"] += 50
print(users[0]["balance"])  # 150
