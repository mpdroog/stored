// Insipired on the dotReader from textproto but way simpler
// to improve performance.
package rawio

import (
	"bytes"
	"io"
)

// End Of Stream
var END = []byte("\r\n.\r\n")
var END_SHORT = []byte(".\r\n")

// Ignore EOF
var testNoEOF bool

type DotReader struct {
	r        io.Reader
	begin    bool
	done     bool
	shortEnd bool

	buf []byte
	pos int
}

func New(r io.Reader, shortEnd bool) *DotReader {
	return &DotReader{r: r, begin: true, shortEnd: shortEnd}
}

// Just keep on reading until we found that END.
func (d *DotReader) Read(b []byte) (n int, err error) {
	if d.done {
		return 0, io.EOF
	}
	n, e := d.r.Read(b)
	if testNoEOF && e == io.EOF {
		// Ignore EOF on unittesting
		e = nil
	}
	if n <= 0 {
		// 0 means EOF?
		return 0, io.ErrUnexpectedEOF
	}

	// Remember last 5 bytes so we can find EOF
	// if we receive the EOF in parts..
	if n >= 5 {
		d.buf = b[n-5 : n]
		d.pos = 5
	} else if n > 0 {
		// TODO: User can OOM when slowly adding data?
		// few bytes
		d.buf = append(d.buf, b[:n]...)
		d.pos += n
	}

	if d.begin && d.shortEnd && bytes.Index(b[0:len(END_SHORT)], END_SHORT) == 0 {
		d.done = true
		e = io.EOF
	}
	d.begin = false
	if idx := bytes.Index(d.buf, END); idx != -1 {
		d.done = true
		e = io.EOF
	}
	if !d.done && e == io.EOF {
		// Did not receive end of stream, error!
		e = io.ErrUnexpectedEOF
	}
	return n, e
}
