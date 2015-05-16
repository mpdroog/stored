<?php
require "./init.php";

Api::check();

$res = Api::json("POST", "msgid", array(
	"msgid" => "aaa@bb.cc",
	"meta" => array( "articleid" => "5050" ),
	"body" => rn("Head: value
Head2: value

Body text here"
)));
var_dump($res);

