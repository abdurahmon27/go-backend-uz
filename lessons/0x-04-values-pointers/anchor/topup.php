<?php
// PHP tomoni (GO-X.04): massiv qiymat bilan nusxalanadi (Go struct kabi), obyekt esa handle bilan (Go pointer kabi).
$users = [["id" => "42", "name" => "Dilnoza", "balance" => 100]];
foreach ($users as $u) {
    if ($u["id"] === "42") {
        $u["balance"] += 50; // $u — nusxa
    }
}
echo $users[0]["balance"], "\n"; // 100

$objs = [(object) ["id" => "42", "name" => "Dilnoza", "balance" => 100]];
foreach ($objs as $o) {
    if ($o->id === "42") {
        $o->balance += 50; // $o — o'sha obyektning handle'i
    }
}
echo $objs[0]->balance, "\n"; // 150
