<?php
Api::check();

foreach (["<aaa@bb.cc>", "<second@bb.cc>", "<third@bb.cc>"] as $msgid) {
	// Add some data in DB
	$res = Api::json("POST", "msgid", array(
		"msgid" => $msgid,
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
	$res = Api::json("DELETE", "msgid?msgid=$msgid", array());
	if (! $res["status"]) {
		var_dump($res);
		echo "DELETE fails";
	}
}