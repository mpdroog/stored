package main

import (
	"stored/config"
	"fmt"
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

func httpListen(listen string) error {
	mux.Title = "StoreD API"
	mux.Desc = "Simple datastore for NNTP Articles."
	mux.Add("/", doc, "This documentation")
	//mux.Add("/meta", article, "PUT Meta set key=?,value=? WHERE msgid=?")
	mux.Add("/msgid", Msgid, "GET message by msgid=? | POST message SET msgid=? AND body=?")
	http.Handle("/", middleware.Use(mux.Mux))

	// TODO: Catch CTRL+C

	if config.Verbose {
		fmt.Println("stored listening on " + listen)
	}
	return http.ListenAndServe(listen, nil)
}