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
	"errors"
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
		panic("Missing msgid or body")
	}

	// TODO: Crash on day change
	today := time.Now().Format("2006-01-02")
	store := config.Stores[today]
	if _, already := store.Files[in.Msgid]; already {
		panic("Already have article " + in.Msgid)
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
	return config.Save(store)
}

// Read msg by msgid
func Get(w http.ResponseWriter, r *http.Request) error {
	msgid := r.URL.Query().Get("msgid")
	if msgid == "" {
		return errors.New("GET msgid not given")
	}
	readType := r.URL.Query().Get("type")
	if readType == "" {
		return errors.New("GET type not given")
	}
	if readType != "HEAD" && readType != "ARTICLE" && readType != "BODY" {
		return errors.New("GET type only support [HEAD, ARTICLE, BODY]")
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
		// TODO: Log so an engineer can fix
		// TODO: Don't report as 'error'
		return errors.New("Article not found msgid=" + msgid)
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
			fmt.Println("WARN: Failed closing file=" + path)
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