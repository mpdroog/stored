<?php
require "./init.php";

Api::check();

// Add some data in DB
$res = Api::json("POST", "msgid", array(
	"msgid" => "<gettest@bb.cc>",
	"meta" => array( "articleid" => "5050" ),
	"body" => base64_encode(rn("Head: value
Head2: value

Body text here"
))));

$res = Api::call("GET", "msgid?msgid=<gettest@bb.cc>&type=ARTICLE", "");
if (md5($res) !== "334d312b3768651d27043e406bbdcb38") {
	var_dump($res);
}
