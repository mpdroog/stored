package client

import (
	"bufio"
	"fmt"
	"net"
	"nntpd/config"
	"strings"
	"io"
	"nntpd/rawio"
)

const EOF = "\r\n"

type Conn struct {
	LoggedIn   bool
	Cmds       int

	conn net.Conn
	r    *bufio.Reader
	w    *bufio.Writer
}

func (c *Conn) Send(cmd string) error {
	if config.Verbose {
		fmt.Println(fmt.Sprintf("C(%s) >> %s", c.RemoteAddr(), cmd))
	}
	_, e := c.w.WriteString(cmd + EOF)
	if e != nil {
		return e
	}
	return c.w.Flush()
}

func (c *Conn) ReadLine() ([]string, error) {
	s, e := c.r.ReadString('\n')
	if e != nil {
		return nil, e
	}
	s = s[:len(s)-2] // Strip \r\n
	if config.Verbose {
		fmt.Println(fmt.Sprintf("C(%s) << %s", c.RemoteAddr(), s))
	}
	tok := strings.Split(s, " ")
	return tok, nil
}

func (c *Conn) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Conn) LocalAddr() string {
	return c.conn.LocalAddr().String()
}

// Get DotReader
func (c *Conn) GetReader() io.Reader {
	return rawio.New(c.r)
}

func (c *Conn) GetWriter() io.Writer {
	return c.w
}

func (c *Conn) Close() error {
	if config.Verbose {
		fmt.Println(fmt.Sprintf("C(%s) CLOSED", c.RemoteAddr()))
	}
	return c.conn.Close()
}

func New(conn net.Conn) *Conn {
	return &Conn{
		conn:      conn,
		r:         bufio.NewReader(conn),
		w:         bufio.NewWriter(conn),
	}
}
