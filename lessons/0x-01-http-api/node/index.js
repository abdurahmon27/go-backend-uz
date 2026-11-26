// Xuddi shu API Express 5 bilan (taqqoslash uchun): node index.js
const express = require("express");
const app = express();

app.get("/health", (req, res) => {
  res.json({ status: "ok" });
});
app.get("/users/:id", (req, res) => {
  res.json({ id: req.params.id });
});

app.listen(3000, () => console.log("listening on :3000"));
