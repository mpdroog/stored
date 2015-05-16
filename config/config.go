package config

// Read datastores from/into memory.
import (
	"encoding/json"
	"os"
	"time"
	"fmt"
	"strings"
	"path/filepath"
	"bufio"
)

type File struct {
	Meta map[string]string
}
type DB struct {
	Begin time.Time
	Basedir string
	Version int
	Files map[string]File
}

// Minutes since Begin
func (d *DB) Since() int64 {
	return int64(time.Now().Sub(d.Begin) / time.Minute)
}

type FileStat struct {
	Age int64                   // minutes since Begin
	Last int64                // Last read (minutes since Begin)
	Count int                 // Amount of times read
}
type Stat struct {
	Version int
	Files map[string]FileStat
}

var (
	Stats       map[string]Stat
	Stores      map[string]DB
	Verbose     bool
	basedir     string
)

func visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !f.IsDir() {
		// Only read dirs
		return nil
	}
	if basedir == path {
		// Ignore
		return nil
	}
	if Verbose {
		fmt.Println("Load datastore=" + path)
	}

	// db.json
	var day string
	{
		var c DB
		if e := read(path + "/db.json", &c); e != nil {
			return e
		}
		// Ensure slash on path
		if !strings.HasSuffix(c.Basedir, "/") {
			c.Basedir += "/"
		}

		day = c.Begin.Format("2006-01-02")
		Stores[day] = c
	}

	// stats.json
	{
		var c Stat
		if e := read(path + "/stats.json", &c); e != nil {
			return e
		}
		Stats[day] = c
	}

	return nil
}

func Init(path string) error {
	basedir = path
	Stores = make(map[string]DB)
	Stats = make(map[string]Stat)

	// Load existing stores
	if e := filepath.Walk(path, visit); e != nil {
		return e
	}

	// Ensure we got a datastore for today
	today := time.Now()
	if _, ok := Stores[today.Format("2006-01-02")]; !ok {
		if e := Create(time.Now()); e != nil {
			return e
		}
	}
	return nil
}

func Create(date time.Time) error {
	today := date.Format("2006-01-02")
	path := basedir + "/" + today + "/"
	if Verbose {
		fmt.Println("Create datastore=" + path)
	}
	if e := os.MkdirAll(path, 0700); e != nil {
		return e
	}

	// db.json
	{
		d := DB{
			Begin: date,
			Basedir: path,
			Version: 1,
			Files: make(map[string]File),
		}

		f, e := os.Create(path + "db.json")
		if e != nil {
			return e
		}
		defer func() {
			if e := f.Close(); e != nil {
				panic(e)
			}
		}()

		w := bufio.NewWriter(f)
		if e := json.NewEncoder(w).Encode(&d); e != nil {
			return e
		}
		w.Flush()
		Stores[today] = d
	}

	// stats.json
	{
		d := Stat{
			Version: 1,
			Files: make(map[string]FileStat),
		}

		f, e := os.Create(path + "stats.json")
		if e != nil {
			return e
		}
		defer func() {
			if e := f.Close(); e != nil {
				panic(e)
			}
		}()

		w := bufio.NewWriter(f)
		if e := json.NewEncoder(w).Encode(&d); e != nil {
			return e
		}
		w.Flush()
		Stats[today] = d
	}
	return nil
}

func read(f string, v interface{}) error {
	r, e := os.Open(f)
	if e != nil {
		return e
	}
	defer func() {
		if e := r.Close(); e != nil {
			panic(e)
		}
	}()

	if e := json.NewDecoder(r).Decode(v); e != nil {
		return e
	}
	if Verbose {
		fmt.Println(v)
	}
	return e
}

func Save(d DB) error {
	f, e := os.Create(d.Basedir + "db.json")
	if e != nil {
		return e
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()

	w := bufio.NewWriter(f)
	if e := json.NewEncoder(w).Encode(&d); e != nil {
		return e
	}
	w.Flush()
	return nil
}

func SaveStats(basedir string, d Stat) error {
	f, e := os.Create(basedir + "stats.json")
	if e != nil {
		return e
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()

	w := bufio.NewWriter(f)
	if e := json.NewEncoder(w).Encode(&d); e != nil {
		return e
	}
	w.Flush()
	return nil
}