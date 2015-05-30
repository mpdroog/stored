package main

import (
	"stored/config"
	"flag"
)

func main() {
	configPath := ""
	http := ""
	nntp := ""

	flag.BoolVar(&config.Verbose, "v", false, "Verbose-mode (log more)")
	flag.BoolVar(&config.MetaOnly, "m", false, "Only save metadata")
	flag.StringVar(&configPath, "c", "./datastore", "Path to datastore")
	flag.StringVar(&http, "h", "0.0.0.0:9090", "HTTP Listen on ip:port")
	flag.StringVar(&nntp, "n", "0.0.0.0:0909", "NNTP Listen on ip:port")
	flag.Parse()

	if e := config.Init(configPath); e != nil {
		panic(e)
	}

	go func() {
		if e := nntpListen(nntp); e != nil {
			panic(e)
		}
	}()

	if e := httpListen(http); e != nil {
		panic(e)
	}
}