package main

import (
	"bytes"
	"stored/config"
	"net"
	"stored/client"
	"strings"
	"stored/db"
	"stored/rawio"
	"io"
	"log"

	"stored/headreader"
	"stored/bodyreader"
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
		conn.Send("400 No such article")
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
		in = bodyreader.New(buf)

	} else {
		panic("Should not get here")
	}

	conn.Send(code + " " + msgid)
	if _, e := io.Copy(conn.GetWriter(), in); e != nil {
		log.Printf("CRIT: %s\n", e.Error())
		// TODO: conn.close?
		return
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
		conn.Send("436 " + msgid + " Transfer not possible; try again later")
		return
	}
	if found {
		conn.Send("435 " + msgid + " Article not wanted (already have it)")
		return
	}

	// Send them we accept it
	conn.Send("335 " + msgid + " Send article to be transferred")

	b := new(bytes.Buffer)
	if _, e := io.Copy(b, conn.GetReader()); e != nil {
		conn.Send("436 Failed reading input")
		return
	}
	r := b.Bytes()
	b = bytes.NewBuffer(r[:len(r) - len(rawio.END)])

	if e := db.Save(msgid, b); e != nil {
		conn.Send("436 Failed storing e=" + e.Error())
		return
	}

	conn.Send("235 " + msgid)
}

func Check(conn *client.Conn, tok []string) {
	if len(tok) != 2 {
		conn.Send("501 Invalid syntax.")
		return
	}
	msgid := tok[1]
	found, e := db.Exists(msgid)
	if e != nil {
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
	b := new(bytes.Buffer)
	if _, e := io.Copy(b, conn.GetReader()); e != nil {
		conn.Send("400 Failed reading input") // TODO: wrong code?
		return
	}

	r := b.Bytes()
	b = bytes.NewBuffer(r[:len(r) - len(rawio.END)])

	if e := db.Save(msgid, b); e != nil {
		conn.Send("400 Failed storing e=" + e.Error()) // TODO: wrong code?
		return
	}

	conn.Send("239 " + msgid)
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
		} else if (cmd == "TAKETHIS") {
			Takethis(conn, tok)
		} else if (cmd == "MODE") {
			Mode(conn, tok)
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