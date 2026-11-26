<?php
// PHP tuzatish (GO-X.04): foreach ($users as &$u) — havola ochiq yoziladi (Go'dagi users[i] / &users[i] kabi).
$users = [["id" => "42", "balance" => 100]];
foreach ($users as &$u) {
    if ($u["id"] === "42") {
        $u["balance"] += 50;
    }
}
unset($u);
echo $users[0]["balance"], "\n"; // 150
