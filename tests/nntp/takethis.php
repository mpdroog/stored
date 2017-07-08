<?php
$msgs = [];
// 10 random msgs
for ($i = 0; $i < 10; $i++) {
	$msgs[] = generateRandomString()."@rootdev.nl";
}

$body = rn("Date: 2017-07-08
X-TEST: YES

Hello world!");

$nntp = conn();
connWrite($nntp, "MODE STREAM");
assertEquals("203 Streaming permitted", connRead($nntp));

// Send
foreach ($msgs as $msg) {
	connWrite($nntp, "TAKETHIS <$msg>");
	connWrite($nntp, $body . "\r\n.\r\n", true);
}

// Check
foreach ($msgs as $msg) {
	$res = connRead($nntp);
	assertPrefix("239 ", $res);
}

connClose($nntp);