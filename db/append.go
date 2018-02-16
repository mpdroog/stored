package db

import (
	"fmt"
	"log"
	"os"
	"path"
	"stored/config"
	"time"
)

func AppendArticleRequestLog(msgid string) error {
	// Add msgid to ArticleRequestLog if enabled
	if config.C.General.EnableArticleRequestLog {
		fpath := fmt.Sprintf("%s/articlerequest.log", config.C.General.ArticleRequestLog)

		isCreated := false
	retry:
		f, e := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if e != nil {
			if !isCreated {
				dir := path.Dir(fpath)
				if _, e := os.Stat(dir); os.IsNotExist(e) {
					isCreated = true
					if e := os.MkdirAll(dir, 0777); e != nil {
						log.Printf("appendArticleRequestLog(%s) failed mkdir err=%s\n", msgid, e.Error())
					} else {
						goto retry
					}
				}
			}
			return e
		}
		tmp := fmt.Sprintf("0000-00-00 00:00:00.000  articleage %s 0\r\n", msgid)
		_, e = f.WriteString(tmp)
		f.Close()
		return e
	}
	return nil
}

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
	_, e = f.WriteString(msg)
	f.Close()
	return e
}
