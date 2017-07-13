package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Disk struct {
	Mountpoint string
	Minfree    string
	MinfreeGB  float64
	Disabled   bool
	Name       string
}
type Config struct {
	General struct {
		HTTPListen              []string `toml:"http_listen"`
		NNTPListen              []string `toml:"nntp_listen"`
		IncomingLog             string   `toml:"incoming_log"`
		ArticleRequestLog       string   `toml:"article_request_log"`
		EnableArticleRequestLog bool     `toml:"enable_article_request_log"`
	}
	Storage []Disk
}

var (
	C        Config
	Verbose  bool
	Appstart time.Time
	Hostname string
	L        *log.Logger
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
	if !strings.HasSuffix(C.General.IncomingLog, "/") {
		C.General.IncomingLog += "/"
	}

	names := make(map[string]bool)
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
		if !strings.HasSuffix(disk.Minfree, "GB") {
			return fmt.Errorf("Mountpoint(%s) only support GB at the moment", disk.Mountpoint)
		}
		minfree, e := strconv.Atoi(C.Storage[i].Minfree[:len(C.Storage[i].Minfree)-2])
		if e != nil {
			return e
		}
		C.Storage[i].MinfreeGB = float64(minfree)
		if e != nil {
			return e
		}

		if _, ok := names[disk.Name]; ok {
			return fmt.Errorf("Mountpoint(%s) has a duplicate name=%s", disk.Mountpoint, disk.Name)
		}
		names[disk.Name] = true
	}
	return
}
