package main

import (
	"stored/config"
	"fmt"
	"net/http"
	"github.com/xsnews/webutils/middleware"
	"github.com/xsnews/webutils/muxdoc"
	"github.com/xsnews/webutils/httpd"

	"encoding/json"
	"stored/db"
	"io"
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

// Add msg to DB
func Post(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	var in db.SaveInput
	if e := json.NewDecoder(r.Body).Decode(&in); e != nil {
		return e
	}
	usrErr, sysErr := db.Save(in)
	if sysErr != nil {
		return sysErr
	}

	httpd.FlushJson(w, httpd.DefaultResponse{
		Status: true, Text: usrErr.Error(),
	})
	return nil
}

// Read msg by msgid
func Get(w http.ResponseWriter, r *http.Request) error {
	var in db.ReadInput
	in.Msgid = r.URL.Query().Get("msgid")
	in.Type = r.URL.Query().Get("type")

	read, usrErr, sysErr := db.Read(in)
	if sysErr != nil {
		return sysErr
	}
	if usrErr != nil {
		httpd.FlushJson(w, httpd.DefaultResponse{
			Status: false, Text: usrErr.Error(),
		})
	}
	defer read.Close()

	_, e := io.Copy(w, read)
	return e
}

func Msgid(w http.ResponseWriter, r *http.Request) {
	var e error
	if r.Method == "GET" {
		e = Get(w, r)
	} else if r.Method == "POST" {
		e = Post(w, r)
	} else {
		httpd.FlushJson(w, httpd.DefaultResponse{Status: false, Text: "Unsupported HTTP Method=" + r.Method})
	}

	if e != nil {
		fmt.Println("ERR: " + e.Error())
		httpd.FlushJson(w, httpd.DefaultResponse{Status: false, Text: "Processing error"})
	}
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