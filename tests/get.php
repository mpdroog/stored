<?php
require "./init.php";

Api::check();

$res = Api::call("GET", "msgid?msgid=<aaa@bb.cc>&type=ARTICLE", "");
var_dump($res);

