<?php
/**
 * Global settings.
 */
define("API_KEY", "Abc");

/* Prepare environment */
error_reporting(E_ALL ^ E_STRICT);
ini_set('display_errors', 1);
if (phpversion() < "5.3.0") {
    throw new \Exception("Need at least PHP 5.3.0 for namespaces");
}

// Convert \n to \r\n
function rn($subject) {
    return str_replace("\n", "\r\n", $subject);
}

/**
 * Simple API abstraction.
 *
 * @author mdroog <pw.droog@quicknet.nl>
 */
class Api {
    /** @var string Base URL */
    private static $URL = "http://localhost:9090/";
    /** @var string User agent */
    private static $UA = "UA";

    /** @var FALSE|resource */
    private static $_contect = FALSE;

    /**
     * Create new persistent connection
     *  to server.
     */
    public static function init() {
        self::$_contect = curl_init(NULL);
    }

    /**
     * Close connection.
     */
    public static function close() {
        curl_close(self::$_contect);
        unset(self::$_contect);
    }

    /**
     * Check if environment is set up correctly.
     *
     * @throws Exception
     */
    public static function check() {
        if (! function_exists("curl_init")) {
            throw new \Exception("Missing php5-curl: apt-get install php5-curl");
        }
    }

    /**
     * Send JSON-request to API.
     *
     * @param string $method GET/POST/PUT/DELETE
     * @param string $path Like shown on finance.itshosted.nl/doc
     * @param array $fields Data to store/process
     *
     * @return array Json decoded array ready for use
     * @throws Exception On any problem
     */
    public static function json($method, $path, array $fields) {
        $json = json_encode($fields);
        var_dump($json);
        $raw = self::call($method, $path, $json);
        if ($raw === "null") {
            return null;
        }
        $ret = json_decode($raw, TRUE);
        if (! is_array($ret)) {
            throw new \Exception("Failed decoding server-msg: $raw");
        }
        return $ret;
    }

    /**
     * Read customer/server IP.
     *
     * @param boolean $remote TRUE=customerIP, FALSE=serverIP
     *
     * @return string IP-address or "NO"
     */
    private static function getIp($remote) {
        if ($remote) {
            if (isset($_SERVER["REMOTE_ADDR"])) {
                return $_SERVER["REMOTE_ADDR"];
            }
            return "NO";
        }
        if (isset($_SERVER['SERVER_ADDR'])) {
            return $_SERVER['SERVER_ADDR'];
        } else {
            return gethostbyname(gethostname());
        }
    }

    /**
     * HTTP Request.
     *
     * @param string $method HTTP Method
     * @param string $path Suffix for HTTP URL
     * @param string $body Data to write (if len($body) > 0)
     * @param string $contentType JSON
     *
     * @return mixed Response
     * @throws Exception On any problem
     */
    private static function call($method, $path, $body, $contentType = "application/json") {
        if (self::$_contect === FALSE) {
            self::init();
        }
        $method = strtoupper($method);
        $ch = self::$_contect;
        if ($ch === FALSE) {
            throw new \Exception("Error calling curl_init");
        }

        $opt = TRUE;
        $key = (strpos($path, "?") === FALSE ? "?" : "&") . "key=" . API_KEY;
        $opt &= curl_setopt($ch, CURLOPT_URL, self::$URL . $path . $key);
        if (strlen($body) > 0) {
            $opt &= curl_setopt($ch, CURLOPT_POSTFIELDS, $body);
        }
        $opt &= curl_setopt($ch, CURLOPT_USERAGENT, self::$UA);
        $opt &= curl_setopt($ch, CURLOPT_RETURNTRANSFER, TRUE);
        $opt &= curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 1000);
        $opt &= curl_setopt($ch, CURLOPT_TIMEOUT, 1000);
        $opt &= curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, 2);
        $opt &= curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $method);
        $opt &= curl_setopt($ch, CURLOPT_HTTPHEADER, array(
            "X-IP-Customer: " . self::getIp(TRUE),
            "X-IP-Server: " . self::getIp(FALSE),
            "X-PHP-Ver: " . phpversion()
        ));
        $opt &= curl_setopt($ch, CURLOPT_HTTPHEADER, array(
            "Content-Type: " . $contentType
        ));

        if ($opt == FALSE) {
            throw new \Exception(
                "One or more cURL option-flags failed"
            );
        }
        $result = curl_exec($ch);
        if ($result === FALSE) {
            throw new \Exception(
                "Error calling curl_exec #" . curl_errno($ch) . curl_error($ch)
            );
        }
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);

        if ($httpCode !== 200) {
            if ($httpCode === 401) {
                throw new \Exception(
                    "Missing key-argument in URL"
                );
            } else if ($httpCode === 403) {
                throw new \Exception(
                    "Invalid API key"
                );
            } else if ($httpCode === 429) {
                throw new \Exception(
                    "HTTP Ratelimit (429) reached, doing too much HTTP-requests/per minute!"
                );
            } else if ($httpCode === 502) {
                throw new \Exception(
                    "Server down? (502 means proxy error)"
                );
            }
            throw new \Exception("HTTP Wrong response-code ($httpCode) $result");
        }

        return $result;
    }
}
