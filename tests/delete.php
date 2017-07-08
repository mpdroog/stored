<?php
require "./init.php";

Api::check();

// Add some data in DB
$res = Api::json("POST", "msgid", array(
	"msgid" => "<aaa@bb.cc>",
	"meta" => array( "articleid" => "5050" ),
	"body" => base64_encode(rn("Head: value
Head2: value

Body text here"
))));
if (!$res["status"] && $res["text"] !== "Already have this msg") {
	var_dump($res);
	echo "POST fails";
	exit;
}

// Now delete it
$res = Api::json("DELETE", "msgid?msgid=<aaa@bb.cc>", array());
if (! $res["status"]) {
	var_dump($res);
	echo "DELETE fails";
}
