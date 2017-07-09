<?php
$msgs = [];
// 10 random msgs
for ($i = 0; $i < 10; $i++) {
	$msgs[] = generateRandomString()."@rootdev.nl";
}

$head = rn("Date: 2017-07-08
X-TEST: YES");
$body = rn("Hello world!");
$article = $head."\r\n\r\n".$body;

$nntp = conn();
connWrite($nntp, "MODE STREAM");
assertEquals("203 Streaming permitted", connRead($nntp));

// Send pipelined
foreach ($msgs as $msg) {
	connWrite($nntp, "TAKETHIS <$msg>");
	connWrite($nntp, $article . "\r\n.\r\n", true);
}

// Check response
foreach ($msgs as $msg) {
	$res = connRead($nntp);
	assertPrefix("239 ", $res);
}

// Read articles again
foreach ($msgs as $msg) {
	connWrite($nntp, "BODY <$msg>");
	assertEquals("222 <$msg>", connRead($nntp));
	assertEquals("Hello world!", connRead($nntp));
	assertEquals(".", connRead($nntp));
}

connClose($nntp);