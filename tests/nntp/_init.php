<?php
function conn() {
	$nntp = stream_socket_client("tcp://127.0.0.1:9091", $errno, $errorMessage, 2);
	if ($nntp === false) {
		user_error("stream_socket_client failed: $errno $errorMessage");
	}
	// default to 120seconds
	if (stream_set_timeout($nntp, 120) === false) {
		user_error("stream_set_timeout failed");
	}
	assertEquals("200 StoreD", connRead($nntp));
	return $nntp;
}

function connClose($nntp) {
	connWrite($nntp, "QUIT");
	assertEquals("205 Bye.", connRead($nntp));
	fclose($nntp);
}

function connWrite($nntp, $msg, $bin=false) {
	if (VERBOSE && !$bin) {
		echo ">> $msg\n";
	}
	$eol = "\r\n";
	if ($bin) {
		$eol = "";
	}

	return fwrite($nntp, "$msg$eol");
}
function connRead($nntp) {
	$msg = stream_get_line($nntp, 999999999999, "\r\n");
	if (VERBOSE) {
		echo "<< $msg\n";
	}
	return $msg;
}

function generateRandomString($length = 10) {
    $characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    $charactersLength = strlen($characters);
    $randomString = '';
    for ($i = 0; $i < $length; $i++) {
        $randomString .= $characters[rand(0, $charactersLength - 1)];
    }
    return $randomString;
}