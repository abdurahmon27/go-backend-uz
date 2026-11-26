// defer'ning JS'dagi tanish tomoni: try…finally. Ishga tushirish: node js/finally.mjs
import { open } from "node:fs/promises";

async function loadUsers(path) {
  const fh = await open(path);
  try {
    return JSON.parse(await fh.readFile("utf8"));
  } finally {
    await fh.close();
  }
}

console.log(await loadUsers("js/users.json"));
