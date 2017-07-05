package main

import (
	"stored/config"
	"flag"
	"time"
	"fmt"
	"io/ioutil"
	"os"
	"bufio"
)

// TODO: Is float64 smart?
func needCleanup(mountpoint string, minfree float64) bool {
	use := DiskUsage(mountpoint)
	if config.Verbose {
		fmt.Printf("Disk(%s) usage=%+v\n", mountpoint, use)
	}

	free := float64(use.Free)/float64(GB)
	if free < minfree {
		fmt.Printf("Disk(%s) cleanup (free=%sGB)\n", mountpoint, free)
		return true
	}
	return false
}

type Message struct {
	Status bool
	Text string
}

func cleanup(disk config.Disk) {
	if !needCleanup(disk.Mountpoint, disk.MinfreeGB) {
		return
	}

	files, e := ioutil.ReadDir(config.C.General.IncomingLog+disk.Name)
	if e != nil {
		fmt.Printf("cleanup e=%s\n", e.Error())
		return
	}

	for _, f := range files {
		if !needCleanup(disk.Mountpoint, disk.MinfreeGB) {
			return
		}
		incomingPath := config.C.General.IncomingLog+disk.Name + "/" + f.Name() +"/incoming.log"
        fmt.Printf("Parse %s/\n", incomingPath)

        fd, e := os.Open(incomingPath)
	    if e != nil {
	        fmt.Printf("cleanup e=%s\n", e.Error())
	    }
	    defer fd.Close()

	    scanner := bufio.NewScanner(fd)
	    for scanner.Scan() {
	    	line := scanner.Text()
	    	if config.Verbose {
		        fmt.Printf("Line=%s\n", line)
		    }

		    msg := new(Message)
	        if e := Delete(line, msg); e != nil {
	        	fmt.Printf("cleanup(msgid=%s) err=%s\n", line, e.Error())
	        }
	        if !msg.Status {
	        	fmt.Printf("cleanup(msgid=%s) res=%s\n", line, msg.Text)
	        }
	    }

	    if e := scanner.Err(); e != nil {
	        fmt.Printf("cleanup(%s) err=%s\n", incomingPath, e.Error())
	    }
    }
}

func main() {
	configPath := ""
	flag.BoolVar(&config.Verbose, "v", false, "Verbose-mode (log more)")
	flag.StringVar(&configPath, "c", "./config.toml", "Path to config.toml")
	flag.Parse()

	if e := config.Init(configPath); e != nil {
		panic(e)
	}

	t := time.Tick(time.Minute * 1)
	select {
		case <- t:
			if config.Verbose {
				fmt.Printf("Minute passed, check disk usage\n")
			}

			for _, disk := range config.C.Storage {
				cleanup(disk)
			}
	}
}