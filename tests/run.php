<?php
/**
 * Simplified alternative to PHPUnit?
 */

/* Prepare environment */
error_reporting(E_ALL ^ E_STRICT);
ini_set('display_errors', 1);
if (phpversion() < "5.3.0") {
    throw new \Exception("Need at least PHP 5.3.0 for namespaces");
}
$options = getopt("v");
define("VERBOSE", isset($options["v"]));

$ignore = [".DS_Store", ".", "..", "run.php" /* lazy :p */];
$iterate = function($path) use ($ignore, &$iterate) {
	$ok = true;
	foreach(scandir($path) as $file) {
		if (in_array($file, $ignore)) {
			if (VERBOSE) echo "ignore $path/$file\n";
			continue;
		}
		if (substr($file, 0, 1) === "_") {
			if (VERBOSE) echo "ignore $path/$file\n";
			continue;
		}

		if (is_dir($file)) {
			$ok = $iterate($file);
			continue;
		}
		if (VERBOSE) echo "run $path/$file\n";
		//ob_start();
		require $path."/".$file;
		/*$res = ob_get_flush();
		if (strlen($res) > 0) {
			$ok = false;
			echo $res;
		}*/
	}
	return $ok;
};

function assertEquals($a, $b) {
	if ($a !== $b) {
		echo "mismatch $a !== $b\n";
	}
}
function assertPrefix($prefix, $haystack) {
	return mb_substr($haystack, 0, strlen($prefix) ) === $prefix;
}

require "http/_init.php";
require "nntp/_init.php";

if (! $iterate(".")) {
	echo "ERR\n";
	exit(1);
}
echo "OK\n";