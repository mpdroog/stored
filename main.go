package main

import (
	"fmt"
	"stored/config"
	"flag"
	"net/http"
	"github.com/xsnews/webutils/middleware"
	"github.com/xsnews/webutils/muxdoc"
)

var (
	mux    muxdoc.MuxDoc
)

// Return API Documentation (paths)
func doc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, mux.String())
}

func main() {
	configPath := ""
	listen := ""
	flag.BoolVar(&config.Verbose, "v", false, "Verbose-mode (log more)")
	flag.StringVar(&configPath, "c", "./db.json", "Path to config.json")
	flag.StringVar(&listen, "l", "0.0.0.0:9090", "Listen on ip:port")
	flag.Parse()

	if e := config.Init(configPath); e != nil {
		panic(e)
	}	

	mux.Title = "StoreD API"
	mux.Desc = "Simple datastore for NNTP Articles."
	mux.Add("/", doc, "This documentation")
	mux.Add("/article", article, "GET article by id=? | POST article SET id=? AND msgid=?")
	mux.Add("/msgid", msgid, "GET message by msgid=? | POST message SET msgid=? AND body=?")
	http.Handle("/", middleware.Use(mux.Mux))

	if config.Verbose {
		fmt.Println("stored listening on " + listen)
	}
	if e := http.ListenAndServe(listen, nil); e != nil {
		panic(e)
	}
}