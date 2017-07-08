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
$v = isset($options["v"]);

$ignore = [".DS_Store", ".", "..", "run.php" /* lazy :p */];
$iterate = function($path) use ($ignore, $v, &$iterate) {
	foreach(scandir($path) as $file) {
		if (in_array($file, $ignore)) {
			if ($v) echo "ignore $path/$file\n";
			continue;
		}
		if (substr($file, 0, 1) === "_") {
			if ($v) echo "ignore $path/$file\n";
			continue;
		}

		if (is_dir($file)) {
			$iterate($file);
			continue;
		}
		if ($v) echo "run $path/$file\n";
		require $path."/".$file;
	}
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

$iterate(".");
echo "OK\n";