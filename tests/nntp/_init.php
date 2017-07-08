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
	assertEquals("200 StoreD", stream_get_line($nntp, 999999999999, "\r\n"));
	return $nntp;
}

function connClose($nntp) {
	fwrite($nntp, "QUIT\r\n");
	assertEquals("205 Bye.", stream_get_line($nntp, 999999999999, "\r\n"));
	fclose($nntp);
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