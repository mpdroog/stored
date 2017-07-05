package db

import (
	"bytes"
	"log"
	"stored/config"
	"fmt"
	"os"
	"io/ioutil"
	"path"
)

// Find file for msgid
func lookup(msgid string) (string, error) {
	if len(msgid) < 3 || msgid[0] != '<' || msgid[len(msgid)-1] != '>' {
		return "", fmt.Errorf("msgid(%s) not valid", msgid)
	}

	msgHash := hash(msgid)
	beginOffset := lookupDisk(msgHash)
	path := lookupPath(msgHash)

	disks := config.C.Storage
	for i := 0; i < len(disks); i++ {
		pos := (i+beginOffset) % len(disks)
		disk := config.C.Storage[pos]
		if disk.Disabled {
			// Disk not used
			continue
		}

		_, e := os.Stat(disk.Mountpoint+path)
		if os.IsNotExist(e) {
			// File not found, try next disk
			continue
		}
		if e != nil {
			log.Printf("DISK_ERR: %s %s\n", disk.Mountpoint, e.Error())
			continue
		}

		// TODO: Save in metric?
		if i == 0 {
			log.Printf("msgid(%s) hash MATCH\n", msgid)
		} else {
			log.Printf("msgid(%s) hash MISMATCH\n", msgid)
		}

		if config.Verbose {
			log.Printf("msgid(%s) resolved to %s\n", msgid, disk.Mountpoint+path)
		}
		return disk.Mountpoint+path, nil
	}

	if config.Verbose {
		log.Printf("msgid(%s) not found\n", msgid)
	}
	return "", nil
}

func Exists(msgid string) (bool, error) {
	path, e := lookup(msgid)
	exists := false
	if path != "" {
		exists = true
	}
	return exists, e
}

// Load message from disk
func Load(msgid string) (*bytes.Buffer, error) {
	path, e := lookup(msgid)
	if e != nil {
		return nil, e
	}
	if path == "" {
		return nil, nil
	}

	// TODO: Re-use buf for perf?
	buf := bytes.NewBuffer(make([]byte, 0, 1024*1024)) //1MB buf
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	// TODO: Defer the read?
	_, e = buf.ReadFrom(f)
	return buf, e
}

// Save message to disk
func Save(msgid string, buf *bytes.Buffer) error {
	if len(msgid) < 3 || msgid[0] != '<' || msgid[len(msgid)-1] != '>' {
		return fmt.Errorf("msgid(%s) not valid", msgid)
	}

	msgHash := hash(msgid)
	beginOffset := lookupDisk(msgHash)
	filepath := lookupPath(msgHash)

	disks := config.C.Storage
	for i := 0; i < len(disks); i++ {
		pos := (i+beginOffset % len(disks))
		log.Printf("X -> %d\n", pos)
		disk := config.C.Storage[pos]
		if disk.Disabled {
			// Skip disk
			continue
		}

		isCreated := false
retry:
		// disk.Mountpoint+path
		e := ioutil.WriteFile(disk.Mountpoint+filepath, buf.Bytes(), 0644)
		if e == nil {
			// Saved!
			if config.Verbose {
				log.Printf("msgid(%s) saved to %s\n", msgid, disk.Mountpoint+filepath)
			}
			return nil
		}

		// Create directory on error (Trick to keep disk I/O low)
		if !isCreated {
			dir := path.Dir(disk.Mountpoint+filepath)
			if config.Verbose {
				log.Printf("msgid(%s) create dir=%s\n", msgid, dir)
			}
			if _, e := os.Stat(dir); os.IsNotExist(e) {
				isCreated = true
				if e := os.MkdirAll(dir, 0777); e != nil {
					log.Printf("msgid(%s) failed mkdir err=%s\n", msgid, e.Error())
				} else {
					goto retry
				}
			}
		}

		log.Printf("msgid(%s) failed writing to %s err=%s\n", msgid, disk.Mountpoint+filepath, e.Error())
	}

	log.Printf("CRIT: Failed writing msgid(%s) to ANY disk\n", msgid)
	return nil
}

// Delete message from disk
func Delete(msgid string) error {
	path, e := lookup(msgid)
	if e != nil {
		return e
	}
	if path == "" {
		return fmt.Errorf("No such msgid")
	}
	e = os.Remove(path)
	if config.Verbose && e == nil {
		log.Printf("msgid(%s) deleted\n", msgid)
	}
	return e
}