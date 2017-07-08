package main

import (
	"stored/config"
	"flag"
	"log"
	_ "net/http/pprof"
)

func main() {
	configPath := ""
	flag.BoolVar(&config.Verbose, "v", false, "Verbose-mode (log more)")
	flag.StringVar(&configPath, "c", "./config.toml", "Path to config.toml")
	flag.Parse()

	if e := config.Init(configPath); e != nil {
		panic(e)
	}
	if config.Verbose {
		log.Printf("Config=%+v\n", config.C)
	}

	go func() {
		if e := nntpListen(config.C.General.NNTPListen[0]); e != nil {
			panic(e)
		}
	}()

	if e := httpListen(config.C.General.HTTPListen[0]); e != nil {
		panic(e)
	}
}