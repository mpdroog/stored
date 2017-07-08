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
fwrite($nntp, "MODE STREAM\r\n");
assertEquals("203 Streaming permitted", stream_get_line($nntp, 999999999999, "\r\n"));

// Send
foreach ($msgs as $msg) {
	echo "TAKETHIS <$msg>\r\n";
	fwrite($nntp, "TAKETHIS <$msg>\r\n");
	fwrite($nntp, $body . "\r\n.\r\n");
}
// Check
foreach ($msgs as $msg) {
	$res = stream_get_line($nntp, 999999999999, "\r\n");
	echo $res . "\r\n";
	assertPrefix("239 ", $res);
}

connClose($nntp);