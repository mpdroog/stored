package config

import (
	"github.com/BurntSushi/toml"
	"time"
	"os"
	"log"
	"fmt"
	"strings"
)

type Disk struct {
	Mountpoint string
	Minfree string
	Disabled bool
}
type Config struct {
	General struct {
		HTTPListen []string `toml:"http_listen"`
		NNTPListen []string `toml:"nntp_listen"`
	}
	Storage []Disk
}

var (
	C           Config
	Verbose     bool
	Appstart    time.Time
	Hostname    string
	L           *log.Logger
)

func Init(f string) error {
	Appstart = time.Now()
	r, e := os.Open(f)
	if e != nil {
		return e
	}
	defer r.Close()
	if _, e := toml.DecodeReader(r, &C); e != nil {
		return fmt.Errorf("TOML: %s", e)
	}

	if e := parseConfig(); e != nil {
		return e
	}

	Hostname, e = os.Hostname()
	if e != nil {
		return e
	}

	L = log.New(os.Stdout, "", log.LstdFlags)
	return nil
}

func parseConfig() (e error) {
	if len(C.General.HTTPListen) != 1 {
		return fmt.Errorf("HttpListen only supports 1 listener")
	}
	if len(C.General.NNTPListen) != 1 {
		return fmt.Errorf("HttpListen only supports 1 listener")
	}

	for i, disk := range C.Storage {
		if disk.Disabled {
			continue
		}
		stat, e := os.Stat(disk.Mountpoint)
		if os.IsNotExist(e) {
			return fmt.Errorf("Mountpoint(%s) invalid path", disk.Mountpoint)
		}
		if !stat.IsDir() {
			return fmt.Errorf("Mountpoint(%s) not directory", disk.Mountpoint)
		}

		if !strings.HasSuffix(disk.Mountpoint, "/") {
			C.Storage[i].Mountpoint += "/"
		}
	}
	return
}