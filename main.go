package main

import (
	"stored/config"
	"flag"
)

func main() {
	configPath := ""
	listen := ""
	flag.BoolVar(&config.Verbose, "v", false, "Verbose-mode (log more)")
	flag.StringVar(&configPath, "c", "./datastore", "Path to datastore")
	flag.StringVar(&listen, "l", "0.0.0.0:9090", "Listen on ip:port")
	flag.Parse()

	if e := config.Init(configPath); e != nil {
		panic(e)
	}

	if e := httpListen(listen); e != nil {
		panic(e)
	}
}