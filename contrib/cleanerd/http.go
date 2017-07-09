package main

// Simple JSON/HTTP abstraction.
// So we can remove a lot of duplicate code
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"stored/config"
	"encoding/json"
)

var client *http.Client

func init() {
	client = &http.Client{}
}

func Delete(url string, out interface{}) error {
	url = "http://127.0.0.1/msgid?msgid=" + url
	if config.Verbose {
		fmt.Printf("HTTP.Delete url=%s\n", url)
	}

	req, e := http.NewRequest("DELETE", url, nil)
	if e != nil {
		return e
	}
	res, e := client.Do(req)
	if e != nil {
		return e
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		msg, e := ioutil.ReadAll(res.Body)
		if e != nil {
			return fmt.Errorf("ReadAll(HTTP.Body) err=%s", e.Error())
		}
		return fmt.Errorf("StatusCode 200 expected, received=%s with msg=%s", res.Status, msg)
	}
	
	dec := json.NewDecoder(res.Body)
	return dec.Decode(&out)
}
