<?php
require "./init.php";

Api::check();

$res = Api::json("DELETE", "msgid?msgid=<aaa@bb.cc>", array());
var_dump($res);
