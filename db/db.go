package db

import (
	"fmt"
	"time"
	"stored/config"
	"os"
	"bufio"
	"io"
	"bytes"

	"stored/headreader"
	"stored/bodyreader"
)

type SaveInput struct {
	Msgid string
	Body string
	Meta map[string]string
}

type ReadInput struct {
	Msgid string
	Type string
}

type ClosingDB struct {
	fd *os.File
	io.Reader
}
func (cb *ClosingDB) Close() error { 
	return cb.fd.Close()
}

func Save(in SaveInput) (error, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	if in.Msgid == "" || len(in.Body) == 0 {
		return fmt.Errorf("Missing msgid or body"), nil
	}

	store, hasStore := config.Stores[today]
	if !hasStore {
		if e := config.Create(now); e != nil {
			return nil, e
		}
	}
	if _, already := store.Files[in.Msgid]; already {
		return fmt.Errorf("Already have article %s", in.Msgid), nil
	}

	// Write to FS
	if !config.MetaOnly {
		f, e := os.Create(store.Basedir + in.Msgid + ".txt")
		if e != nil {
			return nil, e
		}
		defer func() {
			if e := f.Close(); e != nil {
				panic(e)
			}
		}()

		w := bufio.NewWriter(f)
		if _, e := io.Copy(w, bytes.NewBufferString(in.Body)); e != nil {
			return nil, e
		}

		if e := w.Flush(); e != nil {
			return nil, e
		}
	}

	config.Stores[today].Files[in.Msgid] = config.File{
		Meta: in.Meta,
	}
	if e := config.Save(store); e != nil {
		return nil, e
	}
	stat := config.Stats[today].Files[in.Msgid]
	stat.Age = store.Since()
	config.Stats[today].Files[in.Msgid] = stat
	if e := config.SaveStats(store.Basedir, config.Stats[today]); e != nil {
		fmt.Println("WARN: Failed saving stats: " + e.Error())
	}

	if config.Verbose {
		fmt.Println("Saved " + in.Msgid)
	}
	return fmt.Errorf("Saved %s", in.Msgid), nil
}

func Read(in ReadInput) (io.ReadCloser, error, error) {
	msgid := in.Msgid
	readType := in.Type

	if msgid == "" {
		return nil, fmt.Errorf("No msgid given"), nil
	}
	if readType == "" {
		return nil, fmt.Errorf("No type given"), nil
	}
	if readType != "HEAD" && readType != "ARTICLE" && readType != "BODY" {
		return nil, fmt.Errorf("Type invalid value, valid=[HEAD, ARTICLE, BODY]"), nil
	}

	// Check if data in one of the datasets
	var (
		basedir string
		//item config.File
		ok bool
		date string
		store config.DB
	)
	for date, store = range config.Stores {
		_, ok = store.Files[msgid]
		if ok {
			basedir = store.Basedir
			break
		}
	}
	if !ok {
		msg := "Article not found msgid=" + msgid
		fmt.Println("CACHE_MISS: " + msg)
		return nil, fmt.Errorf(msg), nil
	}	

	var f *os.File
	var r io.Reader
	if config.MetaOnly {
		path := basedir + msgid + ".txt"
		if config.Verbose {
			fmt.Println("CACHE_HIT: Read " + path)
		}
		var e error
		f, e = os.Open(path)
		if e != nil {
			return nil, nil, e
		}

		r = bufio.NewReader(f)
		if readType == "HEAD" {
			r = headreader.New(r)
		} else if readType == "BODY" {
			r = bodyreader.New(r)
		}
	}

	// Collect stats
	s, ok := config.Stats[date].Files[msgid]
	if !ok {
		s = config.FileStat{}
	}
	s.Count++
	s.Last = store.Since()
	config.Stats[date].Files[msgid] = s
	if e := config.SaveStats(store.Basedir, config.Stats[date]); e != nil {
		fmt.Println("WARN: Failed saving stats: " + e.Error())
	}

	return &ClosingDB{f, r}, nil, nil
}