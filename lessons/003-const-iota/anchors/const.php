<?php
// PHP const: Go kabi, faqat kompilyatsiya paytida maʼlum ifoda
const MAX_SIZE = 10 * 1024 * 1024; // ✅
echo MAX_SIZE, "\n";
const START = time(); // ❌ Fatal error: Constant expression contains invalid operations
