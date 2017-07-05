package main

import (
	"bytes"
	"stored/config"
	"fmt"
	"net/http"
	"github.com/itshosted/webutils/middleware"
	"github.com/itshosted/webutils/muxdoc"
	"github.com/itshosted/webutils/httpd"

	"encoding/json"
	"stored/db"
	"io"

	"stored/headreader"
	"stored/bodyreader"
	"log"
	"encoding/base64"
)

var (
	mux    muxdoc.MuxDoc
)

type SaveInput struct {
	Msgid string
	Body string
}

// Return API Documentation (paths)
func doc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, mux.String())
}

// Add msg to DB
func Post(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	var (
		in SaveInput
		e error
	)
	if e = json.NewDecoder(r.Body).Decode(&in); e != nil {
		return e
	}
	if config.Verbose {
		log.Printf("http.Post %+v\n", in)
	}

	raw, e := base64.StdEncoding.DecodeString(in.Body)
	if e != nil {
		return e
	}
	if e := db.Save(in.Msgid, bytes.NewBuffer(raw)); e != nil {
		return e
	}

	httpd.FlushJson(w, httpd.DefaultResponse{
		Status: true, Text: "Saved",
	})
	return nil
}

// Read msg by msgid
func Get(w http.ResponseWriter, r *http.Request) error {
	msgid := r.URL.Query().Get("msgid")
	msgtype := r.URL.Query().Get("type")

	//Load(msgid string) (*bytes.Buffer, error) {
	buf, e := db.Load(msgid)
	if e != nil {
		return e
	}
	if buf == nil {
		// Nothing to send
		httpd.FlushJson(w, httpd.DefaultResponse{
			Status: false, Text: "No such article",
		})
		return nil
	}

	var in io.Reader
	if msgtype == "ARTICLE" {
		in = buf

	} else if msgtype == "HEAD" {
		in = headreader.New(buf)

	} else if msgtype == "BODY" {
		in = bodyreader.New(buf)

	} else {
		httpd.FlushJson(w, httpd.DefaultResponse{
			Status: false, Text: "Invalid msgtype, valid=[ARTICLE, HEAD, BODY]",
		})
		return nil
	}

	_, e = io.Copy(w, in)
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
		log.Printf("ERR: %s\n", e.Error())
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
		log.Printf("httpd listening on %s\n", listen)
	}
	return http.ListenAndServe(listen, nil)
}
