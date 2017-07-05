package db

import (
	"os"
	"stored/config"
	"fmt"
	"time"
	"path"
	"log"
)

func appendLog(disk, msg string) error {
	fpath := fmt.Sprintf("%s%s/%s/incoming.log", config.C.General.IncomingLog, disk, time.Now().Format("2006-01-02"))

	isCreated := false
retry:
	f, e := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if e != nil {
		if !isCreated {
			dir := path.Dir(fpath)
			if _, e := os.Stat(dir); os.IsNotExist(e) {
				isCreated = true
				if e := os.MkdirAll(dir, 0777); e != nil {
					log.Printf("appendLog(%s) failed mkdir err=%s\n", msg, e.Error())
				} else {
					goto retry
				}
			}
		}
		return e
	}
	defer f.Close()

	_, e = f.WriteString(msg);
	return e
}