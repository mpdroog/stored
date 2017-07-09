<?php
Api::check();

$msgs = [];
// 10 random msgs
for ($i = 0; $i < 10; $i++) {
	$msgs[] = "<" . generateRandomString()."@rootdev.nl" . ">";
}

$body = rn("Date: 2017-07-08
X-TEST: YES

Hello world!");

// Add some data in DB
foreach ($msgs as $msg) {
	$res = Api::json("POST", "msgid", array(
		"msgid" => $msg,
		"meta" => array( "articleid" => "5050" ),
		"body" => base64_encode(rn("Date: 2017-07-08
Head: value
Head2: value

Body text here"
	))));
	if (! $res["status"]) {
		echo "HTTP POST fail for msgid=$msg\n";
	}
}

$nntp = conn();
foreach ($msgs as $msg) {
	connWrite($nntp, "BODY $msg");
	assertEquals("222 $msg", connRead($nntp));
	assertEquals("Body text here", connRead($nntp));
	assertEquals(".", connRead($nntp));
}

connClose($nntp);