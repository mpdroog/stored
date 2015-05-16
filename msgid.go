package main

import (
	"encoding/json"
	"io"
	"bufio"
	"fmt"
	"net/http"
	"github.com/xsnews/webutils/httpd"
	"stored/config"
	"stored/headreader"
	"stored/bodyreader"
	"os"
	"time"
	"bytes"
)

type PostInput struct {
	Msgid string
	Body string
	Meta map[string]string
}

// Add msg to DB
func Post(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	var in PostInput
	if e := json.NewDecoder(r.Body).Decode(&in); e != nil {
		return e
	}
	if in.Msgid == "" || len(in.Body) == 0 {
		httpd.FlushJson(w, httpd.DefaultResponse{
			Status: false, Text: "Missing msgid or body",
		})
		return nil
	}

	// TODO: Crash on day change
	today := time.Now().Format("2006-01-02")
	store := config.Stores[today]
	if _, already := store.Files[in.Msgid]; already {
		httpd.FlushJson(w, httpd.DefaultResponse{
			Status: false, Text: "Already have article " + in.Msgid,
		})
		return nil
	}

	// Write to FS
	{
		f, e := os.Create(store.Basedir + in.Msgid + ".txt")
		if e != nil {
			return e
		}
		defer func() {
			if e := f.Close(); e != nil {
				panic(e)
			}
		}()

		w := bufio.NewWriter(f)
		if _, e := io.Copy(w, bytes.NewBufferString(in.Body)); e != nil {
			return e
		}

		w.Flush()
	}

	config.Stores[today].Files[in.Msgid] = config.File{
		File: in.Msgid,
		Meta: in.Meta,
	}
	if e := config.Save(store); e != nil {
		return e
	}

	if config.Verbose {
		fmt.Println("Saved " + in.Msgid)
	}
	httpd.FlushJson(w, httpd.DefaultResponse{
		Status: true, Text: "Saved",
	})
	return nil
}

// Read msg by msgid
func Get(w http.ResponseWriter, r *http.Request) error {
	msgid := r.URL.Query().Get("msgid")
	if msgid == "" {
		httpd.FlushJson(w, httpd.DefaultResponse{
			Status: false, Text: "msgid not given",
		})
		return nil
	}
	readType := r.URL.Query().Get("type")
	if readType == "" {
		httpd.FlushJson(w, httpd.DefaultResponse{
			Status: false, Text: "type not given",
		})
		return nil
	}
	if readType != "HEAD" && readType != "ARTICLE" && readType != "BODY" {
		httpd.FlushJson(w, httpd.DefaultResponse{
			Status: false, Text: "Type invalid value, valid=[HEAD, ARTICLE, BODY]",
		})
		return nil
	}

	// Check if data in one of the datasets
	var (
		basedir string
		item config.File
		ok bool
		date string
	)
	for d, store := range config.Stores {
		item, ok = store.Files[msgid]
		if ok {
			date = d
			basedir = store.Basedir
			break
		}
	}
	if !ok {
		msg := "Article not found msgid=" + msgid
		fmt.Println("WARN: " + msg)
		httpd.FlushJson(w, httpd.DefaultResponse{
			Status: false, Text: msg,
		})
		return nil
	}	

	path := basedir + item.File + ".txt"
	if config.Verbose {
		fmt.Println("Read " + path)
	}
	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()

	var in io.Reader	
	in = bufio.NewReader(f)
	if readType == "HEAD" {
		in = headreader.New(in)
	} else if readType == "BODY" {
		in = bodyreader.New(in)
	}

	w.Header().Set("Content-Type", "text/plain")
	_, e = io.Copy(w, in)
	if e != nil {
		return e
	}

	// Collect stats
	s, ok := config.Stats[date].Files[msgid]
	if !ok {
		s = config.FileStat{}
	}
	s.Count++
	s.Last = 12 // TODO
	config.Stats[date].Files[msgid] = s
	return nil
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