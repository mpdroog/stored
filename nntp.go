package main

import (
	"fmt"
	"stored/config"
	"net"
	"stored/client"
	"strings"
	"stored/db"
	"io"
)

func Quit(conn *client.Conn, tok []string) {
	conn.Send("205 Bye.")
}

func Unsupported(conn *client.Conn, tok []string) {
	fmt.Println(fmt.Sprintf("WARN: C(%s): Unsupported cmd %s", conn.RemoteAddr(), tok[0]))
	conn.Send("500 Unsupported command")
}

func read(conn *client.Conn, msgid string, msgtype string) {
	read, usrErr, sysErr := db.Read(
		db.ReadInput{Msgid: msgid[1:len(msgid)-1], Type: msgtype},
	)
	if sysErr != nil {
		fmt.Println("WARN: " + sysErr.Error())
		conn.Send("500 Failed processing")
		return
	}
	if usrErr != nil {
		conn.Send("400 " + usrErr.Error())
		return
	}
	defer read.Close()

	var code string
	if msgtype == "ARTICLE" {
		code = "220"
	} else if msgtype == "HEAD" {
		code = "221"
	} else if msgtype == "BODY" {
		code = "222"
	} else {
		panic("Should not get here")
	}

	conn.Send(code + " " + msgid)
	if _, e := io.Copy(conn.GetWriter(), read); e != nil {
		fmt.Println("WARN: " + e.Error())
		conn.Send("500 Failed forwarding")
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

func req(conn *client.Conn) {
	conn.Send("200 StoreD")
	for {
		tok, e := conn.ReadLine()
		if e != nil {
			fmt.Println(fmt.Sprintf("WARN: C(%s): %s", conn.RemoteAddr(), e.Error()))
			break
		}

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
		} else {
			Unsupported(conn, tok)
			break
		}
	}

	conn.Close()
	if config.Verbose {
		fmt.Println(fmt.Sprintf("C(%s) Closed", conn.RemoteAddr()))
	}
}

func nntpListen(listen string) error {
	sock, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	if config.Verbose {
		fmt.Println("nntpd listening on " + listen)
	}

	for {
		conn, err := sock.Accept()
		if err != nil {
			panic(err)
		}
		if config.Verbose {
			fmt.Println(fmt.Sprintf("C(%s) New", conn.RemoteAddr()))
		}

		go req(client.New(conn))
	}
}