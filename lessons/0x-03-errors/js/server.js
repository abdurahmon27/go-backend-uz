// GO-X.03 hook: Node'dagi xuddi shu xizmat. config.json'da ortiqcha vergul bor.
const fs = require("node:fs");
const http = require("node:http");

const cfg = JSON.parse(fs.readFileSync("config.json", "utf8"));

http
  .createServer((req, res) => res.end("ok"))
  .listen(cfg.port, () => console.log("listening on :" + cfg.port));
