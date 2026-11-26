// XATO: json (kichik harf) aniqlanmagan, lekin server baribir ishga tushadi: node xato.js
const express = require("express");
const app = express();

app.get("/health", (req, res) => {
  res.send(json.stringify({ status: "ok" }));
});
app.get("/users/:id", (req, res) => {
  res.json({ id: req.params.id });
});

app.listen(3000, () => console.log("listening on :3000"));
