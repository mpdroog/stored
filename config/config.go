package config

// Read config.json
import (
	"encoding/json"
	"os"
	"time"
	"fmt"
	"time"
)

type File struct {
	File string
	Age int
}
type DB struct {
	Begin int64
	Basedir string
	Version int
	Files map[string]File
	Refs map[string]map[string]string
}

var (
	C           DB
	Verbose     bool
	Begin       time.Time
)

func Init(f string) error {
	r, e := os.Open(f)
	if e != nil {
		return e
	}
	if e := json.NewDecoder(r).Decode(&C); e != nil {
		return e
	}
	if Verbose {
		fmt.Println(C)
	}
	return nil	
}
