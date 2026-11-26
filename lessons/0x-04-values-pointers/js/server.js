// Hook'ning Node tomoni (GO-X.04): xuddi shu xizmat, kutubxonasiz. for…of bergan user — massivdagi O'SHA obyekt.
// Ishga tushirish: node js/server.js
const http = require("node:http");

const users = [{ id: "42", balance: 100 }];

http
  .createServer((req, res) => {
    const m = req.url.match(/^\/users\/(\w+)(\/topup)?$/);
    if (!m) return res.writeHead(404).end();
    for (const user of users) {
      if (user.id === m[1]) {
        if (req.method === "POST" && m[2]) user.balance += 50;
        return res.end(JSON.stringify(user) + "\n");
      }
    }
    res.writeHead(404).end();
  })
  .listen(3000, () => console.log("listening on :3000"));
