package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"stored/client"
	"stored/config"
	"stored/db"
	"stored/rawio"
	"strings"
	"time"

	"stored/bodyreader"
	"stored/headreader"
)

func Quit(conn *client.Conn, tok []string) {
	conn.Send("205 Bye.")
}

func Unsupported(conn *client.Conn, tok []string) {
	log.Printf("WARN: C(%s): Unsupported cmd=%s\n", conn.RemoteAddr(), tok)
	conn.Send("500 Unsupported command")
}

func read(conn *client.Conn, msgid string, msgtype string) {
	// Load(msgid string) (*bytes.Buffer, error) {
	buf, e := db.Load(msgid)
	if e != nil {
		log.Printf("db.Load(%s) e=%s\n", msgid, e.Error())
		conn.Send("500 Failed loading")
		return
	}
	if buf == nil {
		conn.Send("430 No such article")
		return
	}

	// Put reader around it?
	var in io.Reader
	var code string
	if msgtype == "ARTICLE" {
		code = "220"
		in = buf

	} else if msgtype == "HEAD" {
		code = "221"
		in = headreader.New(buf)

	} else if msgtype == "BODY" {
		code = "222"
		in = bodyreader.New(buf, bodyreader.DefaultBufSize)

	} else {
		panic("Should not get here")
	}

	conn.Send(code + " " + msgid)
	if config.Verbose {
		log.Printf("read(%s) start streamreader\n", msgid)
	}
	if _, e := io.Copy(conn.GetWriter(), in); e != nil {
		log.Printf("read(%s) conn.GetWriter=%s\n", msgid, e.Error())
		conn.Close()
		return
	}
	if config.Verbose {
		log.Printf("read(%s) finish streamreader\n", msgid)
	}
	conn.Send("\r\n.") // additional \r\n auto-added
}

func Article(conn *client.Conn, tok []string) {
	if len(tok) != 2 {
		conn.Send("501 Invalid syntax.")
		return
	}
	read(conn, tok[1], "ARTICLE")
}

func Head(conn *client.Conn, tok []string) {
	if len(tok) != 2 {
		conn.Send("501 Invalid syntax.")
		return
	}
	read(conn, tok[1], "HEAD")
}

func Body(conn *client.Conn, tok []string) {
	if len(tok) != 2 {
		conn.Send("501 Invalid syntax.")
		return
	}
	read(conn, tok[1], "BODY")
}

func Ihave(conn *client.Conn, tok []string) {
	if len(tok) != 2 {
		conn.Send("501 Invalid syntax.")
		return
	}
	msgid := tok[1]
	found, e := db.Exists(msgid)
	if e != nil {
		log.Printf("ihave(%s) db.Exists=%s\n", msgid, e.Error())
		conn.Send("436 " + msgid + " Transfer not possible; try again later")
		return
	}
	if found {
		conn.Send("435 " + msgid + " Article not wanted (already have it)")
		return
	}

	// Send them we accept it
	conn.Send("335 " + msgid + " Send article to be transferred")

	if config.Verbose {
		log.Printf("ihave(%s) start streamreader\n", msgid)
	}
	b := new(bytes.Buffer)
	if e := conn.GetDataBlock(b); e != nil {
		log.Printf("Ihave(%s) io.Copy=%s\n", msgid, e.Error())
		conn.Send("400 Failed reading input")
		conn.Close()
		return
	}
	if config.Verbose {
		log.Printf("ihave(%s) finish streamreader\n", msgid)
	}

	r := b.Bytes()
	if len(r)-len(rawio.END) <= 0 {
		log.Printf("takethis(%s) broken msg received\n", msgid)
		conn.Send("436 Failed reading input")
		conn.Close()
		return
	}
	if !bytes.Contains(r, bodyreader.SEP) {
		log.Printf("takethis(%s) no head/body separator found\n", msgid)
		conn.Send("436 No head/body separation found")
		conn.Close()
		return
	}
	b = bytes.NewBuffer(r[:len(r)-len(rawio.END)])

	if e := db.Save(msgid, b); e != nil {
		log.Printf("ihave(%s) db.Save=%s\n", msgid, e.Error())
		conn.Send("436 Failed storing")
		return
	}

	conn.Send("235 " + msgid + " Article accepted")
}

func Check(conn *client.Conn, tok []string) {
	if len(tok) != 2 {
		conn.Send("501 Invalid syntax.")
		return
	}
	msgid := tok[1]
	found, e := db.Exists(msgid)
	if e != nil {
		log.Printf("Check(%s) db.Exists=%s\n", msgid, e.Error())
		conn.Send("431 " + msgid + " Transfer not possible; try again later")
		return
	}
	if found {
		conn.Send("438 " + msgid + " Article not wanted (already have it)")
		return
	}

	// Start reading input
	conn.Send("238 " + msgid + " Send article to be transferred")
}

func Takethis(conn *client.Conn, tok []string) {
	if len(tok) != 2 {
		conn.Send("501 Invalid syntax.")
		return
	}
	msgid := tok[1]

	if config.Verbose {
		log.Printf("takethis(%s) start streamreader\n", msgid)
	}

	b := new(bytes.Buffer)
	if e := conn.GetDataBlock(b); e != nil {
		log.Printf("Takethis(%s) io.Copy=%s\n", msgid, e.Error())
		conn.Send("400 Failed reading input")
		conn.Close()
		return
	}
	if config.Verbose {
		log.Printf("takethis(%s) finish streamreader\n", msgid)
	}

	r := b.Bytes()
	if len(r)-len(rawio.END) <= 0 {
		log.Printf("takethis(%s) broken msg received\n", msgid)
		conn.Send("400 Failed reading input")
		conn.Close()
		return
	}
	if !bytes.Contains(r, bodyreader.SEP) {
		log.Printf("takethis(%s) no head/body separator found\n", msgid)
		conn.Send("436 No head/body separation found")
		conn.Close()
		return
	}
	b = bytes.NewBuffer(r[:len(r)-len(rawio.END)])

	if e := db.Save(msgid, b); e != nil {
		log.Printf("Takethis(%s) db.Save=%s\n", msgid, e.Error())
		conn.Send("400 Failed storing")
		conn.Close()
		return
	}

	conn.Send("239 " + msgid + " Article transferred OK")
}

func Mode(conn *client.Conn, tok []string) {
	if len(tok) != 2 {
		conn.Send("501 Invalid syntax.")
		return
	}
	if strings.ToUpper(tok[1]) != "STREAM" {
		conn.Send("501 Unknown MODE variant")
		return
	}

	conn.Send("203 Streaming permitted")
}

func Date(conn *client.Conn, tok []string) {
	conn.Send("111 " + time.Now().UTC().Format("20060102150405"))
}

func Stat(conn *client.Conn, tok []string) {
	if len(tok) != 2 {
		conn.Send("501 Invalid syntax.")
		return
	}
	msgid := tok[1]
	found, e := db.Exists(msgid)
	if e != nil {
		log.Printf("Stat(%s) db.Exists=%s\n", msgid, e.Error())
		conn.Send("400 " + msgid + " Transfer not possible; try again later")
		return
	}
	if found {
		conn.Send("223 0 " + msgid)
		return
	} else {
		conn.Send("430 Not Found")
	}

}

func req(conn *client.Conn) {
	conn.Send("200 StoreD")
	for {
		tok, e := conn.ReadLine()
		if e != nil {
			log.Printf("WARN: C(%s): %s\n", conn.RemoteAddr(), e.Error())
			break
		}

		// TODO: close conn on error?
		cmd := strings.ToUpper(tok[0])
		if cmd == "QUIT" {
			Quit(conn, tok)
			break
		} else if cmd == "ARTICLE" {
			Article(conn, tok)
		} else if cmd == "HEAD" {
			Head(conn, tok)
		} else if cmd == "BODY" {
			Body(conn, tok)
		} else if cmd == "IHAVE" {
			Ihave(conn, tok)
		} else if cmd == "CHECK" {
			Check(conn, tok)
		} else if cmd == "TAKETHIS" {
			Takethis(conn, tok)
		} else if cmd == "MODE" {
			Mode(conn, tok)
		} else if cmd == "DATE" {
			Date(conn, tok)
		} else if cmd == "STAT" {
			Stat(conn, tok)
		} else {
			Unsupported(conn, tok)
			break
		}
	}

	conn.Close()
	if config.Verbose {
		log.Printf("C(%s) Closed\n", conn.RemoteAddr())
	}
}

func nntpListen(listen string) error {
	sock, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	if config.Verbose {
		log.Printf("nntpd listening on %s\n", listen)
	}

	for {
		conn, err := sock.Accept()
		if err != nil {
			panic(err)
		}
		if config.Verbose {
			log.Printf("C(%s) New\n", conn.RemoteAddr())
		}

		go req(client.New(conn))
	}
}
